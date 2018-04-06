package core

import (
	"fmt"
	"errors"
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

	lastSent *discordgo.Message

	info *cInfo
}

// Send sends a message to the requesting channel.
func (c *Context) Send(message string) {
	msg, err := c.Session.ChannelMessageSend(c.Event.ChannelID, Truncate(FilterMessage(message)))
	if err != nil {
		panic(err)
	}

	c.lastSent = msg
}

// SendFmt sends a formatted message to the requesting channel.
func (c *Context) SendFmt(message string, obj ...interface{}) {
	c.Send(fmt.Sprintf(message, obj...))
}

// Ok sends a message to the requesting channel with a prefix indicating success.
func (c *Context) Ok(message string) {
	c.Send(c.Client.EmoteOk + " " + message)
}

// Fail sends a message to the requesting channel with a prefix indicating failure.
func (c *Context) Fail(message string) {
	c.Send(c.Client.EmoteFail + " " + message)
}

// Loading sends a message to the requesting channel with a prefix indicating loading.
func (c *Context) Loading(message string) {
	c.Send(c.Client.EmoteLoading + " " + message)
}

// OkFmt sends a formatted message to the requesting channel with a prefix indicating success.
func (c *Context) OkFmt(message string, obj ...interface{}) {
	c.SendFmt(c.Client.EmoteOk + " " + message, obj...)
}

// FailFmt sends a formatted message to the requesting channel with a prefix indicating failure.
func (c *Context) FailFmt(message string, obj ...interface{}) {
	c.SendFmt(c.Client.EmoteFail + " " + message, obj...)
}

// LoadingFmt sends a formatted message to the requesting channel with a prefix indicating loading.
func (c *Context) LoadingFmt(message string, obj ...interface{}) {
	c.SendFmt(c.Client.EmoteLoading + " " + message, obj...)
}

// React adds a reaction to the message that invoked the command.
func (c *Context) React(emote string) {
	err := c.Session.MessageReactionAdd(c.Event.ChannelID, c.Event.ID, emote)
	if err != nil {
		panic(err)
	}
}

// Edit edits the last sent message.
func (c *Context) Edit(message string) {
	if c.lastSent == nil {
		panic(errors.New("No message sent"))
	}

	msg, err := c.Session.ChannelMessageEdit(c.Event.ChannelID, c.lastSent.ID, message)
	if err != nil {
		panic(err)
	}

	c.lastSent = msg
}

// EditFmt edits the last sent message with a formatted message.
func (c *Context) EditFmt(format string, obj ...interface{}) {
	c.Edit(fmt.Sprintf(format, obj...))
}

// EditOk edits the last sent message with a prefix indicating success.
func (c *Context) EditOk(message string) {
	c.Edit(c.Client.EmoteOk + " " + message)
}

// EditFail edits the last sent message with a prefix indicating failure.
func (c *Context) EditFail(message string) {
	c.Edit(c.Client.EmoteFail + " " + message)
}

// EditLoading edits the last sent message with a prefix indicating loading.
func (c *Context) EditLoading(message string) {
	c.Edit(c.Client.EmoteLoading + " " + message)
}

// EditOkFmt edits the last sent message, formatted, with a prefix indicating success.
func (c *Context) EditOkFmt(message string, obj ...interface{}) {
	c.EditFmt(c.Client.EmoteOk + " " + message, obj...)
}

// EditFailFmt edits the last sent message, formatted, with a prefix indicating failure.
func (c *Context) EditFailFmt(message string, obj ...interface{}) {
	c.EditFmt(c.Client.EmoteFail + " " + message, obj...)
}

// EditLoadingFmt edits the last sent message, formatted, with a prefix indicating loading.
func (c *Context) EditLoadingFmt(message string, obj ...interface{}) {
	c.EditFmt(c.Client.EmoteLoading + " " + message, obj...)
}

// Info is used to set command information.
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

// Embed returns a new embed builder.
func (c *Context) Embed() *Embed {
	return &Embed{
		MessageEmbed: &discordgo.MessageEmbed{},
		context: c,
	}
}

// SendEmbed sends an embed as a message.
func (c *Context) SendEmbed(embed *Embed) {
	msg, err := c.Session.ChannelMessageSendEmbed(c.Event.ChannelID, embed.MessageEmbed)
	if err != nil {
		panic(err)
	}

	c.lastSent = msg
}

// DirectSendEmbed DMs an embed to the requesting user.
func (c *Context) DirectSendEmbed(embed *Embed) {
	channel, err := c.Session.UserChannelCreate(c.Event.Author.ID)
	if err != nil {
		panic(err)
	}

	msg, err := c.Session.ChannelMessageSendEmbed(channel.ID, embed.MessageEmbed)
	if err != nil {
		panic(err)
	}

	c.lastSent = msg
}
