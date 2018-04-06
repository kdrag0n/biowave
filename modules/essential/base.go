package essential

import (
	"github.com/kdrag0n/biowave/core"
)

func init() {
	core.RegisterModule("Essential", C{})
}

// C contains the module's commands.
type C struct{}
