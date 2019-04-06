package value

import "errors"

// Errors that may occur inside of behaviors.
var (
	ErrTypeMismatch     = errors.New("type mismatch")
	ErrInvalidField     = errors.New("invalid field")
	ErrInvalidIndexType = errors.New("invalid index type")
)

// Orderer is an interface for types that can be ordered.
type Orderer interface {
	Value
	GreaterThan(right Orderer) (Value, error)
	LessThan(right Orderer) (Value, error)
	GreaterThanOrEqual(right Orderer) (Value, error)
	LessThanOrEqual(right Orderer) (Value, error)
}

// Comparer is an interface for types that can be compared.
type Comparer interface {
	Value
	Equal(right Comparer) (Value, error)
	NotEqual(right Comparer) (Value, error)
}

// Bitwise is an interface for types that can be operated on bitwise.
type Bitwise interface {
	Value
	BitwiseNot() (Bitwise, error)
	BitwiseAnd(right Bitwise) (Bitwise, error)
	BitwiseXor(right Bitwise) (Bitwise, error)
	BitwiseOr(right Bitwise) (Bitwise, error)
	LeftShift(right Bitwise) (Bitwise, error)
	RightShift(right Bitwise) (Bitwise, error)
}

// Logical is an interface for types that can be operated on with logical operations.
type Logical interface {
	Value
	LogicalNot() (Value, error)
	LogicalAnd(right Logical) (Value, error)
	LogicalOr(right Logical) (Value, error)
}

// Adder is an interface for types that can handle addition.
type Adder interface {
	Value
	Negate() (Value, error)
	Add(right Adder) (Value, error)
	Sub(right Adder) (Value, error)
}

// Multiplier is an interface for types that can handle multiplication.
type Multiplier interface {
	Value
	Multiply(right Multiplier) (Value, error)
	Divide(right Multiplier) (Value, error)
}

// Moduler is an interface for types that can handle modulo.
type Moduler interface {
	Value
	Modulo(right Moduler) (Value, error)
}

// Caller is an interface for types that can be called.
type Caller interface {
	Value
	Call(args []Value) (Value, error)
}

// Descender is an interface for types that can be descended.
type Descender interface {
	Value
	Descend(member string) (Value, error)
}

// Indexer is an interface for types that can be indexed.
type Indexer interface {
	Value
	Index(index Value) (Value, error)
}
