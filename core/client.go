package core

import (
	"sync/atomic"
	"go.uber.org/zap"
	"errors"
	"fmt"
	"strings"

	"github.com/kdrag0n/discordgo"
	"github.com/dgraph-io/badger"
	"github.com/getsentry/raven-go"
)

// Client is the full client context
type Client struct {
	Config     Config
	SentryTags map[string]string
	Sessions   []*discordgo.Session
	Commands   map[string]*Command

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
func NewClient(config Config) *Client {
	sessions := make([]*discordgo.Session, config.Shards)

	// create DiscordGo sessions for shards
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

		EmoteOk:    "✅",
		EmoteFail:  "❌",
		EmoteGrave: "⚰",
		EmoteBot:   "bot",

		ourMention:      "<@0>",
		ourGuildMention: "<@!0>",
		ownerID:         0,

		isReady: 0,
	}
}

// ForSessions executes a function with every session
func (c *Client) ForSessions(iter func(*discordgo.Session)) {
	for _, dg := range c.Sessions {
		iter(dg)
	}
}

// Start initiates the client's connections
func (c *Client) Start() error {
	c.LoadModules()
	for _, dg := range c.Sessions {
		err := dg.Open()
		if err != nil {
			return err
		}

		dg.AddHandler(c.OnMessage)
		dg.AddHandler(func(session *discordgo.Session, event *discordgo.Ready) {
			defer c.ErrorHandler("ready handler")
			
			if atomic.LoadUint32(&c.isReady) != 1 {
				atomic.StoreUint32(&c.isReady, 1)

				// first to be ready
				sID := idToStr(dg.State.User.ID)
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
				c.EmoteOk = "<:ok:428754235799240708>"
				// c.EmoteOk = "<:ok2:428754249027944458>"
				c.EmoteFail = "<:fail:428754265486655498>"
				// c.EmoteFail = "<:fail2:428754276777459712>"
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

	return nil
}

// Stop stops the client and all associated Discord sessions.
func (c *Client) Stop() {
	for _, dg := range c.Sessions {
		dg.Close()
	}
	// TODO: unload modules
}

// OnMessage handles an incoming message.
func (c *Client) OnMessage(session *discordgo.Session, event *discordgo.MessageCreate) {
	defer c.ErrorHandler("message handler")

	if event.Author.ID == session.State.User.ID {
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
			}() // TODO: pool
		}
	} else if strings.HasPrefix(event.Content, c.ourMention) || strings.HasPrefix(event.Content, c.ourGuildMention) {
		request := strings.TrimSpace(event.Content[min(len(event.Content), 22):])
		if strings.EqualFold(request, "prefix") {

		}
	} // TODO: private channel
}

// LoadModules loads all the built in modules
func (c *Client) LoadModules() error {
	for _, module := range modules {
		err := c.LoadModule(module)
		if err != nil {
			return err
		}
	}

	return nil
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
		Log.Error("Error in " + scope, zap.Error(rval))

		packet := raven.NewPacket(rval.Error(), raven.NewException(rval, raven.NewStacktrace(2, 3, nil)))
		raven.DefaultClient.Capture(packet, c.SentryTags)
	default:
		rvalStr := fmt.Sprint(rval) // stringify
		Log.Error("Error in " + scope, zap.String("error", rvalStr))

		packet := raven.NewPacket(rvalStr, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)))
		raven.DefaultClient.Capture(packet, c.SentryTags)
	}
}
