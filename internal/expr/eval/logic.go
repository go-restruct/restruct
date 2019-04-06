package eval

import (
	"github.com/go-restruct/restruct/internal/expr/ast"
	"github.com/go-restruct/restruct/internal/expr/value"
)

func evaluateOperandsLogical(context Context, left ast.Node, right ast.Node) (value.Logical, value.Logical, error) {
	l, err := Evaluate(context, left)
	if err != nil {
		return nil, nil, err
	}
	bl, ok := l.(value.Logical)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	r, err := Evaluate(context, right)
	if err != nil {
		return nil, nil, err
	}
	br, ok := r.(value.Logical)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	return bl, br, nil
}

func evaluateLogicalNot(context Context, expr ast.LogicalNotExpression) (value.Value, error) {
	val, err := Evaluate(context, expr.Operand)
	if err != nil {
		return nil, err
	}

	l, ok := val.(value.Logical)
	if !ok {
		return nil, ErrInvalidType
	}

	return l.LogicalNot()
}

func evaluateLogicalAnd(context Context, expr ast.LogicalAndExpression) (value.Value, error) {
	left, right, err := evaluateOperandsLogical(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.LogicalAnd(right)
}

func evaluateLogicalOr(context Context, expr ast.LogicalOrExpression) (value.Value, error) {
	left, right, err := evaluateOperandsLogical(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.LogicalOr(right)
}
