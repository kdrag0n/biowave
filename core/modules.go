package core

import (
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"unsafe"
)

var modules = make(map[string]Module, 10)

// A Module contains commands.
type Module struct {
	Name     string
	Commands map[string]*Command

	// Standard actions
	Cleanup ModuleCleanupFunc
	Unload ModuleUnloadFunc
}

// A ModuleCleanupFunc cleans up a module.
type ModuleCleanupFunc func() error

// A ModuleUnloadFunc unloads a module from a Client.
type ModuleUnloadFunc func(c *Client) error

// Add registers commands in a Module.
func (m *Module) Add(name, desc string, aliases []string, usage []Argument, hidden, guildOnly bool, exec CommandFunc) {
	if len(desc) == 0 {
		desc = "No description."
	}

	if aliases == nil {
		aliases = []string{}
	}

	if usage == nil {
		usage = []Argument{}
	}

	m.Commands[name] = &Command{
		Name:        name,
		Description: desc,
		Aliases:     aliases,
		Usage:       usage,
		Hidden:      hidden,
		GuildOnly:   guildOnly,
		Permissions: []Permission{},
		Function:    exec,
	}
}

// RegisterModule registers a module to be loaded into Clients.
func RegisterModule(name string, cmdStruct interface{}) {
	module := &Module{
		Name:     name,
		Commands: make(map[string]*Command, 20),

		Cleanup: func() error {return nil},
		Unload: func (*Client) error {return nil},
	}

	t := reflect.TypeOf(cmdStruct)
	for f := 0; f < t.NumMethod(); f++ {
		funk := t.Method(f)
		funcName := strings.ToLower(funk.Name)

		funcValue := funk.Func
		funcAddressable := reflect.New(funcValue.Type()).Elem()
		funcAddressable.Set(funcValue)
		funcPtr := unsafe.Pointer(funcAddressable.UnsafeAddr()) // #nosec

		if funcName == "mcleanup" {
			module.Cleanup = *(*ModuleCleanupFunc)(funcPtr)
			continue
		} else if funcName == "munload" {
			module.Unload = *(*ModuleUnloadFunc)(funcPtr)
			continue
		}

		exec := *(*CommandFunc)(funcPtr)

		ctx := &Context{
			Args: nil,
			info: &cInfo{
				name:      funcName,
				desc:      "No description.",
				aliases:   nil,
				usage:     nil,
				hidden:    false,
				guildOnly: false,
				perms:     nil,
			},
		}

		// get information
		func() {
			defer func() { recover() }()
			exec(ctx)
		}()
		info := ctx.info

		module.Add(info.name, info.desc, info.aliases, info.usage, info.hidden, info.guildOnly, exec)
	}

	modules[name] = *module
}

// ModuleCleanup cleans up after modules. Call before shutting down the program.
func ModuleCleanup() (result error) {
	for _, module := range modules {
		err := module.Cleanup()
		if err != nil {
			Log.Error("Error on module cleanup", zap.String("module", module.Name), zap.Error(err))
			multierr.Append(result, err)
		}
	}

	return
}
