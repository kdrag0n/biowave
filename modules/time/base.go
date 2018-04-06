package time

import (
	"github.com/kdrag0n/biowave/core"
)

func init() {
	core.RegisterModule("Time", C{})
}

// C contains the module's commands.
type C struct{}
