package ast

import (
	"fmt"
	"strings"

	"github.com/go-restruct/restruct/internal/expr/value"
)

// This asserts that all nodes fulfill the node interface.
var (
	_ = Node(ExpressionList{})
	_ = Node(ParenExpression{})
	_ = Node(IdentifierExpression{})
	_ = Node(FunctionCallExpression{})
	_ = Node(IndexExpression{})
	_ = Node(DotExpression{})
	_ = Node(NegateExpression{})
	_ = Node(LogicalNotExpression{})
	_ = Node(BitwiseNotExpression{})
	_ = Node(MultiplyExpression{})
	_ = Node(DivideExpression{})
	_ = Node(ModuloExpression{})
	_ = Node(AddExpression{})
	_ = Node(SubtractExpression{})
	_ = Node(BitwiseLeftShiftExpression{})
	_ = Node(BitwiseRightShiftExpression{})
	_ = Node(GreaterThanExpression{})
	_ = Node(LessThanExpression{})
	_ = Node(GreaterThanOrEqualExpression{})
	_ = Node(LessThanOrEqualExpression{})
	_ = Node(EqualExpression{})
	_ = Node(NotEqualExpression{})
	_ = Node(BitwiseAndExpression{})
	_ = Node(BitwiseXorExpression{})
	_ = Node(BitwiseOrExpression{})
	_ = Node(LogicalAndExpression{})
	_ = Node(LogicalOrExpression{})
	_ = Node(ConditionalExpression{})
	_ = Node(Constant{})
)

// Node is the interface for any kind of node.
type Node interface {
	// Source returns source code for the given node.
	Source() string

	// ConstantFold returns an equivalent node with context folding performed recursively.
	ConstantFold() Node
}

// ExpressionList represents a list of expressions.
type ExpressionList struct {
	Nodes []Node
}

// NewEmptyExpressionList creates a new empty expression list.
func NewEmptyExpressionList() ExpressionList {
	return ExpressionList{[]Node{}}
}

// NewExpressionList creates a new expression list.
func NewExpressionList(node Node) ExpressionList {
	return ExpressionList{[]Node{node}}
}

// AppendExpression appends an expression to an expression list.
func AppendExpression(list ExpressionList, node Node) ExpressionList {
	return ExpressionList{append(list.Nodes, node)}
}

// Source implements Node
func (e ExpressionList) Source() string {
	sources := []string{}
	for _, node := range e.Nodes {
		sources = append(sources, node.Source())
	}
	return strings.Join(sources, ", ")
}

// ConstantFold implements Node
func (e ExpressionList) ConstantFold() Node {
	nodes := []Node{}
	for _, node := range e.Nodes {
		nodes = append(nodes, node.ConstantFold())
	}
	return ExpressionList{nodes}
}

// ParenExpression represents an expression in parenthesis.
type ParenExpression struct {
	Node Node
}

// NewParenExpression creates a new expression list.
func NewParenExpression(node Node) ParenExpression {
	return ParenExpression{node}
}

// Source implements Node
func (e ParenExpression) Source() string {
	return "(" + e.Node.Source() + ")"
}

// ConstantFold implements Node
func (e ParenExpression) ConstantFold() Node {
	return ParenExpression{e.Node.ConstantFold()}
}

// IdentifierExpression represents an identifier.
type IdentifierExpression struct {
	Name string
}

// NewIdentifierExpression creates a new identifier expression.
func NewIdentifierExpression(name string) IdentifierExpression {
	return IdentifierExpression{name}
}

// Source implements Node
func (e IdentifierExpression) Source() string { return e.Name }

// ConstantFold implements Node
func (e IdentifierExpression) ConstantFold() Node { return e }

// FunctionCallExpression represents a function call expression.
type FunctionCallExpression struct {
	Function  Node
	Arguments ExpressionList
}

// NewFunctionCallExpression creates a new function call expression.
func NewFunctionCallExpression(fn Node, args ExpressionList) FunctionCallExpression {
	return FunctionCallExpression{fn, args}
}

// Source implements Node
func (e FunctionCallExpression) Source() string {
	return fmt.Sprintf("%s(%s)", e.Function.Source(), e.Arguments.Source())
}

// ConstantFold implements Node
func (e FunctionCallExpression) ConstantFold() Node { return e }

// IndexExpression represents an indexing expression.
type IndexExpression struct {
	Operand Node
	Index   Node
}

// NewIndexExpression creates a new indexing expression.
func NewIndexExpression(op Node, index Node) IndexExpression {
	return IndexExpression{op, index}
}

