package restruct

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTag(t *testing.T) {
	tests := []struct {
		input  string
		opts   tagOptions
		errstr string
	}{
		// Blank
		{"", tagOptions{}, ""},

		// Gibberish
		{"!#%,1245^#df,little,&~@~@~@~@", tagOptions{}, "parsing error"},
		{"&paREZysLu&@,83D9I!,9OsQ56BLD", tagOptions{}, "parsing error"},
		{"B7~,H0IDSxDlJ,#xa$kgDEL%Ts,88", tagOptions{}, "parsing error"},
		{"fio&eQ8xwbhAWR*!CRlL2XBDG$45s", tagOptions{}, "parsing error"},
		{"IcPyRJ#EV@a4QAb9wENk4Zq9MpX$p", tagOptions{}, "parsing error"},

		// Conflicting byte order
		{"little,big", tagOptions{Order: binary.BigEndian}, ""},
		{"big,little", tagOptions{Order: binary.LittleEndian}, ""},

		// Byte order
		{"msb", tagOptions{Order: binary.BigEndian}, ""},
		{"lsb", tagOptions{Order: binary.LittleEndian}, ""},
		{"network", tagOptions{Order: binary.BigEndian}, ""},

		// Ignore
		{"-", tagOptions{Ignore: true}, ""},
		{"-,test", tagOptions{}, "extra options on ignored field"},

		// Bad types
		{"invalid", tagOptions{}, "unknown type invalid"},
		{"chan int8", tagOptions{}, "channel type not allowed"},
		{"map[byte]byte", tagOptions{}, "map type not allowed"},

		// Types
		{"uint8", tagOptions{Type: reflect.TypeOf(uint8(0))}, ""},
		{"uint16", tagOptions{Type: reflect.TypeOf(uint16(0))}, ""},
		{"uint32", tagOptions{Type: reflect.TypeOf(uint32(0))}, ""},
		{"int8", tagOptions{Type: reflect.TypeOf(int8(0))}, ""},
		{"int16", tagOptions{Type: reflect.TypeOf(int16(0))}, ""},
		{"int32", tagOptions{Type: reflect.TypeOf(int32(0))}, ""},

		// Bitfields
		{"uint8:3", tagOptions{Type: reflect.TypeOf(uint8(0)), BitSize: 3}, ""},
		{"uint16:15", tagOptions{Type: reflect.TypeOf(uint16(0)), BitSize: 15}, ""},
		{"uint32:31", tagOptions{Type: reflect.TypeOf(uint32(0)), BitSize: 31}, ""},
		{"uint64:63", tagOptions{Type: reflect.TypeOf(uint64(0)), BitSize: 63}, ""},
		{"uint8:1", tagOptions{Type: reflect.TypeOf(uint8(0)), BitSize: 1}, ""},
		{"uint16:1", tagOptions{Type: reflect.TypeOf(uint16(0)), BitSize: 1}, ""},
		{"uint32:1", tagOptions{Type: reflect.TypeOf(uint32(0)), BitSize: 1}, ""},
		{"uint64:1", tagOptions{Type: reflect.TypeOf(uint64(0)), BitSize: 1}, ""},

		// Wrong bitfields
		{"uint8:0", tagOptions{}, "Bad value on bitfield"},
		{"uint16:0", tagOptions{}, "Bad value on bitfield"},
		{"uint32:0", tagOptions{}, "Bad value on bitfield"},
		{"uint64:0", tagOptions{}, "Bad value on bitfield"},
		{"int8:0", tagOptions{}, "Bad value on bitfield"},
		{"int16:0", tagOptions{}, "Bad value on bitfield"},
		{"int32:0", tagOptions{}, "Bad value on bitfield"},
		{"int64:0", tagOptions{}, "Bad value on bitfield"},

		{"uint8:8", tagOptions{BitSize: 0}, "Too high value on bitfield"},
		{"uint16:16", tagOptions{BitSize: 0}, "Too high value on bitfield"},
		{"uint32:32", tagOptions{BitSize: 0}, "Too high value on bitfield"},
		{"uint64:64", tagOptions{BitSize: 0}, "Too high value on bitfield"},
		{"int8:8", tagOptions{BitSize: 0}, "Too high value on bitfield"},
		{"int16:16", tagOptions{BitSize: 0}, "Too high value on bitfield"},
		{"int32:32", tagOptions{BitSize: 0}, "Too high value on bitfield"},
		{"int64:64", tagOptions{BitSize: 0}, "Too high value on bitfield"},

		{"uint8:XX", tagOptions{BitSize: 0}, "Bad value on bitfield"},
		{"uint16:XX", tagOptions{BitSize: 0}, "Bad value on bitfield"},
		{"uint32:XX", tagOptions{BitSize: 0}, "Bad value on bitfield"},
		{"uint64:XX", tagOptions{BitSize: 0}, "Bad value on bitfield"},
		{"int8:XX", tagOptions{BitSize: 0}, "Bad value on bitfield"},
		{"int16:XX", tagOptions{BitSize: 0}, "Bad value on bitfield"},
		{"int32:XX", tagOptions{BitSize: 0}, "Bad value on bitfield"},
		{"int64:XX", tagOptions{BitSize: 0}, "Bad value on bitfield"},

		{"uint8:X:X", tagOptions{BitSize: 0}, "extra options on type field"},
		{"uint16:X:X", tagOptions{BitSize: 0}, "extra options on type field"},
		{"uint32:X:X", tagOptions{BitSize: 0}, "extra options on type field"},
		{"uint64:X:X", tagOptions{BitSize: 0}, "extra options on type field"},
		{"int8:X:X", tagOptions{BitSize: 0}, "extra options on type field"},
		{"int16:X:X", tagOptions{BitSize: 0}, "extra options on type field"},
		{"int32:X:X", tagOptions{BitSize: 0}, "extra options on type field"},
		{"int64:X:X", tagOptions{BitSize: 0}, "extra options on type field"},

		// Sizeof
		{"sizeof=OtherField", tagOptions{SizeOf: "OtherField"}, ""},
		{"sizeof=日本", tagOptions{SizeOf: "日本"}, ""},

		// Skip
		{"skip=4", tagOptions{Skip: 4}, ""},
		{"skip=字", tagOptions{}, "bad skip amount"},

		// Composite
		{"uint16,little,sizeof=test,skip=5",
			tagOptions{
				Type:   reflect.TypeOf(uint16(0)),
				Order:  binary.LittleEndian,
				SizeOf: "test",
				Skip:   5,
			},
			"",
		},
	}

	for _, test := range tests {
		opts, err := parseTag(test.input)
		assert.Equal(t, test.opts, opts)
		if err != nil {
			assert.Equal(t, test.errstr, err.Error())
		}
	}
}

func TestMustParseTagPanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Invalid tag did not panic.")
		}
	}()
	mustParseTag("???")
}

func TestMustParseTagReturnsOnSuccess(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Valid tag panicked.")
		}
	}()
	mustParseTag("[128]byte,little,sizeof=Test")
}
