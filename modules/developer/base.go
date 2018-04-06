package developer

import (
	"github.com/kdrag0n/biowave/core"
)

func init() {
	core.RegisterModule("Developer", C{})
}

// C contains the module's commands.
type C struct{}
