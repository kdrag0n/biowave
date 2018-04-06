package moderation

import (
	"github.com/kdrag0n/biowave/core"
)

func init() {
	core.RegisterModule("Moderation", C{})
}

// C contains the module's commands.
type C struct{}
