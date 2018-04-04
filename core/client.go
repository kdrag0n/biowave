package core

import (
	"errors"
	"fmt"
	"github.com/iancmcc/go-datastructures/bitarray"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/getsentry/raven-go"
	"github.com/kdrag0n/discordgo"

	// HTTP pprof handlers
	_ "net/http/pprof"
)

// Client is the full client context.
type Client struct {
	Config     Config
	Sessions   []*discordgo.Session
	Commands   map[string]*Command
	Developers bitarray.BitArray

	StartTime  time.Time

	// Data
	DB         *badger.DB
	IsDBClosed bool

	// Emotes
	EmoteOk      string
	EmoteFail    string
	EmoteGrave   string
	EmoteLoading string
	EmoteBot     string

	ourMention      string
	ourGuildMention string
	ownerID         uint64

	isReady uint32
}

// NewClient creates a new bot client.
func NewClient(config Config) (*Client, error) {
	// verify config
	if l := len(config.Token); l < 56 || l > 64 {
		return nil, errors.New("Invalid token")
	} else if config.Shards < 1 || config.Shards == 2 {
		return nil, errors.New("Invalid shard count, cannot be 2 or below 1")
	} else if len(config.DatabasePath) == 0 {
		return nil, errors.New("Invalid database path")
	} else if l := len(config.DefaultPrefix); l == 0 || l > 32 {
		return nil, errors.New("Invalid default prefix, empty or over 32 characters")
	}

	// open database
	opts := badger.DefaultOptions
	opts.Dir = config.DatabasePath
	opts.ValueDir = config.DatabasePath

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	// create DiscordGo sessions for shards
	sessions := make([]*discordgo.Session, config.Shards)

	for id := range sessions {
		dg, _ := discordgo.New() // error is impossible
		dg.Token = "Bot " + config.Token
		dg.MFA = false
		dg.ShouldReconnectOnError = true
		dg.Compress = true
		dg.ShardID = id
		dg.ShardCount = config.Shards
		dg.StateEnabled = true
		dg.SyncEvents = true
		dg.MaxRestRetries = 3

		sessions[id] = dg
	}

	devSet := bitarray.NewSparseBitArray()
	for _, devID := range config.Developers {
		devSet.SetBit(devID)
	}

	return &Client{
		Config:     config,
		Sessions:   sessions,
		Commands:   make(map[string]*Command, 120),
		Developers: devSet,

		StartTime: time.Now(),

		DB:         db,
		IsDBClosed: false,

		EmoteOk:      "✅",
		EmoteFail:    "❌",
		EmoteGrave:   "⚰",
		EmoteLoading: "⌛",
		EmoteBot:     "bot",

		ourMention:      "<@0>",
		ourGuildMention: "<@!0>",
		ownerID:         0,

		isReady: 0,
	}, nil
}

// ForSessions executes a function with every session.
func (c *Client) ForSessions(iter func(*discordgo.Session)) {
	for _, dg := range c.Sessions {
		iter(dg)
	}
}

// Start initiates the client's connections.
func (c *Client) Start() (result error) {
	c.LoadModules()
	go c.housekeeper()

	if c.Config.Pprof != 0 {
		go func() {
			defer c.ErrorHandler("pprof http server")

			err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", c.Config.Pprof), nil)
			if err != nil {
				panic(err)
			}
		}()
		Log.Info("pprof server started", zap.Uint16("port", c.Config.Pprof))
	}

	var wg sync.WaitGroup
	wg.Add(len(c.Sessions))

	for idx, dg := range c.Sessions {
		go func(idx int, dg *discordgo.Session) {
			defer wg.Done()

			dg.AddHandler(c.onMessage)
			dg.AddHandlerOnce(c.onReady) // don't want to fire again on connect

			err := dg.Open()
			if err != nil {
				Log.Error("error opening session", zap.Int("shard", idx), zap.Error(err))
				multierr.Append(result, err)
			}

			Log.Info("started", zap.Int("shard", idx))
		}(idx, dg)

		time.Sleep(4 * time.Second)
	}

	wg.Wait()

	if result == nil {
		c.UpdateStatus()
	}

	return
}

// Stop stops the client and all associated Discord sessions.
func (c *Client) Stop() (result error) {
	// close sessions
	var wg sync.WaitGroup
	wg.Add(len(c.Sessions))

	for idx, dg := range c.Sessions {
		go func(idx int, dg *discordgo.Session) {
			defer wg.Done()

			err := dg.Close()
			if err != nil {
				Log.Error("error closing session", zap.Int("shard", idx), zap.Error(err))
				multierr.Append(result, err)
			}

			Log.Info("stopped", zap.Int("shard", idx))
		}(idx, dg)
	}

	wg.Wait()

	// unload modules
	for _, module := range modules {
		err := module.Unload(c)
		if err != nil {
			Log.Error("error unloading module", zap.String("module", module.Name), zap.Error(err))
			multierr.Append(result, err)
		}
	}

	// close database
	err := c.DB.Close()
	if err != nil {
		Log.Error("error closing database", zap.Error(err))
		multierr.Append(result, err)
	}
	c.IsDBClosed = true

	return
}

// LoadModules loads all the built in modules.
func (c *Client) LoadModules() (result error) {
	for _, module := range modules {
		err := c.LoadModule(module)
		if err != nil {
			Log.Error("error loading module", zap.String("module", module.Name), zap.Error(err))
			multierr.Append(result, err)
		}
	}

	return
}

// LoadModule loads a Module.
func (c *Client) LoadModule(m Module) error {
	for name, command := range m.Commands {
		if _, ok := c.Commands[name]; ok {
			return errors.New("Command '" + name + "' already exists")
		}

		c.Commands[name] = command

		// register aliases
		for _, alias := range command.Aliases {
			if _, ok := c.Commands[alias]; ok {
				return errors.New("Command '" + name + "' already exists")
			}

			c.Commands[alias] = command
		}
	}

	return nil
}

// ErrorHandler recovers from panics and reports them when deferred.
func (c *Client) ErrorHandler(scope string, handlers ...func(error)) {
	err := recover()

	switch rval := err.(type) {
	case nil:
		return
	case error:
		Log.Error("error in "+scope, zap.Error(rval))

		packet := raven.NewPacket(rval.Error(), raven.NewException(rval, raven.NewStacktrace(3, 5, appPkgPrefixes)))
		raven.DefaultClient.Capture(packet, map[string]string{
			"scope": scope,
		})

		for _, handler := range handlers {
			handler(rval)
		}
	default:
		rvalStr := fmt.Sprint(rval) // stringify
		Log.Error("error in "+scope, zap.String("error", rvalStr))

		packet := raven.NewPacket(rvalStr, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(3, 5, appPkgPrefixes)))
		raven.DefaultClient.Capture(packet, map[string]string{
			"scope": scope,
		})

		for _, handler := range handlers {
			handler(errors.New(rvalStr))
		}
	}
}
