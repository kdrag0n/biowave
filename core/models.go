package core

import (
	"time"
)


// Format is a message format
type Format uint8

const (
	// FormatMarkdown represents the Markdown formatting syntax
	FormatMarkdown Format = iota
	// FormatPlain represents plain text
	FormatPlain
)

// Context provides a context for commands
type Context struct {
	Client *Client
}

// User represents an user on the chat platform
type User struct {
	Name         string
	DisplayName  string
	ID           uint64
	Mention      string
	CreationTime time.Time
}

// Member represents an user wrapped with group-specific information
type Member struct {
	Nickname string
	User     *User
}

// Message represents a message in chat
type Message struct {
	Text         string
	Format       Format
	Sender       *User
	CreationTime time.Time
}

// Command is a chat command usable by users  
type Command struct {
	Name string
	Description string
	Aliases []string
	Permissions []Permission
	Function CommandFunc
}

// CommandFunc is the internal function for a command
type CommandFunc func(*Context)
