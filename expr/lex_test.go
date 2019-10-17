package expr

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanValid(t *testing.T) {
	tests := []struct {
		input    string
		expected []token
	}{
		{
			input: "a + b - c * d / e ^ f | g & h &^ i % k",
			expected: []token{
				{kind: identtoken, sval: "a", pos: 0},
				{kind: addtoken, sval: "+", pos: 2},
				{kind: identtoken, sval: "b", pos: 4},
				{kind: subtoken, sval: "-", pos: 6},
				{kind: identtoken, sval: "c", pos: 8},
				{kind: multoken, sval: "*", pos: 10},
				{kind: identtoken, sval: "d", pos: 12},
				{kind: quotoken, sval: "/", pos: 14},
				{kind: identtoken, sval: "e", pos: 16},
				{kind: xortoken, sval: "^", pos: 18},
				{kind: identtoken, sval: "f", pos: 20},
				{kind: ortoken, sval: "|", pos: 22},
				{kind: identtoken, sval: "g", pos: 24},
				{kind: andtoken, sval: "&", pos: 26},
				{kind: identtoken, sval: "h", pos: 28},
				{kind: andnottoken, sval: "&^", pos: 30},
				{kind: identtoken, sval: "i", pos: 33},
				{kind: remtoken, sval: "%", pos: 35},
				{kind: identtoken, sval: "k", pos: 37},
				{kind: eoftoken, pos: 38},
			},
		},
		{
			input: "l > m < n >= o <= p == q != r >> s << t || u && v",
			expected: []token{
				{kind: identtoken, sval: "l", pos: 0},
				{kind: greatertoken, sval: ">", pos: 2},
				{kind: identtoken, sval: "m", pos: 4},
				{kind: lessertoken, sval: "<", pos: 6},
				{kind: identtoken, sval: "n", pos: 8},
				{kind: greaterequaltoken, sval: ">=", pos: 10},
				{kind: identtoken, sval: "o", pos: 13},
				{kind: lesserequaltoken, sval: "<=", pos: 15},
				{kind: identtoken, sval: "p", pos: 18},
				{kind: equaltoken, sval: "==", pos: 20},
				{kind: identtoken, sval: "q", pos: 23},
				{kind: notequaltoken, sval: "!=", pos: 25},
				{kind: identtoken, sval: "r", pos: 28},
				{kind: shrtoken, sval: ">>", pos: 30},
				{kind: identtoken, sval: "s", pos: 33},
				{kind: shltoken, sval: "<<", pos: 35},
				{kind: identtoken, sval: "t", pos: 38},
				{kind: logicalortoken, sval: "||", pos: 40},
				{kind: identtoken, sval: "u", pos: 43},
				{kind: logicalandtoken, sval: "&&", pos: 45},
				{kind: identtoken, sval: "v", pos: 48},
				{kind: eoftoken, pos: 49},
			},
		},
		{
			input: "(A + Test) % C[0]",
			expected: []token{
				{kind: leftparentoken, pos: 0, sval: "("},
				{kind: identtoken, pos: 1, sval: "A"},
				{kind: addtoken, pos: 3, sval: "+"},
				{kind: identtoken, pos: 5, sval: "Test"},
				{kind: rightparentoken, pos: 9, sval: ")"},
				{kind: remtoken, pos: 11, sval: "%"},
				{kind: identtoken, pos: 13, sval: "C"},
				{kind: leftbrackettoken, pos: 14, sval: "["},
				{kind: inttoken, pos: 15, sval: "0"},
				{kind: rightbrackettoken, pos: 16, sval: "]"},
				{kind: eoftoken, pos: 17},
			},
		},
		{
			input: "0x80000000 + 0.1",
			expected: []token{
				{kind: inttoken, pos: 0, sval: "0x80000000", ival: 0x80000000, uval: 0x80000000, fval: 0x80000000},
				{kind: addtoken, pos: 11, sval: "+"},
				{kind: floattoken, pos: 13, sval: "0.1", ival: 0, uval: 0, fval: 0.1},
				{kind: eoftoken, pos: 16},
			},
		},
		{
			input: "0x80000000 + 0.1",
			expected: []token{
				{kind: inttoken, pos: 0, sval: "0x80000000", ival: 0x80000000, uval: 0x80000000, fval: 0x80000000},
				{kind: addtoken, pos: 11, sval: "+"},
				{kind: floattoken, pos: 13, sval: "0.1", ival: 0, uval: 0, fval: 0.1},
				{kind: eoftoken, pos: 16},
			},
		},
		{
			input: "-0x80000000 - -0.1",
			expected: []token{
				{kind: inttoken, pos: 0, sval: "-0x80000000", ival: -0x80000000, uval: 0xFFFFFFFF80000000, fval: -0x80000000, sign: true},
				{kind: subtoken, pos: 12, sval: "-"},
				{kind: floattoken, pos: 14, sval: "-0.1", ival: 0, uval: 0, fval: -0.1},
				{kind: eoftoken, pos: 18},
			},
		},
		{
			input: `"a b c" + " d e f"`,
			expected: []token{
				{kind: strtoken, pos: 0, sval: "a b c"},
				{kind: addtoken, pos: 8, sval: "+"},
				{kind: strtoken, pos: 10, sval: " d e f"},
				{kind: eoftoken, pos: 18},
			},
		},
		{
			input: `'a' + 0`,
			expected: []token{
				{kind: runetoken, pos: 0, ival: 'a'},
				{kind: addtoken, pos: 4, sval: "+"},
				{kind: inttoken, pos: 6, sval: "0", ival: 0, uval: 0, fval: 0},
				{kind: eoftoken, pos: 7},
			},
		},
		{
			input: `"aä本\t\000\007\377\x07\xff\u12e4\U00101234\""`,
			expected: []token{
				{kind: strtoken, pos: 0, sval: "aä本\t\000\007\377\x07\xff\u12e4\U00101234\""},
				{kind: eoftoken, pos: 45},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual := []token{}
			s := newscanner(bytes.NewBufferString(test.input))
			for {
				t := s.scan()
				actual = append(actual, t)
				if t.kind == eoftoken || t.kind == errtoken {
					break
				}
				if len(actual) > 500 {
					panic("Maximum test token limit exceeded.")
				}
			}
			assert.Equal(t, test.expected, actual)
		})
	}
}
