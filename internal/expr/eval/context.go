package eval

import "reflect"

// Context contains an execution context for expressions.
type Context struct {
	// Self contains a reference to the current struct. If it is nil,
	// execution will fail on non-constant expressions.
	Self reflect.Value

	// Global contains any global variables or functions.
	Global map[string]interface{}
}

// NewContext creates a new context in the default state.
func NewContext(s interface{}) Context {
	var v reflect.Value
	if s != nil {
		v = reflect.ValueOf(s)
	}
	return Context{
		Self:   v,
		Global: builtins,
	}
}
