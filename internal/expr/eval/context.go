package eval

import (
	"github.com/go-restruct/restruct/internal/expr/value"
)

// Context contains an execution context for expressions.
type Context struct {
	// Self contains a reference to the current struct. If it is nil,
	// execution will fail on non-constant expressions.
	Self value.Struct

	// Global contains any global variables or functions.
	Global map[string]value.Value
}

// NewContext creates a new context in the default state.
func NewContext(s interface{}) (Context, error) {
	self := value.NewStruct(struct{}{})
	if s != nil {
		self = value.NewStruct(s)
	}
	return Context{
		Self:   self,
		Global: builtins,
	}, nil
}

// Resolve resolves an identifier contextually.
func (context Context) Resolve(ident string) (value.Value, error) {
	if val, ok := context.Global[ident]; ok {
		return val, nil
	}

	return context.Self.Descend(ident)
}