// Source implements Node
func (e IndexExpression) Source() string {
	return fmt.Sprintf("%s[%s]", e.Operand.Source(), e.Index.Source())
}

// ConstantFold implements Node
func (e IndexExpression) ConstantFold() Node { return e }

// DotExpression represents a dot descend expression.
type DotExpression struct {
	Left  Node
	Right IdentifierExpression
}

// NewDotExpression creates a new dot descend expression.
func NewDotExpression(left Node, right IdentifierExpression) DotExpression {
	return DotExpression{left, right}
}

// Source implements Node
func (e DotExpression) Source() string { return e.Left.Source() + "." + e.Right.Source() }

// ConstantFold implements Node
func (e DotExpression) ConstantFold() Node { return e }

// NegateExpression represents a unary negation expression.
type NegateExpression struct {
	Operand Node
}

// NewNegateExpression creates a new unary negation expression.
func NewNegateExpression(op Node) NegateExpression {
	return NegateExpression{op}
}

// Source implements Node
func (e NegateExpression) Source() string { return "-" + e.Operand.Source() }

// ConstantFold implements Node
func (e NegateExpression) ConstantFold() Node {
	v := e.Operand.ConstantFold()

	if op, ok := v.(Constant); ok {
		if a, ok := op.Value.(value.Adder); ok {
			if c, err := a.Negate(); err == nil {
				return NewConstant(c)
			}
		}
	}
	return NegateExpression{v}
}

// LogicalNotExpression represents a logical not expression.
type LogicalNotExpression struct {
	Operand Node
}

// NewLogicalNotExpression creates a new logical not expression.
func NewLogicalNotExpression(op Node) LogicalNotExpression {
	return LogicalNotExpression{op}
}

// Source implements Node
func (e LogicalNotExpression) Source() string { return "!" + e.Operand.Source() }

// ConstantFold implements Node
func (e LogicalNotExpression) ConstantFold() Node {
	v := e.Operand.ConstantFold()

	if op, ok := v.(Constant); ok {
		if a, ok := op.Value.(value.Logical); ok {
			if c, err := a.LogicalNot(); err == nil {
				return NewConstant(c)
			}
		}
	}
	return NewLogicalNotExpression(v)
}

// BitwiseNotExpression represents a bitwise not expression.
type BitwiseNotExpression struct {
	Operand Node
}

// NewBitwiseNotExpression creates a new bitwise not expression.
func NewBitwiseNotExpression(op Node) BitwiseNotExpression {
	return BitwiseNotExpression{op}
}

// Source implements Node
func (e BitwiseNotExpression) Source() string { return "~" + e.Operand.Source() }

// ConstantFold implements Node
func (e BitwiseNotExpression) ConstantFold() Node {
	v := e.Operand.ConstantFold()

	if op, ok := v.(Constant); ok {
		if a, ok := op.Value.(value.Bitwise); ok {
			if c, err := a.BitwiseNot(); err == nil {
				return NewConstant(c)
			}
		}
	}
	return NewBitwiseNotExpression(v)
}

// MultiplyExpression represents a multiplication expression.
type MultiplyExpression struct {
	Left  Node
	Right Node
}

// NewMultiplyExpression creates a new multiplication expression.
func NewMultiplyExpression(left Node, right Node) MultiplyExpression {
	return MultiplyExpression{left, right}
}

