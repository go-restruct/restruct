package value

import "strconv"

var (
	_ = Value(Int{})
	_ = Comparer(Int{})
	_ = Orderer(Int{})
	_ = Bitwise(Int{})
	_ = Adder(Int{})
	_ = Multiplier(Int{})
	_ = Moduler(Int{})
)

// Int represents a signed integer constant.
type Int struct {
	value int64
}

// NewInt creates a new signed integer constant.
func NewInt(value int64) Int {
	return Int{value}
}

func (c Int) String() string { return strconv.FormatInt(c.value, 10) }

// Value implements Value
func (c Int) Value() interface{} { return c.value }

// Equal implements Comparer
func (c Int) Equal(right Comparer) (Value, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value == r.value), nil
}

// NotEqual implements Comparer
func (c Int) NotEqual(right Comparer) (Value, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value != r.value), nil
}

// GreaterThan implements Orderer
func (c Int) GreaterThan(right Orderer) (Value, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value > r.value), nil
}

// LessThan implements Orderer
func (c Int) LessThan(right Orderer) (Value, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value < r.value), nil
}

// GreaterThanOrEqual implements Orderer
func (c Int) GreaterThanOrEqual(right Orderer) (Value, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value >= r.value), nil
}

// LessThanOrEqual implements Orderer
func (c Int) LessThanOrEqual(right Orderer) (Value, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewBoolean(c.value <= r.value), nil
}

// BitwiseNot implements Bitwise
func (c Int) BitwiseNot() (Bitwise, error) {
	return NewInt(^c.value), nil
}

// BitwiseAnd implements Bitwise
func (c Int) BitwiseAnd(right Bitwise) (Bitwise, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewInt(c.value & r.value), nil
}

// BitwiseXor implements Bitwise
func (c Int) BitwiseXor(right Bitwise) (Bitwise, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewInt(c.value ^ r.value), nil
}

// BitwiseOr implements Bitwise
func (c Int) BitwiseOr(right Bitwise) (Bitwise, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewInt(c.value | r.value), nil
}

// LeftShift implements Bitwise
func (c Int) LeftShift(right Bitwise) (Bitwise, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewInt(c.value << r.value), nil
}

// RightShift implements Bitwise
func (c Int) RightShift(right Bitwise) (Bitwise, error) {
	r, err := asUint(right)
	if err != nil {
		return nil, err
	}
	return NewInt(c.value >> r.value), nil
}

// Negate implements Adder
func (c Int) Negate() (Value, error) {
	return NewInt(-c.value), nil
}

// Add implements Adder
func (c Int) Add(right Adder) (Value, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewInt(c.value + r.value), nil
}

// Sub implements Adder
func (c Int) Sub(right Adder) (Value, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewInt(c.value - r.value), nil
}

// Multiply implements Multiplier
func (c Int) Multiply(right Multiplier) (Value, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewInt(c.value * r.value), nil
}

// Divide implements Multiplier
func (c Int) Divide(right Multiplier) (Value, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewInt(c.value / r.value), nil
}

// Modulo implements Moduler
func (c Int) Modulo(right Moduler) (Value, error) {
	r, err := asInt(right)
	if err != nil {
		return nil, err
	}
	return NewInt(c.value % r.value), nil
}
