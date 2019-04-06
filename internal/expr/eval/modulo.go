package eval

import (
	"github.com/go-restruct/restruct/internal/expr/ast"
	"github.com/go-restruct/restruct/internal/expr/value"
)

func evaluateOperandsModuler(context Context, left ast.Node, right ast.Node) (value.Moduler, value.Moduler, error) {
	l, err := Evaluate(context, left)
	if err != nil {
		return nil, nil, err
	}
	bl, ok := l.(value.Moduler)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	r, err := Evaluate(context, right)
	if err != nil {
		return nil, nil, err
	}
	br, ok := r.(value.Moduler)
	if !ok {
		return nil, nil, ErrInvalidType
	}

	return bl, br, nil
}

func evaluateModulo(context Context, expr ast.ModuloExpression) (value.Value, error) {
	left, right, err := evaluateOperandsModuler(context, expr.Left, expr.Right)
	if err != nil {
		return nil, err
	}
	return left.Modulo(right)
}
