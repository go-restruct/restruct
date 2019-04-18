package value

import (
	"reflect"

	"github.com/go-restruct/restruct/internal/expr/typing"
)

var (
	_ = Value(Struct{})
	_ = Descender(Struct{})
	_ = Indexer(Struct{})
)

// Struct represents a struct.
type Struct struct {
	value reflect.Value
}

// NewStruct creates a new struct value.
func NewStruct(value interface{}) Struct {
	rval := reflect.ValueOf(value)
	for {
		if rval.Kind() != reflect.Ptr {
			break
		}
		rval = rval.Elem()
	}
	if rval.Kind() != reflect.Struct {
		panic("NewStruct called on non-struct")
	}
	return Struct{rval}
}

func (c Struct) String() string {
	return "<struct>"
}

// Value implements Value
func (c Struct) Value() interface{} { return c.value }

// Type implements Value
func (c Struct) Type() (typing.Type, error) {
	typ, err := typing.FromReflectType(c.value.Type())
	if err != nil {
		return nil, err
	}
	return typ, nil
}

// Descend implements Descender
func (c Struct) Descend(member string) (Value, error) {
	mval := c.value.FieldByName(member)
	if !mval.IsValid() {
		return nil, ErrInvalidField
	}
	return FromValue(mval.Interface())
}

// Index implements Indexer
func (c Struct) Index(index Value) (Value, error) {
	member, ok := index.Value().(string)
	if !ok {
		return nil, ErrInvalidIndexType
	}
	mval := c.value.FieldByName(member)
	if !mval.IsValid() {
		return nil, ErrInvalidField
	}
	return FromValue(mval.Interface())
}
