package compiler

import (
	"bytes"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/go-restruct/restruct/internal/expr"
	"github.com/go-restruct/restruct/internal/expr/typing"
	"github.com/stretchr/testify/assert"
)

type Struct struct {
	Field string
}

func genTest(in string) (string, error) {
	var buf bytes.Buffer

	rootType, err := typing.FromValue(Struct{})
	if err != nil {
		return "", err
	}
	context := typing.Context{
		Self:   rootType,
		Global: map[string]typing.Type{},
	}

	expr, err := expr.ParseString(in)
	if err != nil {
		return "", err
	}
	expr = expr.ConstantFold()

	result, err := Compile("expr", context, expr)
	if err != nil {
		return "", err
	}

	f := jen.NewFilePathName("github.com/go-restruct/restruct/internal/expr/compiler", "compiler")
	f.Add(result)

	err = f.Render(&buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func TestCompile(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{{"Field", `package compiler

func expr(self Struct) string {
	return self.Field
}
`}, {"((Field[len(Field)-1] * 2 > 4 ? 0 : 1) == 0 ? (true) : (false))", `package compiler

func expr(self Struct) bool {
	var temp1 int64
	if self.Field[len(self.Field)-1]*2 > 4 {
		temp1 = 0
	} else {
		temp1 = 1
	}
	if (temp1) == 0 {
		return true
	} else {
		return false
	}
}
`}}

	for _, test := range tests {
		result, err := genTest(test.in)
		if err != nil {
			t.Error(err)
			continue
		}
		assert.Equal(t, test.out, result)
	}
}
