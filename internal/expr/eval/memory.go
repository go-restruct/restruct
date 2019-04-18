package eval

import (
	"fmt"

	"github.com/go-restruct/restruct/internal/expr/ast"
	"github.com/go-restruct/restruct/internal/expr/value"
)

func evaluateExpressionList(context Context, exprs ast.ExpressionList) ([]value.Value, error) {
	values := []value.Value{}

	for _, arg := range exprs.Nodes {
		c, err := Evaluate(context, arg)
		if err != nil {
			return nil, err
		}
		values = append(values, c)
	}

	return values, nil
}

func evaluateIdentifier(context Context, expr ast.IdentifierExpression) (value.Value, error) {
	return context.Resolve(expr.Name)
}

func evaluateFunctionCall(context Context, expr ast.FunctionCallExpression) (value.Value, error) {
	fn, err := Evaluate(context, expr.Function)
	if err != nil {
		return nil, err
	}

	if c, ok := fn.(value.Caller); ok {
		arguments, err := evaluateExpressionList(context, expr.Arguments)
		if err != nil {
			return nil, err
		}
		return c.Call(arguments)
	}

	return nil, ErrInvalidType
}

func evaluateIndex(context Context, expr ast.IndexExpression) (value.Value, error) {
	operand, err := Evaluate(context, expr.Operand)
	if err != nil {
		return nil, err
	}
	if i, ok := operand.(value.Indexer); ok {
		index, err := Evaluate(context, expr.Index)
		if err != nil {
			return nil, err
		}
		return i.Index(index)
	}
	return nil, ErrInvalidType
}

func evaluateDot(context Context, expr ast.DotExpression) (value.Value, error) {
	left, err := Evaluate(context, expr.Left)
	if err != nil {
		return nil, err
	}

	if d, ok := left.(value.Descender); ok {
		return d.Descend(expr.Right.Name)
	}

	return nil, fmt.Errorf("cannot descend %t", left)
}
