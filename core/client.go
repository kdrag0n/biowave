package core

import (
	"fmt"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/raven-go"
)

// Client is connected to the platform
type Client struct {
	Config Config
	SentryTags map[string]string
	Sessions []*discordgo.Session
}

// NewClient creates a new Discord client
func NewClient(config Config) *Client {
	sessions := make([]*discordgo.Session, config.Shards)

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
	}
}

// ForSessions executes a function with every platform session
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
	}

	return nil
}

// LoadModules loads all the built in modules
func (c *Client) LoadModules() error {
	return nil
}

// ErrorHandler recovers from panics and reports them when deferred
func (c *Client) ErrorHandler() {
	var packet *raven.Packet
	err := recover()

	switch rval := err.(type) {
	case nil:
		return
	case error:
		packet = raven.NewPacket(rval.Error(), raven.NewException(rval, raven.NewStacktrace(2, 3, nil)))
	default:
		rvalStr := fmt.Sprint(rval)

		packet = raven.NewPacket(rvalStr, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)))
	}

	_, _ = raven.DefaultClient.Capture(packet, c.SentryTags)
}
