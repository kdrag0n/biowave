package core

import (
	"reflect"
)

// An Argument describes the purpose of a command argument.
type Argument struct {
	// Name stores the full name of the argument.
	Name string
	// Inline stores a short, lower-case description of the argument.
	Inline string
	// Type stores the intended type of the argument.
	Type reflect.Type
	// Optional stores whether the argument is optional.
	Optional bool
}
