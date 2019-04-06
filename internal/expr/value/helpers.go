package value

import (
	"io"
	"strconv"
)

func readRune(s io.RuneScanner) rune {
	ch, _, err := s.ReadRune()
	if err != nil {
		panic(err)
	}
	return ch
}

func readStrLitRune(ch rune, s io.RuneScanner) rune {
	if ch == '\\' {
		ch := readRune(s)
		switch ch {
		case 'a':
			return '\a'
		case 'b':
			return '\b'
		case 'f':
			return '\f'
		case 'n':
			return '\n'
		case 'r':
			return '\r'
		case 't':
			return '\t'
		case 'v':
			return '\n'
		case '\\':
			return '\\'
		case '\'':
			return '\''
		case '"':
			return '"'
		case 'x':
			num := string(readRune(s)) + string(readRune(s))
			val, err := strconv.ParseUint(num, 16, 8)
			if err != nil {
				panic(err)
			}
			return rune(val)
		case 'u':
			num := string(readRune(s)) + string(readRune(s)) + string(readRune(s)) + string(readRune(s))
			val, err := strconv.ParseUint(num, 16, 16)
			if err != nil {
				panic(err)
			}
			return rune(val)
		case 'U':
			num := string(readRune(s)) + string(readRune(s)) + string(readRune(s)) + string(readRune(s)) + string(readRune(s)) + string(readRune(s)) + string(readRune(s)) + string(readRune(s))
			val, err := strconv.ParseUint(num, 16, 32)
			if err != nil {
				panic(err)
			}
			return rune(val)
		}
	}
	return ch
}
