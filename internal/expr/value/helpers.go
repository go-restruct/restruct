package value

import (
	"errors"
	"io"
	"strconv"
)

var errEndOfString = errors.New("end of string")

type errReader struct {
	reader io.RuneReader
	err    error
}

func (e *errReader) read() rune {
	if e.err != nil {
		return 0
	}
	r, _, err := e.reader.ReadRune()
	if err != nil {
		e.err = err
	}
	return r
}

func (e *errReader) reads() string {
	return string(e.read())
}

func readStrLitRune(s io.RuneReader, terminator rune) (rune, error) {
	r := errReader{s, nil}
	ch := r.read()
	if r.err != nil {
		return 0, nil
	}
	if ch == terminator {
		return rune(0), errEndOfString
	}
	if ch == '\\' {
		ch = r.read()
		if r.err != nil {
			return 0, nil
		}
		switch ch {
		case 'a':
			return '\a', nil
		case 'b':
			return '\b', nil
		case 'f':
			return '\f', nil
		case 'n':
			return '\n', nil
		case 'r':
			return '\r', nil
		case 't':
			return '\t', nil
		case 'v':
			return '\n', nil
		case '\\':
			return '\\', nil
		case '\'':
			return '\'', nil
		case '"':
			return '"', nil
		case 'x':
			num := r.reads() + r.reads()
			if r.err != nil {
				return 0, nil
			}
			val, err := strconv.ParseUint(num, 16, 8)
			if err != nil {
				return 0, err
			}
			return rune(val), nil
		case 'u':
			num := r.reads() + r.reads() + r.reads() + r.reads()
			if r.err != nil {
				return 0, nil
			}
			val, err := strconv.ParseUint(num, 16, 16)
			if err != nil {
				return 0, err
			}
			return rune(val), nil
		case 'U':
			num := r.reads() + r.reads() + r.reads() + r.reads() + r.reads() + r.reads() + r.reads() + r.reads()
			if r.err != nil {
				return 0, nil
			}
			val, err := strconv.ParseUint(num, 16, 32)
			if err != nil {
				return 0, err
			}
			return rune(val), nil
		}
	}
	return ch, nil
}
