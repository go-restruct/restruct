package value

import (
	"reflect"

	"github.com/go-restruct/restruct/internal/expr/typing"
)

var (
	_ = Value(Map{})
	_ = Indexer(Map{})
)

// Map represents a map type.
type Map struct {
	value reflect.Value
}

// NewMap creates a new map value.
func NewMap(value interface{}) Struct {
	rval := reflect.ValueOf(value)
	if rval.Kind() != reflect.Map {
		panic("NewMap called on non-map")
	}
	return Struct{rval}
}

func (c Map) String() string {
	return "<map>"
}

// Value implements Value
func (c Map) Value() interface{} { return c.value }

// Type implements Value
func (c Map) Type() (typing.Type, error) {
	typ, err := typing.FromReflectType(c.value.Type())
	if err != nil {
		return nil, err
	}
	return typ, nil
}

// Index implements Indexer
func (c Map) Index(index Value) (Value, error) {
	return FromValue(c.value.MapIndex(reflect.ValueOf(index.Value())))
}
