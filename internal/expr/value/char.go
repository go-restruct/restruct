package value

import (
	"errors"
	"strings"
)

// ParseCharLiteral parses a character literal.
func ParseCharLiteral(literal string) (value Uint, err error) {
	r := strings.NewReader(literal)
	ch, _, err := r.ReadRune()
	if err != nil {
		return Uint{}, err
	} else if ch != '\'' {
		return Uint{}, errors.New("syntax error: expected '")
	}

	lit, err := readStrLitRune(r, '\'')
	if err != nil {
		return Uint{}, err
	}

	ch, _, err = r.ReadRune()
	if err != nil {
		return Uint{}, err
	} else if ch != '\'' {
		return Uint{}, errors.New("syntax error: expected '")
	}

	return Uint{uint64(lit)}, nil
}
