package essential

import (
	"github.com/kdrag0n/biowave/core"
)

func init() {
	core.RegisterModule("Essential", C{})
}

// C contains commands.
type C struct{}

// Test makes sure the bot is working.
func (C) Test(c *core.Context) {
	c.Info("Test the bot.")

	c.Session.MessageReactionAdd(c.Event.ChannelID, c.Event.ID, "👍")
}

// Owner explains the role of bot owner to the user.
func (C) Owner(c *core.Context) {
	c.Info("Become bot owner.") // Lure them in

	c.Session.ChannelMessageSend(c.Event.ChannelID, `My owner is **kdragon#1337**. It's something that affects entire bot, and is the person who actually owns the bot, not the owner of a server.
**No**, you may not have bot owner because it allows full control. In your server, being server owner is sufficient.`)
}
