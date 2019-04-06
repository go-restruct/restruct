package eval

import (
	"github.com/go-restruct/restruct/internal/expr/ast"
	"github.com/go-restruct/restruct/internal/expr/value"
)

func evaluateOperandsMultiplier(context Context, left ast.Node, right ast.Node) (value.Multiplier, value.Multiplier, error) {
	l, err := Evaluate(context, left)
	if err != nil {
		return nil, nil, err
	}
	bl, ok := l.(value.Multiplier)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	r, err := Evaluate(context, right)
	if err != nil {
		return nil, nil, err
	}
	br, ok := r.(value.Multiplier)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	return bl, br, nil
}

func evaluateMultiply(context Context, expr ast.MultiplyExpression) (value.Value, error) {
	left, right, err := evaluateOperandsMultiplier(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.Multiply(right)
}

func evaluateDivide(context Context, expr ast.DivideExpression) (value.Value, error) {
	left, right, err := evaluateOperandsMultiplier(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.Divide(right)
}
