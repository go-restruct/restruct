package expr

import (
	"github.com/go-restruct/restruct/internal/expr"
	"github.com/go-restruct/restruct/internal/expr/eval"
)

// Eval evaluates a restruct expr.
func Eval(s interface{}, input string) (interface{}, error) {
	expr, err := expr.ParseString(input)
	if err != nil {
		return nil, err
	}
	expr = expr.ConstantFold()

	context, err := eval.NewContext(s)
	if err != nil {
		return nil, err
	}

	result, err := eval.Evaluate(context, expr)
	if err != nil {
		return nil, err
	}

	return result.Value(), nil
}
