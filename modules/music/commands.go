package music

import (
	"github.com/kdrag0n/biowave/core"
)

// Play plays a requested song.
func (C) Play(c *core.Context) {
	c.Info("music")
}
