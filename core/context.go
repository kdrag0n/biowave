package core

import (
	"github.com/bwmarrin/discordgo"
)

// Context provides a context for commands.
type Context struct {
	Client *Client
	Session *discordgo.Session
	Event *discordgo.MessageCreate
	Invoker string
	Args []string
	RawArgs string
}

func (c *Context) Send(message string) {
	c.Session.ChannelMessageSend(c.Event.ChannelID, message)
}

func (c *Context) Ok(message string) {
	c.Session.ChannelMessageSend(c.Event.ChannelID, "")
}
