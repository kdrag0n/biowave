package core

import (
	"sync/atomic"
	"github.com/kdrag0n/discordgo"
	"go.uber.org/zap"
	"github.com/dgraph-io/badger"
	"strings"
	"runtime/debug"
)

func (c *Client) onMessage(session *discordgo.Session, event *discordgo.MessageCreate) {
	defer c.ErrorHandler("message handler")

	if event.Author.ID == session.State.User.ID || len(event.Content) == 1 {
		return
	}

	channel, err := session.State.Channel(event.ChannelID)
	if err != nil {
		panic(err)
	}

	prefix, err := c.GetByID("prefix", channel.GuildID)
	if err == badger.ErrKeyNotFound {
		prefix = c.Config.DefaultPrefix
		go c.SetByID("prefix", channel.GuildID, prefix)
	} else if err != nil {
		Log.Error("error getting prefix, aborting handler", zap.Error(err))
		return
	}

	if strings.HasPrefix(event.Content, prefix) {
		split := strings.Fields(event.Content)
		commandName := strings.ToLower(split[0][len(prefix):])

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
				defer c.ErrorHandler("command", func(err error) {
					context.Fail("Error: " + err.Error())
				})

				debug.SetPanicOnFault(true)

				command.Function(context)
			}()
		}
	} else if (strings.HasPrefix(event.Content, c.ourMention) || strings.HasPrefix(event.Content, c.ourGuildMention)) && !event.MentionEveryone {
		request := strings.TrimSpace(event.Content[min(len(event.Content), 22):])
		if strings.EqualFold(request, "prefix") {

		}
	} // else if session.State.Channel(channelID uint64)
}

func (c *Client) onReady(session *discordgo.Session, event *discordgo.Ready) {
	defer c.ErrorHandler("ready handler")

	if atomic.LoadUint32(&c.isReady) != 1 {
		atomic.StoreUint32(&c.isReady, 1)

		// first to be ready
		sID := StrID(session.State.User.ID)
		c.ourMention = "<@" + sID + ">"
		c.ourGuildMention = "<@!" + sID + ">"

		app, err := session.Application(0)
		if err != nil {
			Log.Error("error getting bot application", zap.Error(err))
			c.ownerID = UserOriginalOwner
		} else {
			c.ownerID = app.Owner.ID
		}

		c.Developers.SetBit(c.ownerID)
	}

	_, err := session.State.Guild(GuildPrivate)
	if err == nil {
		c.EmoteOk = "<:ok:428754249027944458>"
		c.EmoteFail = "<:fail:428754276777459712>"
		c.EmoteBot = "<:bot:428754293156216834>"
		c.EmoteGrave = "<:rip:337405147347025930>"
	}
}
