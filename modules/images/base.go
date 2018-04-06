package images

import (
	"github.com/kdrag0n/biowave/core"
)

func init() {
	core.RegisterModule("Images", C{})
}

// C contains the module's commands.
type C struct{}
