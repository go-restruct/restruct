package value

import (
	"errors"
	"reflect"
)

var (
	_ = Value(Func{})
	_ = Caller(Func{})
)

// Func represents a function.
type Func struct {
	value reflect.Value
}

// NewFunc creates a new struct value.
func NewFunc(value interface{}) Func {
	rval := reflect.ValueOf(value)
	if rval.Kind() != reflect.Func {
		panic("NewFunc called on non-function")
	}
	return Func{rval}
}

func (c Func) String() string {
	return "<func>"
}

// Value implements Value
func (c Func) Value() interface{} { return c.value }

// Call implements Caller
func (c Func) Call(args []Value) (Value, error) {
	vals := make([]reflect.Value, len(args))
	for i := range args {
		vals[i] = reflect.ValueOf(args[i].Value())
	}

	retvals := c.value.Call(vals)
	if len(retvals) != 1 {
		return nil, errors.New("")
	}

	return FromValue(retvals[0].Interface())
}
