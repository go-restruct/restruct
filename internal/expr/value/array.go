package value

import "reflect"

var (
	_ = Value(Array{})
	_ = Indexer(Array{})
)

// Array represents an array-like type.
type Array struct {
	value reflect.Value
}

// NewArray creates a new array value.
func NewArray(value interface{}) Struct {
	rval := reflect.ValueOf(value)
	if rval.Kind() != reflect.Array && rval.Kind() != reflect.Slice {
		panic("NewArray called on incompatible type")
	}
	return Struct{rval}
}

func (c Array) String() string {
	return "<array>"
}

// Value implements Value
func (c Array) Value() interface{} { return c.value }

// Index implements Indexer
func (c Array) Index(index Value) (Value, error) {
	switch t := index.Value().(type) {
	case uint64:
		return FromValue(c.value.Index(int(t)).Interface())
	case int64:
		return FromValue(c.value.Index(int(t)).Interface())
	default:
		return nil, ErrInvalidIndexType
	}
}
