package compiler

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"

	"github.com/go-restruct/restruct/internal/expr/ast"
	"github.com/go-restruct/restruct/internal/expr/typing"
	"github.com/go-restruct/restruct/internal/expr/value"
)

type genLocation int

const (
	genReturn genLocation = iota
	genTopLevel
	genSubExpr
)

type compilerContext struct {
	TypingContext  typing.Context
	TempVarCount   int
	CondStatements []jen.Code
}

func genType(typ typing.Type) *jen.Statement {
	// For named types.
	name := typ.Name()
	if name != "" {
		pkgPath := typ.PkgPath()
		if pkgPath == "" {
			return jen.Id(name)
		}
		return jen.Qual(pkgPath, name)
	}

	switch typ.Kind() {
	case typing.Boolean:
		return jen.Bool()
	case typing.Int:
		return jen.Int64()
	case typing.Uint:
		return jen.Uint64()
	case typing.Float:
		return jen.Float64()
	case typing.String:
		return jen.String()
	case typing.Array:
		elem, err := typ.Elem()
		if err != nil {
			panic(err)
		}
		return jen.Index().Add(genType(elem))
	case typing.Func:
		v, err := typ.IsVariadic()
		if err != nil {
			panic(err)
		}
		n, err := typ.NumParams()
		if err != nil {
			panic(err)
		}
		ret, err := typ.Return()
		if err != nil {
			panic(err)
		}

		return jen.Func().ParamsFunc(func(g *jen.Group) {
			for i := 0; i < n; i++ {
				param, err := typ.Param(i)
				if err != nil {
					panic(err)
				}
				if i == n-1 && v {
					g.Lit("...")
				}
				g.Add(genType(param))
			}
		}).Add(genType(ret))
	case typing.Map:
		key, err := typ.Key()
		if err != nil {
			panic(err)
		}
		elem, err := typ.Elem()
		if err != nil {
			panic(err)
		}
		return jen.Map(genType(key)).Add(genType(elem))
	case typing.Struct:
		n, err := typ.NumFields()
		if err != nil {
			panic(err)
		}
		return jen.StructFunc(func(g *jen.Group) {
			for i := 0; i < n; i++ {
				f, err := typ.Field(i)
				if err != nil {
					panic(err)
				}
				g.Id(f.Name).Add(genType(f.Type))
			}
		})
	default:
		panic("type literal not implemented")
	}
}

func genExprList(context *compilerContext, list ast.ExpressionList) *jen.Statement {
	var exprs []jen.Code
	for _, expr := range list.Nodes {
		exprs = append(exprs, genExpr(context, expr, genTopLevel))
	}
	return jen.Call(exprs...)
}

