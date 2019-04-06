package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicExpr(t *testing.T) {
	tests := []struct {
		script string
		result interface{}
		obj    interface{}
	}{
		{script: "1", result: int64(1)},
		{script: "-1", result: int64(-1)},
		{script: "1 - 2", result: int64(-1)},
		{script: "uint(-1)", result: uint64(0xffffffffffffffff)},
		{script: `"test"`, result: "test"},
		{script: `'a'`, result: uint64('a')},
		{script: `"test"[0]`, result: uint64('t')},
		{script: `len("test")`, result: uint64(4)},
		{script: `0.0`, result: float64(0.0)},
		{script: `0.0 == 0.0`, result: true},
		{script: `0.0 != 0.0`, result: false},
		{script: `.1`, result: float64(0.1)},
		{script: `1.`, result: float64(1.0)},
		{script: `1 | 2 | 4 | 8`, result: int64(15)},
		{script: `15 &^ 8`, result: int64(7)},
		{script: `2 << 2`, result: int64(8)},
		{script: `2 >> 1`, result: int64(1)},
		{script: `len("test") * 2`, result: uint64(8)},
		{script: `3 / 2`, result: int64(1)},
		{script: `3 % 2`, result: int64(1)},
		{script: `^0`, result: int64(-1)},
		{script: `1.2e2`, result: float64(120)},
		{script: `1.1 + 1`, result: float64(2.1)},
		{script: `uint(1.1 + 1)`, result: uint64(2)},
		{script: `16 << 1 == 32`, result: true},
		{script: `8 << 1 == 16 && 16 << 1 == 32`, result: true},
		{script: `"a"[0] > 0x20`, result: true},
		{script: `len(Val) + Val[0]`, result: uint64(0x78), obj: struct{ Val string }{Val: "test"}},
		{script: `Val`, result: int64(-31), obj: struct{ Val int32 }{Val: -31}},
		{script: `Val * 2`, result: int64(-62), obj: struct{ Val int32 }{Val: -31}},
		{script: `!!Nested.Value`, result: true, obj: struct{ Nested struct{ Value bool } }{Nested: struct{ Value bool }{true}}},
		{script: `!!Nested.Value`, result: true, obj: struct{ Nested *struct{ Value bool } }{Nested: &struct{ Value bool }{true}}},
		{script: `((("test")))`, result: "test"},
		{script: `sum(1, 2, 3)`, result: int64(6)},
		{script: `usum(1, 2, 3)`, result: uint64(6)},
		{script: `fsum(1, 2, 3)`, result: float64(6)},
	}

	for _, test := range tests {
		result, err := Eval(test.obj, test.script)
		if err != nil {
			t.Error(err)
			t.Fail()
			continue
		}
		assert.Equal(t, test.result, result)
	}
}
