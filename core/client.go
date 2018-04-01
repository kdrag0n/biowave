package core

import (
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"strings"
	"sync/atomic"

	"github.com/dgraph-io/badger"
	"github.com/getsentry/raven-go"
	"github.com/kdrag0n/discordgo"
)

// Client is the full client context
type Client struct {
	Config     Config
	SentryTags map[string]string
	Sessions   []*discordgo.Session
	Commands   map[string]*Command
	
	// Data
	DB *badger.DB
	IsDBClosed bool

	// Emotes
	EmoteOk    string
	EmoteFail  string
	EmoteGrave string
	EmoteBot   string

	ourMention      string
	ourGuildMention string
	ownerID         uint64

	isReady uint32
}

// NewClient creates a new Discord client
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

	return &Client{
		Config:     config,
		SentryTags: nil,
		Sessions:   sessions,
		Commands:   make(map[string]*Command, 120),

		DB: db,
		IsDBClosed: false,

		EmoteOk:    "✅",
		EmoteFail:  "❌",
		EmoteGrave: "⚰",
		EmoteBot:   "bot",

		ourMention:      "<@0>",
		ourGuildMention: "<@!0>",
		ownerID:         0,

		isReady: 0,
	}, nil
}

// ForSessions executes a function with every session
func (c *Client) ForSessions(iter func(*discordgo.Session)) {
	for _, dg := range c.Sessions {
		iter(dg)
	}
}

// Start initiates the client's connections
func (c *Client) Start() (result error) {
	c.LoadModules()

	for idx, dg := range c.Sessions {
		err := dg.Open()
		if err != nil {
			Log.Error("Error opening session", zap.Int("shard", idx), zap.Error(err))
			multierr.Append(result, err)
		}

		dg.AddHandler(c.OnMessage)
		dg.AddHandler(func(session *discordgo.Session, event *discordgo.Ready) {
			defer c.ErrorHandler("ready handler")

			if atomic.LoadUint32(&c.isReady) != 1 {
				atomic.StoreUint32(&c.isReady, 1)

				// first to be ready
				sID := StrID(dg.State.User.ID)
				c.ourMention = "<@" + sID + ">"
				c.ourGuildMention = "<@!" + sID + ">"

				app, err := dg.Application(dg.State.User.ID)
				if err != nil {
					Log.Error("Error getting bot application", zap.Error(err))
					c.ownerID = UserOriginalOwner
				} else {
					c.ownerID = app.Owner.ID
				}
			}

			_, err = dg.State.Guild(GuildPrivate)
			if err == nil {
				c.EmoteOk = "<:ok:428754249027944458>"
				c.EmoteFail = "<:fail:428754276777459712>"
				c.EmoteBot = "<:bot:428754293156216834>"
				c.EmoteGrave = "<:rip:337405147347025930>"

				/*
					<a:loading:428754343018364929>
					<a:typing:428754324668022785>
					<a:loading2:428754355911524357>
					<a:download:428754309610733588>
				*/
			}
		})
	}

	return
}

// Stop stops the client and all associated Discord sessions.
func (c *Client) Stop() (result error) {
	// close sessions
	for idx, dg := range c.Sessions {
		err := dg.Close()
		if err != nil {
			Log.Error("Error closing session", zap.Int("shard", idx), zap.Error(err))
			multierr.Append(result, err)
		}
	}

	// unload modules
	for _, module := range modules {
		err := module.Unload(c)
		if err != nil {
			Log.Error("Error on module unload", zap.String("module", module.Name), zap.Error(err))
			multierr.Append(result, err)
		}
	}

	// close database
	err := c.DB.Close()
	if err != nil {
		Log.Error("Error closing database", zap.Error(err))
		multierr.Append(result, err)
	}
	c.IsDBClosed = true

	return
}

// OnMessage handles an incoming message.
func (c *Client) OnMessage(session *discordgo.Session, event *discordgo.MessageCreate) {
	defer c.ErrorHandler("message handler")

	if event.Author.ID == session.State.User.ID || len(event.Content) == 1 {
		return
	}

	prefix := c.Config.DefaultPrefix // TODO: prefix

	if strings.HasPrefix(event.Content, prefix) {
		split := strings.Fields(event.Content)
		commandName := split[0][len(prefix):]

		if command, ok := c.Commands[commandName]; ok {
			context := &Context{
				Client:  c,
				Session: session,
				Event:   event,
				Invoker: commandName,
				Args:    split[1:],
				RawArgs: strings.TrimSpace(event.Content[len(prefix)+len(commandName):]),
				info:    nil,
			}

			go func() {
				defer c.ErrorHandler("command")
				command.Function(context)
			}()
		}
	} else if strings.HasPrefix(event.Content, c.ourMention) || strings.HasPrefix(event.Content, c.ourGuildMention) {
		request := strings.TrimSpace(event.Content[min(len(event.Content), 22):])
		if strings.EqualFold(request, "prefix") {

		}
	} // else if session.State.Channel(channelID uint64)
}

// LoadModules loads all the built in modules
func (c *Client) LoadModules() (result error) {
	for _, module := range modules {
		err := c.LoadModule(module)
		if err != nil {
			Log.Error("Error loading module", zap.String("module", module.Name), zap.Error(err))
			multierr.Append(result, err)
		}
	}

	return
}

// LoadModule loads a Module
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

// ErrorHandler recovers from panics and reports them when deferred
func (c *Client) ErrorHandler(scope string) {
	err := recover()

	switch rval := err.(type) {
	case nil:
		return
	case error:
		Log.Error("Error in "+scope, zap.Error(rval))

		packet := raven.NewPacket(rval.Error(), raven.NewException(rval, raven.NewStacktrace(2, 3, nil)))
		raven.DefaultClient.Capture(packet, c.SentryTags)
	default:
		rvalStr := fmt.Sprint(rval) // stringify
		Log.Error("Error in "+scope, zap.String("error", rvalStr))

		packet := raven.NewPacket(rvalStr, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)))
		raven.DefaultClient.Capture(packet, c.SentryTags)
	}
}
