package restruct

import (
	"reflect"

	"github.com/go-restruct/restruct/expr"
)

var (
	expressionsEnabled = false
)

// EnableExprBeta enables you to use restruct expr while it is still in beta.
// Use at your own risk. Functionality may change in unforeseen, incompatible
// ways at any time.
func EnableExprBeta() {
	expressionsEnabled = true
}

func makeResolver(s reflect.Value) expr.Resolver {
	env := expr.NewMetaResolver()
	env.AddResolver(expr.NewStructResolver(s))
	env.AddResolver(expr.NewMapResolver(exprStdLib))
	return env
}
