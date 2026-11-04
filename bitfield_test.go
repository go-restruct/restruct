package restruct

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertBitfieldRoundTrip(t *testing.T, order binary.ByteOrder, value interface{}, expected []byte) {
	t.Helper()
	input := reflect.New(reflect.TypeOf(value))
	input.Elem().Set(reflect.ValueOf(value))
	data, err := Pack(order, input.Interface())
	assert.Nil(t, err)
	assert.Equal(t, expected, data)

	decoded := reflect.New(reflect.TypeOf(value))
	err = Unpack(expected, order, decoded.Interface())
	assert.Nil(t, err)
	assert.Equal(t, value, decoded.Elem().Interface())

	decoded = reflect.New(reflect.TypeOf(value))
	err = Unpack(data, order, decoded.Interface())
	assert.Nil(t, err)
	assert.Equal(t, value, decoded.Elem().Interface())
}

func TestSignedNonByteAlignedBitfieldRoundTrip(t *testing.T) {
	tests8 := []struct {
		value    int8
		unsigned uint8
		data     []byte
		prefixed []byte
	}{
		{0, 0, []byte{0x00}, []byte{0xa3}},
		{1, 1, []byte{0x20}, []byte{0xa7}},
		{2, 2, []byte{0x40}, []byte{0xab}},
		{3, 3, []byte{0x60}, []byte{0xaf}},
		{-4, 4, []byte{0x80}, []byte{0xb3}},
		{-3, 5, []byte{0xa0}, []byte{0xb7}},
		{-2, 6, []byte{0xc0}, []byte{0xbb}},
		{-1, 7, []byte{0xe0}, []byte{0xbf}},
	}
	for _, test := range tests8 {
		t.Run(fmt.Sprintf("int8/%d", test.value), func(t *testing.T) {
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Value int8 `struct:"int8:3"`
			}{test.value}, test.data)
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Prefix uint8 `struct:"uint8:3"`
				Value  int8  `struct:"int8:3"`
				Suffix uint8 `struct:"uint8:2"`
			}{5, test.value, 3}, test.prefixed)
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Value uint8 `struct:"uint8:3"`
			}{test.unsigned}, test.data)
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Prefix uint8 `struct:"uint8:3"`
				Value  uint8 `struct:"uint8:3"`
				Suffix uint8 `struct:"uint8:2"`
			}{5, test.unsigned, 3}, test.prefixed)
		})
	}

	tests16 := []struct {
		value    int16
		unsigned uint16
		data     []byte
		prefixed []byte
	}{
		{0, 0, []byte{0x00, 0x00}, []byte{0xa0, 0x03}},
		{1, 1, []byte{0x00, 0x20}, []byte{0xa0, 0x07}},
		{2, 2, []byte{0x00, 0x40}, []byte{0xa0, 0x0b}},
		{1023, 1023, []byte{0x7f, 0xe0}, []byte{0xaf, 0xff}},
		{-1024, 1024, []byte{0x80, 0x00}, []byte{0xb0, 0x03}},
		{-1023, 1025, []byte{0x80, 0x20}, []byte{0xb0, 0x07}},
		{-2, 2046, []byte{0xff, 0xc0}, []byte{0xbf, 0xfb}},
		{-1, 2047, []byte{0xff, 0xe0}, []byte{0xbf, 0xff}},
	}
	for _, test := range tests16 {
		t.Run(fmt.Sprintf("int16/%d", test.value), func(t *testing.T) {
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Value int16 `struct:"int16:11"`
			}{test.value}, test.data)
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Prefix uint8 `struct:"uint8:3"`
				Value  int16 `struct:"int16:11"`
				Suffix uint8 `struct:"uint8:2"`
			}{5, test.value, 3}, test.prefixed)
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Value uint16 `struct:"uint16:11"`
			}{test.unsigned}, test.data)
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Prefix uint8  `struct:"uint8:3"`
				Value  uint16 `struct:"uint16:11"`
				Suffix uint8  `struct:"uint8:2"`
			}{5, test.unsigned, 3}, test.prefixed)
		})
	}

	tests32 := []struct {
		value    int32
		unsigned uint32
		data     []byte
		prefixed []byte
	}{
		{0, 0, []byte{0x00, 0x00, 0x00}, []byte{0xa0, 0x00, 0x00, 0xc0}},
		{1, 1, []byte{0x00, 0x00, 0x08}, []byte{0xa0, 0x00, 0x01, 0xc0}},
		{2, 2, []byte{0x00, 0x00, 0x10}, []byte{0xa0, 0x00, 0x02, 0xc0}},
		{1048575, 1048575, []byte{0x7f, 0xff, 0xf8}, []byte{0xaf, 0xff, 0xff, 0xc0}},
		{-1048576, 1048576, []byte{0x80, 0x00, 0x00}, []byte{0xb0, 0x00, 0x00, 0xc0}},
		{-1048575, 1048577, []byte{0x80, 0x00, 0x08}, []byte{0xb0, 0x00, 0x01, 0xc0}},
		{-2, 2097150, []byte{0xff, 0xff, 0xf0}, []byte{0xbf, 0xff, 0xfe, 0xc0}},
		{-1, 2097151, []byte{0xff, 0xff, 0xf8}, []byte{0xbf, 0xff, 0xff, 0xc0}},
	}
	for _, test := range tests32 {
		t.Run(fmt.Sprintf("int32/%d", test.value), func(t *testing.T) {
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Value int32 `struct:"int32:21"`
			}{test.value}, test.data)
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Prefix uint8 `struct:"uint8:3"`
				Value  int32 `struct:"int32:21"`
				Suffix uint8 `struct:"uint8:2"`
			}{5, test.value, 3}, test.prefixed)
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Value uint32 `struct:"uint32:21"`
			}{test.unsigned}, test.data)
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Prefix uint8  `struct:"uint8:3"`
				Value  uint32 `struct:"uint32:21"`
				Suffix uint8  `struct:"uint8:2"`
			}{5, test.unsigned, 3}, test.prefixed)
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Value int64 `struct:"int64:21"`
			}{int64(test.value)}, test.data)
			assertBitfieldRoundTrip(t, binary.BigEndian, struct {
				Prefix uint8 `struct:"uint8:3"`
				Value  int64 `struct:"int64:21"`
				Suffix uint8 `struct:"uint8:2"`
			}{5, int64(test.value), 3}, test.prefixed)
		})
	}
}

