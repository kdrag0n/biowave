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
