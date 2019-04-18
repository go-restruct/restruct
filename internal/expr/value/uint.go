package value

import (
	"strconv"

	"github.com/go-restruct/restruct/internal/expr/typing"
)

var (
	_ = Value(Uint{})
	_ = Comparer(Uint{})
	_ = Orderer(Uint{})
	_ = Bitwise(Uint{})
	_ = Adder(Uint{})
	_ = Multiplier(Uint{})
	_ = Moduler(Uint{})
)

// Uint represents an unsigned integer constant.
type Uint struct {
	value uint64
}

// NewUint creates a new signed integer constant.
func NewUint(value uint64) Uint {
	return Uint{value}
}

func (c Uint) String() string { return strconv.FormatUint(c.value, 10) }

// Value implements Value
func (c Uint) Value() interface{} { return c.value }

// Type implements Value
func (c Uint) Type() (typing.Type, error) { return typing.PrimitiveType(typing.Uint), nil }

// Equal implements Comparer
func (c Uint) Equal(right Comparer) (Value, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value == r.value), nil
}

// NotEqual implements Comparer
func (c Uint) NotEqual(right Comparer) (Value, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value != r.value), nil
}

// GreaterThan implements Orderer
func (c Uint) GreaterThan(right Orderer) (Value, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value > r.value), nil
}

// LessThan implements Orderer
func (c Uint) LessThan(right Orderer) (Value, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value < r.value), nil
}

// GreaterThanOrEqual implements Orderer
func (c Uint) GreaterThanOrEqual(right Orderer) (Value, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value >= r.value), nil
}

// LessThanOrEqual implements Orderer
func (c Uint) LessThanOrEqual(right Orderer) (Value, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value <= r.value), nil
}

// BitwiseNot implements Bitwise
func (c Uint) BitwiseNot() (Bitwise, error) {
	return NewUint(^c.value), nil
}

// BitwiseAnd implements Bitwise
func (c Uint) BitwiseAnd(right Bitwise) (Bitwise, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewUint(c.value & r.value), nil
}

// BitwiseXor implements Bitwise
func (c Uint) BitwiseXor(right Bitwise) (Bitwise, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewUint(c.value ^ r.value), nil
}

// BitwiseOr implements Bitwise
func (c Uint) BitwiseOr(right Bitwise) (Bitwise, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewUint(c.value | r.value), nil
}

// LeftShift implements Bitwise
func (c Uint) LeftShift(right Bitwise) (Bitwise, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewUint(c.value << r.value), nil
}

// RightShift implements Bitwise
func (c Uint) RightShift(right Bitwise) (Bitwise, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewUint(c.value >> r.value), nil
}

// Negate implements Adder
func (c Uint) Negate() (Value, error) {
	return NewInt(int64(-c.value)), nil
}

// Add implements Adder
func (c Uint) Add(right Adder) (Value, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewUint(c.value + r.value), nil
}

// Sub implements Adder
func (c Uint) Sub(right Adder) (Value, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewUint(c.value - r.value), nil
}

// Multiply implements Multiplier
func (c Uint) Multiply(right Multiplier) (Value, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewUint(c.value * r.value), nil
}

// Divide implements Multiplier
func (c Uint) Divide(right Multiplier) (Value, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewUint(c.value / r.value), nil
}

// Modulo implements Moduler
func (c Uint) Modulo(right Moduler) (Value, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewUint(c.value % r.value), nil
}
