package restruct

import (
	"encoding/binary"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// tagOptions represents a parsed struct tag.
type tagOptions struct {
	Ignore bool
	Type   reflect.Type
	SizeOf string
	Skip   int
	Order  binary.ByteOrder
}

// mustParseTag calls ParseTag but panics if there is an error, to help make
// sure programming errors surface quickly.
func mustParseTag(tag string) tagOptions {
	opt, err := parseTag(tag)
	if err != nil {
		panic(err)
	}
	return opt
}

// parseTag parses a struct tag into a TagOptions structure.
func parseTag(tag string) (tagOptions, error) {
	parts := strings.Split(tag, ",")

	if len(tag) == 0 || len(parts) == 0 {
		return tagOptions{}, nil
	}

	// Handle `struct:"-"`
	if parts[0] == "-" {
		if len(parts) > 1 {
			return tagOptions{}, errors.New("extra options on ignored field")
		}
		return tagOptions{Ignore: true}, nil
	}

	result := tagOptions{}
	for _, part := range parts {
		switch part {
		case "lsb", "little":
			result.Order = binary.LittleEndian
			continue
		case "msb", "big", "network":
			result.Order = binary.BigEndian
			continue
		default:
			if strings.HasPrefix(part, "sizeof=") {
				result.SizeOf = part[7:]
				continue
			} else if strings.HasPrefix(part, "skip=") {
				var err error
				result.Skip, err = strconv.Atoi(part[5:])
				if err != nil {
					return tagOptions{}, errors.New("bad skip amount")
				}
			} else {
				typ, err := parseType(part)
				if err != nil {
					return tagOptions{}, err
				}
				result.Type = typ
				continue
			}
		}
	}

	return result, nil
}
