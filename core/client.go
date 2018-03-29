package core

import (
	"strings"
	"fmt"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/raven-go"
)

// Client is the full client context
type Client struct {
	Config Config
	SentryTags map[string]string
	Sessions []*discordgo.Session
	Commands map[string]*Command
	Emotes map[string]string

	ourMention string
	ourGuildMention string
	ownerID string
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
		dg.StateEnabled = false
		dg.SyncEvents = true
		dg.MaxRestRetries = 3

		sessions[id] = dg
	}

	return &Client{
		Config: config,
		SentryTags: nil,
		Sessions: sessions,
		Commands: make(map[string]*Command, 120),
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
	for _, dg := range c.Sessions {
		err := dg.Open()
		if err != nil {
			return err
		}

		dg.AddHandler(c.OnMessage)
	}

	return nil
}

// OnMessage handles an incoming message.
func (c *Client) OnMessage(session *discordgo.Session, event *discordgo.MessageCreate) {
	prefix := c.Config.DefaultPrefix // TODO: prefix

	if strings.HasPrefix(event.Content, prefix) {
		split := strings.Fields(event.Content)
		commandName := split[0][len(prefix):]

		if command, ok := c.Commands[commandName]; ok {
			context := &Context{
				Client: c,
				Session: session,
				Event: event,
				Invoker: commandName,
				Args: split[1:],
				RawArgs: strings.TrimSpace(event.Content[len(prefix)+len(commandName):]),
			}

			command.Function(context)
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
func (c *Client) ErrorHandler() {
	err := recover()

	switch rval := err.(type) {
	case nil:
		return
	case error:
		packet := raven.NewPacket(rval.Error(), raven.NewException(rval, raven.NewStacktrace(2, 3, nil)))
		raven.DefaultClient.Capture(packet, c.SentryTags)
	default:
		rvalStr := fmt.Sprint(rval) // stringify
		packet := raven.NewPacket(rvalStr, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)))
		raven.DefaultClient.Capture(packet, c.SentryTags)
	}
}
