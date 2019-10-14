package restruct

import (
	"reflect"

	"github.com/go-restruct/restruct/expr"
)

type structstack struct {
	stack     []reflect.Value
	exprenv   expr.Resolver
	allowexpr bool
}

func (s *structstack) resolver() expr.Resolver {
	if s.exprenv == nil {
		if !s.allowexpr {
			panic("call restruct.EnableExprBeta() to eanble expressions beta")
		}
		s.exprenv = makeResolver(s.ancestor(0))
	}
	return s.exprenv
}

func (s *structstack) push(v reflect.Value) {
	s.stack = append(s.stack, v)
	s.exprenv = nil
}

func (s *structstack) pop(v reflect.Value) {
	var p reflect.Value
	s.stack, p = s.stack[:len(s.stack)-1], s.stack[len(s.stack)-1]
	if p != v {
		panic("struct stack misaligned")
	}
	s.exprenv = nil
}

func (s *structstack) setancestor(f field, v reflect.Value, ancestor reflect.Value) {
	if ancestor.CanAddr() && ancestor.Kind() != reflect.Ptr {
		ancestor = ancestor.Addr()
	}
	if ancestor.Type().AssignableTo(v.Type()) {
		v.Set(ancestor)
	}
}

func (s *structstack) root() reflect.Value {
	if len(s.stack) > 0 {
		return s.stack[0]
	}
	return reflect.ValueOf(nil)
}

func (s *structstack) ancestor(generation int) reflect.Value {
	if len(s.stack) > generation {
		return s.stack[len(s.stack)-generation-1]
	}
	return reflect.ValueOf(nil)
}
