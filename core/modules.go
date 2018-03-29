package core

var modules = make(map[string]Module, 10)

// A Module contains commands.
type Module struct {
	Name     string
	Commands map[string]*Command
}

// A ModuleLoader loads a Module.
type ModuleLoader func(*Module)

// Add registers commands in a Module.
func (m *Module) Add(name string, desc string, aliases []string, exec CommandFunc) {
	if len(desc) == 0 {
		desc = "No description."
	}

	if aliases == nil {
		aliases = []string{}
	}

	m.Commands[name] = &Command{
		Name:        name,
		Description: desc,
		Aliases:     aliases,
		Function:    exec,
		Permissions: []Permission{},
	}
}

// RegisterModule registers a module to be loaded into Clients.
func RegisterModule(name string, loader ModuleLoader) {
	module := &Module{
		Name:     name,
		Commands: make(map[string]*Command, 20),
	}

	loader(module)
	modules[name] = *module
}
