package value

import "reflect"

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

// Index implements Indexer
func (c Map) Index(index Value) (Value, error) {
	return FromValue(c.value.MapIndex(reflect.ValueOf(index.Value())))
}
