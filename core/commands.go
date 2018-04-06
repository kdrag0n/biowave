package core

// Command is a chat command usable by users.
type Command struct {
	Name        string
	Description string
	Aliases     []string
	Usage       []Argument
	Hidden      bool
	GuildOnly   bool

	Module *Module
	Permissions []Permission
	Function    CommandFunc
}

// CommandFunc is the internal function for a command.
type CommandFunc func(*Context)

// cInfo temporarily holds command information.
type cInfo struct {
	name      string
	desc      string
	aliases   []string
	usage     []Argument
	hidden    bool
	guildOnly bool
	perms     []Permission
}

// ChainInfo provides an easy chained interface to set all command information.
type ChainInfo struct {
	i *cInfo
}

// Requires returns whether the command requires a permission.
func (c *Command) Requires(reqPerm Permission) bool {
	for _, perm := range c.Permissions {
		if perm == reqPerm {
			return true
		}
	}

	return false
}

// Aliases sets the aliases for a command.
func (c *ChainInfo) Aliases(aliases ...string) *ChainInfo {
	c.i.aliases = aliases
	return c
}

// Usage sets the argument usage for a command.
func (c *ChainInfo) Usage(usage []Argument) *ChainInfo {
	c.i.usage = usage
	return c
}

// Hidden sets a command as hidden.
func (c *ChainInfo) Hidden() *ChainInfo {
	c.i.hidden = true
	return c
}

// GuildOnly makes a command only work in guilds.
func (c *ChainInfo) GuildOnly() *ChainInfo {
	c.i.guildOnly = true
	return c
}

// Perms sets permissions required to use command.
func (c *ChainInfo) Perms(perms ...Permission) *ChainInfo {
	c.i.perms = perms
	return c
}

// End ends the chained information sequence.
func (c *ChainInfo) End() {
	panic(0)
}
