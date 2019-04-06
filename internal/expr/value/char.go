package value

import (
	"errors"
	"strings"
)

// ParseCharLiteral parses a character literal.
func ParseCharLiteral(literal string) (value Uint, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	r := strings.NewReader(literal)
	if readRune(r) != '\'' {
		panic(errors.New("syntax error: expected '"))
	}
	lit := readStrLitRune(readRune(r), r)
	if readRune(r) != '\'' {
		panic(errors.New("syntax error: expected '"))
	}

	return Uint{uint64(lit)}, nil
}
