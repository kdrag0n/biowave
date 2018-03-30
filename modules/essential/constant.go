package essential

import (
	"github.com/kdrag0n/biowave/core"
)

func init() {
	core.RegisterModule("Essential", C{})
}

type C struct{}

func (C) Test(c *core.Context) {
	c.Info("Test the bot.")

	c.Session.MessageReactionAdd(c.Event.ChannelID, c.Event.ID, "👍")
}
