package eval

import (
	"github.com/go-restruct/restruct/internal/expr/ast"
	"github.com/go-restruct/restruct/internal/expr/value"
)

func evaluateOperandsAdder(context Context, left ast.Node, right ast.Node) (value.Adder, value.Adder, error) {
	l, err := Evaluate(context, left)
	if err != nil {
		return nil, nil, err
	}
	bl, ok := l.(value.Adder)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	r, err := Evaluate(context, right)
	if err != nil {
		return nil, nil, err
	}
	br, ok := r.(value.Adder)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	return bl, br, nil
}

func evaluateNegate(context Context, expr ast.NegateExpression) (value.Value, error) {
	val, err := Evaluate(context, expr.Operand)
	if err != nil {
		return nil, err
	}
	a, ok := val.(value.Adder)
	if !ok {
		return nil, ErrInvalidType
	}
	return a.Negate()
}

func evaluateAdd(context Context, expr ast.AddExpression) (value.Value, error) {
	left, right, err := evaluateOperandsAdder(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.Add(right)
}

func evaluateSubtract(context Context, expr ast.SubtractExpression) (value.Value, error) {
	left, right, err := evaluateOperandsAdder(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.Sub(right)
}
