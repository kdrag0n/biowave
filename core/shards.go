package core

import (
	"github.com/kdrag0n/discordgo"
)

// UserCount returns the total number of users across all shards.
func (c *Client) UserCount() (total int) {
	c.ForSessions(func(session *discordgo.Session) {
		for _, guild := range session.State.Guilds {
			total += len(guild.Members) // too expensive to fully dedupe users
		}
	})

	return
}

// ChannelCount returns the total number of channels across all shards.
func (c *Client) ChannelCount() (total int) {
	c.ForSessions(func(session *discordgo.Session) {
		total += len(session.State.PrivateChannels)
		for _, guild := range session.State.Guilds {
			total += len(guild.Channels)
		}
	})

	return
}

// GuildCount returns the total number of guilds across all shards.
func (c *Client) GuildCount() (total int) {
	c.ForSessions(func(session *discordgo.Session) {
		total += len(session.State.Guilds)
	})

	return
}