// Source implements Node
func (e MultiplyExpression) Source() string {
	return fmt.Sprintf("%s * %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e MultiplyExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return MultiplyExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// DivideExpression represents a division expression.
type DivideExpression struct {
	Left  Node
	Right Node
}

// NewDivideExpression creates a new division expression.
func NewDivideExpression(left Node, right Node) DivideExpression {
	return DivideExpression{left, right}
}

// Source implements Node
func (e DivideExpression) Source() string {
	return fmt.Sprintf("%s / %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e DivideExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return DivideExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// ModuloExpression represents a modulo expression.
type ModuloExpression struct {
	Left  Node
	Right Node
}

// NewModuloExpression creates a new modulo expression.
func NewModuloExpression(left Node, right Node) ModuloExpression {
	return ModuloExpression{left, right}
}

// Source implements Node
func (e ModuloExpression) Source() string {
	return fmt.Sprintf("%s %% %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e ModuloExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return ModuloExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// AddExpression represents an addition expression.
type AddExpression struct {
	Left  Node
	Right Node
}

// NewAddExpression creates a new addition expression.
func NewAddExpression(left Node, right Node) AddExpression {
	return AddExpression{left, right}
}

// Source implements Node
func (e AddExpression) Source() string {
	return fmt.Sprintf("%s + %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e AddExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return AddExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// SubtractExpression represents a subtraction expression.
type SubtractExpression struct {
	Left  Node
	Right Node
}

// NewSubtractExpression creates a new subtraction expression.
func NewSubtractExpression(left Node, right Node) SubtractExpression {
	return SubtractExpression{left, right}
}

// ConstantFold implements Node
func (e SubtractExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return SubtractExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// Source implements Node
func (e SubtractExpression) Source() string {
	return fmt.Sprintf("%s - %s", e.Left.Source(), e.Right.Source())
}

// BitwiseLeftShiftExpression represents a bitwise left shift expression.
type BitwiseLeftShiftExpression struct {
	Left  Node
	Right Node
}

// NewBitwiseLeftShiftExpression creates a new bitwise left shift expression.
func NewBitwiseLeftShiftExpression(left Node, right Node) BitwiseLeftShiftExpression {
	return BitwiseLeftShiftExpression{left, right}
}

// Source implements Node
func (e BitwiseLeftShiftExpression) Source() string {
	return fmt.Sprintf("%s << %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e BitwiseLeftShiftExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return BitwiseLeftShiftExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// BitwiseRightShiftExpression represents a bitwise right shift expression.
type BitwiseRightShiftExpression struct {
	Left  Node
	Right Node
}

// NewBitwiseRightShiftExpression creates a new bitwise right shift expression.
func NewBitwiseRightShiftExpression(left Node, right Node) BitwiseRightShiftExpression {
	return BitwiseRightShiftExpression{left, right}
}

// Source implements Node
func (e BitwiseRightShiftExpression) Source() string {
	return fmt.Sprintf("%s >> %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e BitwiseRightShiftExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return BitwiseRightShiftExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// GreaterThanExpression represents a greater than relational expression.
type GreaterThanExpression struct {
	Left  Node
	Right Node
}

// NewGreaterThanExpression creates a new greater than relational expression.
func NewGreaterThanExpression(left Node, right Node) GreaterThanExpression {
	return GreaterThanExpression{left, right}
}

// ConstantFold implements Node
func (e GreaterThanExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return GreaterThanExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// Source implements Node
func (e GreaterThanExpression) Source() string {
	return fmt.Sprintf("%s > %s", e.Left.Source(), e.Right.Source())
}

// LessThanExpression represents a less than relational expression.
type LessThanExpression struct {
	Left  Node
	Right Node
}

// NewLessThanExpression creates a new less than relational expression.
func NewLessThanExpression(left Node, right Node) LessThanExpression {
	return LessThanExpression{left, right}
}

// Source implements Node
func (e LessThanExpression) Source() string {
	return fmt.Sprintf("%s < %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e LessThanExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return LessThanExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// GreaterThanOrEqualExpression represents a greater than or equal relational expression.
type GreaterThanOrEqualExpression struct {
	Left  Node
	Right Node
}

// NewGreaterThanOrEqualExpression creates a new greater than or equal relational expression.
func NewGreaterThanOrEqualExpression(left Node, right Node) GreaterThanOrEqualExpression {
	return GreaterThanOrEqualExpression{left, right}
}

// Source implements Node
func (e GreaterThanOrEqualExpression) Source() string {
	return fmt.Sprintf("%s >= %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e GreaterThanOrEqualExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return GreaterThanOrEqualExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// LessThanOrEqualExpression represents a less than or equal relational expression.
type LessThanOrEqualExpression struct {
	Left  Node
	Right Node
}

// NewLessThanOrEqualExpression creates a new less than or equal relational expression.
func NewLessThanOrEqualExpression(left Node, right Node) LessThanOrEqualExpression {
	return LessThanOrEqualExpression{left, right}
}

// Source implements Node
func (e LessThanOrEqualExpression) Source() string {
	return fmt.Sprintf("%s <= %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e LessThanOrEqualExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return LessThanOrEqualExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// EqualExpression represents an equality expression.
type EqualExpression struct {
	Left  Node
	Right Node
}

// NewEqualExpression creates a new equality expression.
func NewEqualExpression(left Node, right Node) EqualExpression {
	return EqualExpression{left, right}
}

// Source implements Node
func (e EqualExpression) Source() string {
	return fmt.Sprintf("%s == %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e EqualExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return EqualExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// NotEqualExpression represents an inequality expression.
type NotEqualExpression struct {
	Left  Node
	Right Node
}

// NewNotEqualExpression creates a new inequality expression.
func NewNotEqualExpression(left Node, right Node) NotEqualExpression {
	return NotEqualExpression{left, right}
}

// Source implements Node
func (e NotEqualExpression) Source() string {
	return fmt.Sprintf("%s != %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e NotEqualExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return NotEqualExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// BitwiseAndExpression represents a bitwise and expression.
type BitwiseAndExpression struct {
	Left  Node
	Right Node
}

// NewBitwiseAndExpression creates a new bitwise and expression.
func NewBitwiseAndExpression(left Node, right Node) BitwiseAndExpression {
	return BitwiseAndExpression{left, right}
}

// Source implements Node
func (e BitwiseAndExpression) Source() string {
	return fmt.Sprintf("%s & %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e BitwiseAndExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return BitwiseAndExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// BitwiseClearExpression represents a bitwise clear expression.
type BitwiseClearExpression struct {
	Left  Node
	Right Node
}

// NewBitwiseClearExpression creates a new bitwise and expression.
func NewBitwiseClearExpression(left Node, right Node) BitwiseClearExpression {
	return BitwiseClearExpression{left, right}
}

// Source implements Node
func (e BitwiseClearExpression) Source() string {
	return fmt.Sprintf("%s &^ %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e BitwiseClearExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return BitwiseClearExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// BitwiseXorExpression represents a bitwise exclusive or expression.
type BitwiseXorExpression struct {
	Left  Node
	Right Node
}

// NewBitwiseXorExpression creates a new bitwise exclusive or expression.
func NewBitwiseXorExpression(left Node, right Node) BitwiseXorExpression {
	return BitwiseXorExpression{left, right}
}

// Source implements Node
func (e BitwiseXorExpression) Source() string {
	return fmt.Sprintf("%s ^ %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e BitwiseXorExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return BitwiseXorExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// BitwiseOrExpression represents an inequality expression.
type BitwiseOrExpression struct {
	Left  Node
	Right Node
}

// NewBitwiseOrExpression creates a new bitwise or expression.
func NewBitwiseOrExpression(left Node, right Node) BitwiseOrExpression {
	return BitwiseOrExpression{left, right}
}

// Source implements Node
func (e BitwiseOrExpression) Source() string {
	return fmt.Sprintf("%s | %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e BitwiseOrExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return BitwiseOrExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// LogicalAndExpression represents a logical and expression.
type LogicalAndExpression struct {
	Left  Node
	Right Node
}

// NewLogicalAndExpression creates a new logical and expression.
func NewLogicalAndExpression(left Node, right Node) LogicalAndExpression {
	return LogicalAndExpression{left, right}
}

// Source implements Node
func (e LogicalAndExpression) Source() string {
	return fmt.Sprintf("%s && %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e LogicalAndExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return LogicalAndExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// LogicalOrExpression represents a logical or expression.
type LogicalOrExpression struct {
	Left  Node
	Right Node
}

// NewLogicalOrExpression creates a new logical or expression.
func NewLogicalOrExpression(left Node, right Node) LogicalOrExpression {
	return LogicalOrExpression{left, right}
}

// Source implements Node
func (e LogicalOrExpression) Source() string {
	return fmt.Sprintf("%s || %s", e.Left.Source(), e.Right.Source())
}

// ConstantFold implements Node
func (e LogicalOrExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return LogicalOrExpression{e.Left.ConstantFold(), e.Right.ConstantFold()}
}

// ConditionalExpression represents a conditional branching expression.
type ConditionalExpression struct {
	Condition Node
	Then      Node
	Else      Node
}

// NewConditionalExpression creates a new logical or expression.
func NewConditionalExpression(condition, then, els Node) ConditionalExpression {
	return ConditionalExpression{condition, then, els}
}

// Source implements Node
func (e ConditionalExpression) Source() string {
	return fmt.Sprintf("%s ? %s : %s", e.Condition.Source(), e.Then.Source(), e.Else.Source())
}

// ConstantFold implements Node
func (e ConditionalExpression) ConstantFold() Node {
	// TODO: implement constant folding for binary operators
	return ConditionalExpression{e.Condition.ConstantFold(), e.Then.ConstantFold(), e.Else.ConstantFold()}
}

// Constant represents a constant.
type Constant struct {
	Value value.Value
}

// NewConstant creates a new constant.
func NewConstant(value value.Value) Constant {
	return Constant{value}
}

// NewConstantErr creates a new constant, or passes through the error.
func NewConstantErr(value value.Value, err error) (Constant, error) {
	if err != nil {
		return NewConstant(nil), err
	}
	return NewConstant(value), nil
}

// Source implements Node
func (e Constant) Source() string {
	return e.Value.String()
}

// ConstantFold implements Node
func (e Constant) ConstantFold() Node {
	return e
}
