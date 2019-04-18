package value

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/go-restruct/restruct/internal/expr/typing"
)

var (
	_ = Value(String{})
	_ = Indexer(String{})
	_ = Comparer(String{})
	_ = Orderer(String{})
)

// String represents a string constant.
type String struct {
	value string
}

// NewString creates a new string constant.
func NewString(value string) String {
	return String{value}
}

// ParseStrLiteral parses a quoted string literal.
func ParseStrLiteral(literal string) (expr String, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	buf := bytes.Buffer{}
	r := strings.NewReader(literal)

	ch, _, err := r.ReadRune()
	if err != nil {
		return String{}, err
	} else if ch != '"' {
		return String{}, errors.New("syntax error: expected '")
	}

	for {
		ch, err := readStrLitRune(r, '"')
		if err == errEndOfString {
			break
		} else if err != nil {
			return String{}, err
		}
		buf.WriteRune(ch)

	}

	return NewString(buf.String()), nil
}

func (c String) String() string { return fmt.Sprintf("%q", c.value) }

// Value implements Value
func (c String) Value() interface{} { return c.value }

// Type implements Value
func (c String) Type() (typing.Type, error) { return typing.PrimitiveType(typing.String), nil }

// Index implements Indexer
func (c String) Index(index Value) (Value, error) {
	switch t := index.Value().(type) {
	case uint64:
		return NewUint(uint64(c.value[int(t)])), nil
	case int64:
		return NewUint(uint64(c.value[int(t)])), nil
	default:
		return nil, ErrInvalidIndexType
	}
}

// Equal implements Comparer
func (c String) Equal(right Comparer) (Value, error) {
	r, ok := right.(String)
	if !ok {
		return nil, ErrTypeMismatch
	}
	return NewBoolean(c.value == r.value), nil
}

// NotEqual implements Comparer
func (c String) NotEqual(right Comparer) (Value, error) {
	r, ok := right.(String)
	if !ok {
		return nil, ErrTypeMismatch
	}
	return NewBoolean(c.value != r.value), nil
}

// GreaterThan implements Orderer
func (c String) GreaterThan(right Orderer) (Value, error) {
	r, ok := right.(String)
	if !ok {
		return nil, ErrTypeMismatch
	}
	return NewBoolean(c.value > r.value), nil
}

// LessThan implements Orderer
func (c String) LessThan(right Orderer) (Value, error) {
	r, ok := right.(String)
	if !ok {
		return nil, ErrTypeMismatch
	}
	return NewBoolean(c.value < r.value), nil
}

// GreaterThanOrEqual implements Orderer
func (c String) GreaterThanOrEqual(right Orderer) (Value, error) {
	r, ok := right.(String)
	if !ok {
		return nil, ErrTypeMismatch
	}
	return NewBoolean(c.value >= r.value), nil
}

// LessThanOrEqual implements Orderer
func (c String) LessThanOrEqual(right Orderer) (Value, error) {
	r, ok := right.(String)
	if !ok {
		return nil, ErrTypeMismatch
	}
	return NewBoolean(c.value <= r.value), nil
}