func TestSignedFullWidthRoundTrip(t *testing.T) {
	tests := []struct {
		a      int8
		b      int16
		c      int32
		d      int64
		big    []byte
		little []byte
	}{
		{1, 2, 2147483647, 9223372036854775807,
			[]byte{0x01, 0x00, 0x02, 0x7f, 0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			[]byte{0x01, 0x02, 0x00, 0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}},
		{-128, -32768, -2147483648, -9223372036854775808,
			[]byte{0x80, 0x80, 0x00, 0x80, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			[]byte{0x80, 0x00, 0x80, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.a), func(t *testing.T) {
			untagged := struct {
				A int8
				B int16
				C int32
				D int64
			}{test.a, test.b, test.c, test.d}
			tagged := struct {
				A int8  `struct:"int8"`
				B int16 `struct:"int16"`
				C int32 `struct:"int32"`
				D int64 `struct:"int64"`
			}{test.a, test.b, test.c, test.d}
			assertBitfieldRoundTrip(t, binary.BigEndian, untagged, test.big)
			assertBitfieldRoundTrip(t, binary.BigEndian, tagged, test.big)
			assertBitfieldRoundTrip(t, binary.LittleEndian, untagged, test.little)
			assertBitfieldRoundTrip(t, binary.LittleEndian, tagged, test.little)
		})
	}
}
