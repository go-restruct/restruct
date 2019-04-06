package expr

import (
	"github.com/go-restruct/restruct/internal/expr/ast"
	"github.com/go-restruct/restruct/internal/expr/lexer"
	"github.com/go-restruct/restruct/internal/expr/parser"
)

// ParseString parses a string into an expression AST.
func ParseString(input string) (expr ast.Node, err error) {
	scanner := lexer.NewLexer([]byte(input))
	parser := parser.NewParser()
	root, err := parser.Parse(scanner)
	if err != nil {
		return nil, err
	}
	return root.(ast.Node), nil
}
