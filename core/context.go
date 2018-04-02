package core

import (
	"github.com/kdrag0n/discordgo"
)

// Context provides a context for commands.
type Context struct {
	Client  *Client
	Session *discordgo.Session
	Event   *discordgo.MessageCreate
	Invoker string
	Args    []string
	RawArgs string

	info *cInfo
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

// React adds a reaction to the message that invoked the command.
func (c *Context) React(emote string) {
	c.Session.MessageReactionAdd(c.Event.ChannelID, c.Event.ID, emote)
}

// Info is used to obtain command information.
func (c *Context) Info(description string, cont ...struct{}) *ChainInfo {
	if c.info == nil {
		return nil
	}

	c.info.desc = description

	if len(cont) == 0 {
		panic(0)
	}

	return &ChainInfo{
		i: c.info,
	}
}
