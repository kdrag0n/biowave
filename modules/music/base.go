package music

import (
	"github.com/kdrag0n/biowave/core"
)

func init() {
	core.RegisterModule("Music", C{})
}

// C contains the module's commands.
type C struct{}
