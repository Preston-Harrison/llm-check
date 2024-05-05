package common

import (
	"fmt"
	"go/token"
)

type Function interface {
	Module() string
	Name() string
	String() string
	// Newline delimited function doc comments.
	Docs() []string
	CallSites() []Function
	Location() (token.Position, token.Position)
}

type Check struct {
	Content string
	Scope   string
}

func (c Check) String() string {
	return fmt.Sprintf("Check{ Scope: \"%s\", Content: \"%s\" }", c.Scope, c.Content)
}

type FunctionWithCheck interface {
	Function
	Check() Check
}
