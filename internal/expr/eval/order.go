package eval

import (
	"github.com/go-restruct/restruct/internal/expr/ast"
	"github.com/go-restruct/restruct/internal/expr/value"
)

func evaluateOperandsOrder(context Context, left ast.Node, right ast.Node) (value.Orderer, value.Orderer, error) {
	l, err := Evaluate(context, left)
	if err != nil {
		return nil, nil, err
	}
	bl, ok := l.(value.Orderer)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	r, err := Evaluate(context, right)
	if err != nil {
		return nil, nil, err
	}
	br, ok := r.(value.Orderer)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	return bl, br, nil
}

func evaluateGreaterThan(context Context, expr ast.GreaterThanExpression) (value.Value, error) {
	left, right, err := evaluateOperandsOrder(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.GreaterThan(right)
}

func evaluateLessThan(context Context, expr ast.LessThanExpression) (value.Value, error) {
	left, right, err := evaluateOperandsOrder(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.LessThan(right)
}

func evaluateGreaterThanOrEqual(context Context, expr ast.GreaterThanOrEqualExpression) (value.Value, error) {
	left, right, err := evaluateOperandsOrder(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.GreaterThanOrEqual(right)
}

func evaluateLessThanOrEqual(context Context, expr ast.LessThanOrEqualExpression) (value.Value, error) {
	left, right, err := evaluateOperandsOrder(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.LessThanOrEqual(right)
}
