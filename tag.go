package restruct

import (
	"encoding/binary"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// TagOptions represents a parsed struct tag.
type TagOptions struct {
	Ignore bool
	Type   reflect.Type
	SizeOf string
	Skip   int
	Order  binary.ByteOrder
}

// MustParseTag calls ParseTag but panics if there is an error, to help make
// sure programming errors surface quickly.
func MustParseTag(tag string) TagOptions {
	opt, err := ParseTag(tag)
	if err != nil {
		panic(err)
	}
	return opt
}

// ParseTag parses a struct tag into a TagOptions structure.
func ParseTag(tag string) (TagOptions, error) {
	parts := strings.Split(tag, ",")

	if len(tag) == 0 || len(parts) == 0 {
		return TagOptions{}, nil
	}

	// Handle `struct:"-"`
	if parts[0] == "-" {
		if len(parts) > 1 {
			return TagOptions{}, errors.New("extra options on ignored field")
		}
		return TagOptions{Ignore: true}, nil
	}

	result := TagOptions{}
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
					return TagOptions{}, errors.New("bad skip amount")
				}
			} else {
				typ, err := ParseType(part)
				if err != nil {
					return TagOptions{}, err
				}
				result.Type = typ
				continue
			}
		}
	}

	return result, nil
}
