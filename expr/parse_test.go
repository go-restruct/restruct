package expr

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBasic(t *testing.T) {
	tests := []struct {
		input  string
		output node
	}{
		{
			"a + b * c",
			binaryexpr{
				binaryadd,
				identnode{0, "a"},
				binaryexpr{
					binarymul,
					identnode{4, "b"},
					identnode{8, "c"},
				},
			},
		},
		{
			"a * b + c",
			binaryexpr{
				binaryadd,
				binaryexpr{
					binarymul,
					identnode{0, "a"},
					identnode{4, "b"},
				},
				identnode{8, "c"},
			},
		},
		{
			"a * (b + c)",
			binaryexpr{
				binarymul,
				identnode{0, "a"},
				binaryexpr{
					binaryadd,
					identnode{5, "b"},
					identnode{9, "c"},
				},
			},
		},
		{
			"a ? b : c",
			ternaryexpr{
				identnode{0, "a"},
				identnode{4, "b"},
				identnode{8, "c"},
			},
		},
		{
			"a == 1 ? b + 1 : c * 1",
			ternaryexpr{
				binaryexpr{
					binaryequal,
					identnode{0, "a"},
					intnode{5, 1, 1, false},
				},
				binaryexpr{
					binaryadd,
					identnode{9, "b"},
					intnode{13, 1, 1, false},
				},
				binaryexpr{
					binarymul,
					identnode{17, "c"},
					intnode{21, 1, 1, false},
				},
			},
		},
		{
			"a(b, c)",
			binaryexpr{
				binarycall,
				identnode{0, "a"},
				binaryexpr{
					binarygroup,
					identnode{2, "b"},
					identnode{5, "c"},
				},
			},
		},
		{
			"a[1]",
			binaryexpr{
				binarysubscript,
				identnode{0, "a"},
				intnode{2, 1, 1, false},
			},
		},
	}

	for _, test := range tests {
		p := newparser(newscanner(bytes.NewBufferString(test.input)))
		assert.Equal(t, test.output, p.parse())
	}
}
