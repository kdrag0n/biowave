package core

import (
	"github.com/bwmarrin/discordgo"
)

// Context provides a context for commands.
type Context struct {
	Client  *Client
	Session *discordgo.Session
	Event   *discordgo.MessageCreate
	Invoker string
	Args    []string
	RawArgs string
}

// Send sends a message to the channel the request came from.
func (c *Context) Send(message string) {
	c.Session.ChannelMessageSend(c.Event.ChannelID, Truncate(FilterMessage(message)))
}

// Ok sends a message to the requesting channel with a prefix indicating success.
func (c *Context) Ok(message string) {
	c.Send(c.Client.EmoteOk + " " + message)
}

// Fail sends a message to the requesting channel with a prefix indicating failure.
func (c *Context) Fail(message string) {
	c.Send(c.Client.EmoteFail + " " + message)
}
