package search

import (
	"github.com/kdrag0n/biowave/core"
)

func init() {
	core.RegisterModule("Search", C{})
}

// C contains the module's commands.
type C struct{}
