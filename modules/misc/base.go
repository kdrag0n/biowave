package misc

import (
	"github.com/kdrag0n/biowave/core"
)

func init() {
	core.RegisterModule("Miscellaneous", C{})
}

// C contains the module's commands.
type C struct{}
