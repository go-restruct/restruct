package value

import "strconv"

var (
	_ = Value(ConstInt{})
	_ = Comparer(ConstInt{})
	_ = Orderer(ConstInt{})
	_ = Bitwise(ConstInt{})
	_ = Adder(ConstInt{})
	_ = Multiplier(ConstInt{})
	_ = Moduler(ConstInt{})
)

// ConstInt represents a signed integer constant.
type ConstInt struct {
	Int
}

// NewConstInt creates a new signed integer constant.
func NewConstInt(value int64) ConstInt {
	return ConstInt{Int{value}}
}

// ParseInt parses an integer literal.
func ParseInt(literal string) (ConstInt, error) {
	val, err := strconv.ParseInt(literal, 0, 64)
	return NewConstInt(val), err
}

func asUint(v Value) (Uint, error) {
	switch t := v.(type) {
	case Uint:
		return t, nil
	case ConstInt:
		return NewUint(uint64(t.Int.value)), nil
	default:
		return NewUint(0), ErrTypeMismatch
	}
}

func asInt(v Value) (Int, error) {
	switch t := v.(type) {
	case Int:
		return t, nil
	case ConstInt:
		return NewInt(int64(t.Int.value)), nil
	default:
		return NewInt(0), ErrTypeMismatch
	}
}

func asFloat(v Value) (Float, error) {
	switch t := v.(type) {
	case Float:
		return t, nil
	case ConstInt:
		return NewFloat(float64(t.Int.value)), nil
	default:
		return NewFloat(0), ErrTypeMismatch
	}
}
