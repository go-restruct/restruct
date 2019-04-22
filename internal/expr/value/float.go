package value

import (
	"fmt"
	"strconv"

	"github.com/go-restruct/restruct/internal/expr/typing"
)

var (
	_ = Value(Float{})
	_ = Comparer(Float{})
	_ = Orderer(Float{})
	_ = Adder(Float{})
	_ = Multiplier(Float{})
)

// Float represents a floating point constant.
type Float struct {
	value float64
}

// NewFloat creates a floating point constant.
func NewFloat(value float64) Float {
	return Float{value}
}

// ParseFloat parses a floating point literal.
func ParseFloat(literal string) (Float, error) {
	val, err := strconv.ParseFloat(literal, 64)
	return NewFloat(val), err
}

func (c Float) String() string { return fmt.Sprintf("%f", c.value) }

// Value implements Value
func (c Float) Value() interface{} { return c.value }

// Type implements Value
func (c Float) Type() (typing.Type, error) {
	return typing.PrimitiveType(typing.Float, typing.TypeInfo{}), nil
}

// Equal implements Comparer
func (c Float) Equal(right Comparer) (Value, error) {
	r, err := asFloat(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value == r.value), nil
}

// NotEqual implements Comparer
func (c Float) NotEqual(right Comparer) (Value, error) {
	r, err := asFloat(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value != r.value), nil
}

// GreaterThan implements Orderer
func (c Float) GreaterThan(right Orderer) (Value, error) {
	r, err := asFloat(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value > r.value), nil
}

// LessThan implements Orderer
func (c Float) LessThan(right Orderer) (Value, error) {
	r, err := asFloat(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value < r.value), nil
}

// GreaterThanOrEqual implements Orderer
func (c Float) GreaterThanOrEqual(right Orderer) (Value, error) {
	r, err := asFloat(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value >= r.value), nil
}

// LessThanOrEqual implements Orderer
func (c Float) LessThanOrEqual(right Orderer) (Value, error) {
	r, err := asFloat(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value <= r.value), nil
}

// Negate implements Adder
func (c Float) Negate() (Value, error) {
	return NewFloat(-c.value), nil
}

// Add implements Adder
func (c Float) Add(right Adder) (Value, error) {
	r, err := asFloat(right)
	if err != nil {
		return nil, err
	}
	return NewFloat(c.value + r.value), nil
}

// Sub implements Adder
func (c Float) Sub(right Adder) (Value, error) {
	r, err := asFloat(right)
	if err != nil {
		return nil, err
	}
	return NewFloat(c.value - r.value), nil
}

// Multiply implements Multiplier
func (c Float) Multiply(right Multiplier) (Value, error) {
	r, err := asFloat(right)
	if err != nil {
		return nil, err
	}
	return NewFloat(c.value * r.value), nil
}

// Divide implements Multiplier
func (c Float) Divide(right Multiplier) (Value, error) {
	r, err := asFloat(right)
	if err != nil {
		return nil, err
	}
	return NewFloat(c.value / r.value), nil
}
