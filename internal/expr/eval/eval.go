package eval

import (
	"errors"
	"fmt"

	"github.com/go-restruct/restruct/internal/expr/ast"
	"github.com/go-restruct/restruct/internal/expr/value"
)

// Evaluation errors.
var (
	ErrInvalidType = errors.New("invalid type")
)

func typeMismatch(a, b value.Value, expr ast.Node) error {
	return fmt.Errorf("type mismatch: %t != %t in expression %s", a.Value(), b.Value(), expr.Source())
}

func evaluateOperands(context Context, left ast.Node, right ast.Node) (value.Value, value.Value, error) {
	l, err := Evaluate(context, left)
	if err != nil {
		return nil, nil, err
	}

	r, err := Evaluate(context, right)
	if err != nil {
		return nil, nil, err
	}

	return l, r, nil
}

// Evaluate evaluates an expression with regards to context.
func Evaluate(context Context, expr ast.Node) (value.Value, error) {
	switch t := expr.(type) {
	case ast.ParenExpression:
		return Evaluate(context, t.Node)
	case ast.IdentifierExpression:
		return evaluateIdentifier(context, t)
	case ast.FunctionCallExpression:
		return evaluateFunctionCall(context, t)
	case ast.IndexExpression:
		return evaluateIndex(context, t)
	case ast.DotExpression:
		return evaluateDot(context, t)
	case ast.NegateExpression:
		return evaluateNegate(context, t)
	case ast.LogicalNotExpression:
		return evaluateLogicalNot(context, t)
	case ast.BitwiseNotExpression:
		return evaluateBitwiseNot(context, t)
	case ast.MultiplyExpression:
		return evaluateMultiply(context, t)
	case ast.DivideExpression:
		return evaluateDivide(context, t)
	case ast.ModuloExpression:
		return evaluateModulo(context, t)
	case ast.AddExpression:
		return evaluateAdd(context, t)
	case ast.SubtractExpression:
		return evaluateSubtract(context, t)
	case ast.BitwiseLeftShiftExpression:
		return evaluateBitwiseLeftShift(context, t)
	case ast.BitwiseRightShiftExpression:
		return evaluateBitwiseRightShift(context, t)
	case ast.GreaterThanExpression:
		return evaluateGreaterThan(context, t)
	case ast.LessThanExpression:
		return evaluateLessThan(context, t)
	case ast.GreaterThanOrEqualExpression:
		return evaluateGreaterThanOrEqual(context, t)
	case ast.LessThanOrEqualExpression:
		return evaluateLessThanOrEqual(context, t)
	case ast.EqualExpression:
		return evaluateEqual(context, t)
	case ast.NotEqualExpression:
		return evaluateNotEqual(context, t)
	case ast.BitwiseAndExpression:
		return evaluateBitwiseAnd(context, t)
	case ast.BitwiseClearExpression:
		return evaluateBitwiseClear(context, t)
	case ast.BitwiseXorExpression:
		return evaluateBitwiseXor(context, t)
	case ast.BitwiseOrExpression:
		return evaluateBitwiseOr(context, t)
	case ast.LogicalAndExpression:
		return evaluateLogicalAnd(context, t)
	case ast.LogicalOrExpression:
		return evaluateLogicalOr(context, t)
	case ast.ConditionalExpression:
		return evaluateConditional(context, t)
	case ast.Constant:
		return t.Value, nil
	default:
		panic(fmt.Errorf("unknown node type %t", t))
	}
}
