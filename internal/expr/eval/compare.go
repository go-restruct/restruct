package eval

import (
	"fmt"

	"github.com/go-restruct/restruct/internal/expr/ast"
	"github.com/go-restruct/restruct/internal/expr/value"
)

func evaluateOperandsComparer(context Context, left ast.Node, right ast.Node) (value.Comparer, value.Comparer, error) {
	l, err := Evaluate(context, left)
	if err != nil {
		return nil, nil, err
	}
	bl, ok := l.(value.Comparer)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	r, err := Evaluate(context, right)
	if err != nil {
		return nil, nil, err
	}
	br, ok := r.(value.Comparer)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	return bl, br, nil
}

func evaluateEqual(context Context, expr ast.EqualExpression) (value.Value, error) {
	left, right, err := evaluateOperandsComparer(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.Equal(right)
}

func evaluateNotEqual(context Context, expr ast.NotEqualExpression) (value.Value, error) {
	left, right, err := evaluateOperandsComparer(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.NotEqual(right)
}

func evaluateConditional(context Context, expr ast.ConditionalExpression) (value.Value, error) {
	condition, err := Evaluate(context, expr.Condition)
	if err != nil {
		return nil, err
	}

	result, ok := condition.Value().(bool)
	if !ok {
		return nil, fmt.Errorf("expected bool, got %t from %s", condition.Value(), condition)
	}

	if result {
		return Evaluate(context, expr.Then)
	}
	return Evaluate(context, expr.Else)
}
