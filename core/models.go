package core

// Command is a chat command usable by users.
type Command struct {
	Name        string
	Description string
	Aliases     []string
	Permissions []Permission
	Function    CommandFunc
}

// CommandFunc is the internal function for a command.
type CommandFunc func(*Context)