func genExpr(context *compilerContext, expr ast.Node, loc genLocation) *jen.Statement {
	switch loc {
	case genReturn:
		// Turn the given expression into a return.
		switch t := expr.(type) {
		case ast.ParenExpression:
			return genExpr(context, t.Node, genReturn)
		case ast.ConditionalExpression:
			condExpr := genExpr(context, t.Condition, genTopLevel)
			thenExpr := genExpr(context, t.Then, genReturn)
			elseExpr := genExpr(context, t.Else, genReturn)
			return jen.If(condExpr).Block(thenExpr).Else().Block(elseExpr)
		default:
			if loc == genReturn {
				return jen.Return(genExpr(context, expr, genSubExpr))
			}
		}
	case genTopLevel:
		// Perform optimizations that are allowed at top level.
		switch t := expr.(type) {
		case ast.ParenExpression:
			return genExpr(context, t.Node, genTopLevel)
		}
	}

	switch t := expr.(type) {
	case ast.ParenExpression:
		return jen.Parens(genExpr(context, t.Node, genSubExpr))
	case ast.IdentifierExpression:
		if t.Name[0] == strings.ToUpper(t.Name)[0] {
			return jen.Id("self").Dot(t.Name)
		}
		return jen.Id(t.Name)
	case ast.FunctionCallExpression:
		return genExpr(context, t.Function, genSubExpr).Add(genExprList(context, t.Arguments))
	case ast.IndexExpression:
		return genExpr(context, t.Operand, genSubExpr).Index(genExpr(context, t.Index, genSubExpr))
	case ast.DotExpression:
		return genExpr(context, t.Left, genSubExpr).Dot(t.Right.Name)
	case ast.NegateExpression:
		return jen.Op("-").Add(genExpr(context, t.Operand, genSubExpr))
	case ast.LogicalNotExpression:
		return jen.Op("!").Add(genExpr(context, t.Operand, genSubExpr))
	case ast.BitwiseNotExpression:
		return jen.Op("^").Add(genExpr(context, t.Operand, genSubExpr))
	case ast.MultiplyExpression:
		return genExpr(context, t.Left, genSubExpr).Op("*").Add(genExpr(context, t.Right, genSubExpr))
	case ast.DivideExpression:
		return genExpr(context, t.Left, genSubExpr).Op("/").Add(genExpr(context, t.Right, genSubExpr))
	case ast.ModuloExpression:
		return genExpr(context, t.Left, genSubExpr).Op("%").Add(genExpr(context, t.Right, genSubExpr))
	case ast.AddExpression:
		return genExpr(context, t.Left, genSubExpr).Op("+").Add(genExpr(context, t.Right, genSubExpr))
	case ast.SubtractExpression:
		return genExpr(context, t.Left, genSubExpr).Op("-").Add(genExpr(context, t.Right, genSubExpr))
	case ast.BitwiseLeftShiftExpression:
		return genExpr(context, t.Left, genSubExpr).Op("<<").Add(genExpr(context, t.Right, genSubExpr))
	case ast.BitwiseRightShiftExpression:
		return genExpr(context, t.Left, genSubExpr).Op(">>").Add(genExpr(context, t.Right, genSubExpr))
	case ast.GreaterThanExpression:
		return genExpr(context, t.Left, genSubExpr).Op(">").Add(genExpr(context, t.Right, genSubExpr))
	case ast.LessThanExpression:
		return genExpr(context, t.Left, genSubExpr).Op("<").Add(genExpr(context, t.Right, genSubExpr))
	case ast.GreaterThanOrEqualExpression:
		return genExpr(context, t.Left, genSubExpr).Op(">=").Add(genExpr(context, t.Right, genSubExpr))
	case ast.LessThanOrEqualExpression:
		return genExpr(context, t.Left, genSubExpr).Op("<=").Add(genExpr(context, t.Right, genSubExpr))
	case ast.EqualExpression:
		return genExpr(context, t.Left, genSubExpr).Op("==").Add(genExpr(context, t.Right, genSubExpr))
	case ast.NotEqualExpression:
		return genExpr(context, t.Left, genSubExpr).Op("!=").Add(genExpr(context, t.Right, genSubExpr))
	case ast.BitwiseAndExpression:
		return genExpr(context, t.Left, genSubExpr).Op("&").Add(genExpr(context, t.Right, genSubExpr))
	case ast.BitwiseClearExpression:
		return genExpr(context, t.Left, genSubExpr).Op("&^").Add(genExpr(context, t.Right, genSubExpr))
	case ast.BitwiseXorExpression:
		return genExpr(context, t.Left, genSubExpr).Op("^").Add(genExpr(context, t.Right, genSubExpr))
	case ast.BitwiseOrExpression:
		return genExpr(context, t.Left, genSubExpr).Op("|").Add(genExpr(context, t.Right, genSubExpr))
	case ast.LogicalAndExpression:
		return genExpr(context, t.Left, genSubExpr).Op("&&").Add(genExpr(context, t.Right, genSubExpr))
	case ast.LogicalOrExpression:
		return genExpr(context, t.Left, genSubExpr).Op("||").Add(genExpr(context, t.Right, genSubExpr))
	case ast.ConditionalExpression:
		context.TempVarCount++
		tempName := fmt.Sprintf("temp%d", context.TempVarCount)

		typ, err := t.Type(context.TypingContext)
		if err != nil {
			panic(err)
		}

		condExpr := genExpr(context, t.Condition, genTopLevel)
		thenExpr := jen.Id(tempName).Op("=").Add(genExpr(context, t.Then, genTopLevel))
		elseExpr := jen.Id(tempName).Op("=").Add(genExpr(context, t.Else, genTopLevel))

		context.CondStatements = append(context.CondStatements, jen.Var().Id(tempName).Add(genType(typ)).Line().If(condExpr).Block(thenExpr).Else().Block(elseExpr))

		return jen.Id(tempName)
	case ast.Constant:
		switch v := t.Value.(type) {
		case value.ConstInt:
			return jen.Lit(int(v.Value().(int64)))
		case value.Boolean, value.Float, value.Int, value.String, value.Uint:
			return jen.Lit(v.Value())
		default:
			panic(fmt.Errorf("generating type literal for type not supported: %v", t.Value.Value()))
		}
	default:
		panic(fmt.Errorf("unknown node type %t", t))
	}
}

// Compile compiles an expression to Go code.
func Compile(functionName string, typctx typing.Context, node ast.Node) (*jen.Statement, error) {
	context := &compilerContext{
		TypingContext:  typctx,
		TempVarCount:   0,
		CondStatements: []jen.Code{},
	}

	typ, err := node.Type(typctx)
	if err != nil {
		return nil, err
	}

	expr := genExpr(context, node, genReturn)

	statements := []jen.Code{}
	statements = append(statements, context.CondStatements...)
	statements = append(statements, expr)
	return jen.Func().Id(functionName).Params(jen.Id("self").Add(genType(typctx.Self))).Add(genType(typ)).Block(statements...), nil
}
