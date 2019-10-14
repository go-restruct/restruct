package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSource(t *testing.T) {
	tests := []struct {
		input  node
		output string
	}{
		{
			binaryexpr{
				binaryadd,
				identnode{0, "a"},
				binaryexpr{
					binarymul,
					identnode{4, "b"},
					identnode{8, "c"},
				},
			},
			"a + b * c",
		},
		{
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
			"a == 1 ? b + 1 : c * 1",
		},
		{
			binaryexpr{
				binarycall,
				identnode{0, "a"},
				binaryexpr{
					binarygroup,
					identnode{2, "b"},
					identnode{5, "c"},
				},
			},
			"a(b, c)",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.output, test.input.source())
	}
}
