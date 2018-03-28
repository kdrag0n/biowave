package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kdrag0n/cyborg/platform"
)

// Adapter is a Discord platform adapter
type Adapter struct {
	dgSessions []*discordgo.Session
}

// NewAdapter creates a new Discord platform adapter with the options in config
func NewAdapter(config platform.Config) *Adapter {
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

	return &Adapter{
		dgSessions: sessions,
	}
}

// Connect connects the adapter to Discord
func (a *Adapter) Connect() error {
	for _, dg := range a.dgSessions {
		err := dg.Open()
		if err != nil {
			return err
		}
	}

	return nil
}
