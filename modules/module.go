package modules

import (
	"github.com/kdrag0n/biowave/core"
)

// A Module contains commands
type Module struct {
	Name     string
	Commands map[string]core.Command
}

// Add registers commands in a Module
func (m *Module) Add(commands ...core.Command) {

}
