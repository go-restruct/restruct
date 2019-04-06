package eval

import (
	"github.com/go-restruct/restruct/internal/expr/ast"
	"github.com/go-restruct/restruct/internal/expr/value"
)

func evaluateOperandsBitwise(context Context, left ast.Node, right ast.Node) (value.Bitwise, value.Bitwise, error) {
	l, err := Evaluate(context, left)
	if err != nil {
		return nil, nil, err
	}
	bl, ok := l.(value.Bitwise)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	r, err := Evaluate(context, right)
	if err != nil {
		return nil, nil, err
	}
	br, ok := r.(value.Bitwise)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	return bl, br, nil
}

func evaluateBitwiseNot(context Context, expr ast.BitwiseNotExpression) (value.Value, error) {
	val, err := Evaluate(context, expr.Operand)
	if err != nil {
		return nil, err
	}
	if b, ok := val.(value.Bitwise); ok {
		return b.BitwiseNot()
	}
	return nil, ErrInvalidType
}

func evaluateBitwiseLeftShift(context Context, expr ast.BitwiseLeftShiftExpression) (value.Value, error) {
	left, right, err := evaluateOperandsBitwise(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.LeftShift(right)
}

func evaluateBitwiseRightShift(context Context, expr ast.BitwiseRightShiftExpression) (value.Value, error) {
	left, right, err := evaluateOperandsBitwise(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.RightShift(right)
}

func evaluateBitwiseAnd(context Context, expr ast.BitwiseAndExpression) (value.Value, error) {
	left, right, err := evaluateOperandsBitwise(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.BitwiseAnd(right)
}

func evaluateBitwiseClear(context Context, expr ast.BitwiseClearExpression) (value.Value, error) {
	left, right, err := evaluateOperandsBitwise(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	right, err = right.BitwiseNot()
	if err != nil {
		return nil, err
	}
	return left.BitwiseAnd(right)
}

func evaluateBitwiseXor(context Context, expr ast.BitwiseXorExpression) (value.Value, error) {
	left, right, err := evaluateOperandsBitwise(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.BitwiseXor(right)
}

func evaluateBitwiseOr(context Context, expr ast.BitwiseOrExpression) (value.Value, error) {
	left, right, err := evaluateOperandsBitwise(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.BitwiseOr(right)
}
