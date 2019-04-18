package restruct

import (
	"encoding/binary"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		data    []byte
		bitsize int
		value   interface{}
	}{
		{
			data: []byte{
				0x12, 0x34, 0x56, 0x78,
			},
			bitsize: 32,
			value: struct {
				Dd uint32
			}{
				Dd: 0x12345678,
			},
		},

		{
			data: []byte{
				0x55, 0x55,
			},
			bitsize: 16,
			value: struct {
				A uint8 `struct:"uint8:3"`
				B uint8 `struct:"uint8:2"`
				C uint8 `struct:"uint8"`
				D uint8 `struct:"uint8:3"`
			}{
				A: 0x02,
				B: 0x02,
				C: 0xAA,
				D: 0x05,
			},
		},
		{
			data: []byte{
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x02,
				0x03, 0x00, 0x00, 0x00,
			},
			bitsize: 96,
			value: struct {
				DefaultOrder uint32
				BigEndian    uint32 `struct:"big"`
				LittleEndian uint32 `struct:"little"`
			}{
				DefaultOrder: 1,
				BigEndian:    2,
				LittleEndian: 3,
			},
		},
		{
			data: []byte{
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x02,
				0x03, 0x00, 0x00, 0x00,
			},
			bitsize: 96,
			value: struct {
				DefaultOrder uint32
				BigSub       struct {
					BigEndian uint32
				} `struct:"big"`
				LittleSub struct {
					LittleEndian uint32
				} `struct:"little"`
			}{
				DefaultOrder: 1,
				BigSub:       struct{ BigEndian uint32 }{2},
				LittleSub:    struct{ LittleEndian uint32 }{3},
			},
		},
		{
			data: []byte{
				0x00, 0x00, 0x00, 0x02,
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x02,
				0x00, 0x00, 0x00, 0x03,
				0x00, 0x00, 0x00, 0x04,
			},
			bitsize: 160,
			value: struct {
				NumStructs int32 `struct:"sizeof=Structs"`
				Structs    []struct{ V1, V2 uint32 }
			}{
				NumStructs: 2,
				Structs: []struct{ V1, V2 uint32 }{
					{V1: 1, V2: 2},
					{V1: 3, V2: 4},
				},
			},
		},
		{
			data: []byte{
				0x3e, 0x00, 0x00, 0x00,
				0x3f, 0x80, 0x00, 0x00,
			},
			bitsize: 64,
			value: struct {
				C64 complex64
			}{
				C64: complex(0.125, 1.0),
			},
		},
		{
			data: []byte{
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x02,

				0x3f, 0x8c, 0xcc, 0xcd,
				0x3f, 0x99, 0x99, 0x9a,
				0x3f, 0xa6, 0x66, 0x66,
				0x3f, 0xb3, 0x33, 0x33,

				0x3f, 0xc0, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,

				0x3e, 0x00, 0x00, 0x00,
				0x3f, 0x80, 0x00, 0x00,

				0x3f, 0xc0, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x3f, 0xf0, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,

				0x3e, 0x00, 0x00, 0x00,
				0x3f, 0x80, 0x00, 0x00,

				0xfc, 0xfd, 0xfe, 0xff,
				0x00, 0x01, 0x02, 0x03,

				0xff, 0xfe, 0xfd, 0xfc,
				0xfb, 0xfa, 0xf9, 0xf8,

				0xff, 0xfe,

				0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0x00,

				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,

				0xe3, 0x82, 0x84, 0xe3,
				0x81, 0xa3, 0xe3, 0x81,
				0x9f, 0xef, 0xbc, 0x81,
			},
			bitsize: 880,
			value: struct {
				NumStructs uint32 `struct:"uint64,sizeof=Structs"`
				Structs    []struct{ V1, V2 float32 }
				Float64    float64
				Complex64  complex64
				Complex128 complex128
				Complex    complex128 `struct:"complex64"`
				SomeInt8s  [8]int8
				SomeUint8s [8]uint8
				AUint16    uint16
				AnInt64    int64
				_          [8]byte
				Message    string `struct:"[12]byte"`
			}{
				NumStructs: 2,
				Structs: []struct{ V1, V2 float32 }{
					{V1: 1.1, V2: 1.2},
					{V1: 1.3, V2: 1.4},
				},
				Float64:    0.125,
				Complex64:  complex(0.125, 1.0),
				Complex128: complex(0.125, 1.0),
				Complex:    complex(0.125, 1.0),
				SomeInt8s:  [8]int8{-4, -3, -2, -1, 0, 1, 2, 3},
				SomeUint8s: [8]uint8{0xff, 0xfe, 0xfd, 0xfc, 0xfb, 0xfa, 0xf9, 0xf8},
				AUint16:    0xfffe,
				AnInt64:    -256,
				Message:    "やった！",
			},
		},
		{
			data: []byte{
				0x00, 0x00, 0x00, 0x04,
				0xf0, 0x9f, 0x91, 0x8c,
			},
			bitsize: 64,
			value: struct {
				StrLen uint32 `struct:"uint32,sizeof=String"`
				String string
			}{
				StrLen: 4,
				String: "👌",
			},
		},
		{
			data: []byte{
				0x00, 0x00, 0x00, 0x04,
				0x00, 0x00, 0x00, 0x00,
				0xf0, 0x9f, 0x91, 0x8c,
			},
			bitsize: 96,
			value: struct {
				StrLen uint32 `struct:"uint32,sizeof=String"`
				String string `struct:"skip=4"`
			}{
				StrLen: 4,
				String: "👌",
			},
		},
		{
			data: []byte{
				0xf0, 0x9f, 0x91, 0x8c,
			},
			bitsize: 32,
			value: struct {
				String string `struct:"[4]byte"`
			}{
				String: "👌",
			},
		},
		{
			data: []byte{
				0xf0, 0x9f, 0x91, 0x8c, 0x00, 0x00, 0x00, 0x01,
			},
			bitsize: 64,
			value: struct {
				String string `struct:"[7]byte"`
				Value  byte
			}{
				String: "👌",
				Value:  1,
			},
		},
		{
			data: []byte{
				0x00, 0x02, 0x00,
				0x00, 0x00,
				0x00, 0x22, 0x18,
				0x00, 0x28, 0x12,
			},
			bitsize: 88,
			value: struct {
				Length int32 `struct:"int16,sizeof=Slice,little,skip=1"`
				Slice  []struct {
					Test int16 `struct:"skip=1"`
				} `struct:"skip=2,lsb"`
			}{
				Length: 2,
				Slice: []struct {
					Test int16 `struct:"skip=1"`
				}{
					{Test: 0x1822},
					{Test: 0x1228},
				},
			},
		},
		{
			data: []byte{
				0x00, 0x01,
				0x00, 0x02,
				0x00, 0x03,
			},
			bitsize: 48,
			value: struct {
				Ints []uint16 `struct:"[3]uint16"`
			}{
				Ints: []uint16{1, 2, 3},
			},
		},
		{
			data: []byte{
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x03,
			},
			bitsize: 64,
			value: struct {
				Size  int `struct:"int32,sizeof=Array"`
				Array []int32
			}{
				Size:  1,
				Array: []int32{3},
			},
		},
		{
			data: []byte{
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x03,
			},
			bitsize: 64,
			value: struct {
				_     struct{}
				Size  int `struct:"int32"`
				_     struct{}
				Array []int32 `struct:"sizefrom=Size"`
				_     struct{}
			}{
				Size:  1,
				Array: []int32{3},
			},
		},
		{
			data: []byte{
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x03,
				0x00, 0x00, 0x00, 0x04,
			},
			bitsize: 96,
			value: struct {
				_      struct{}
				Size   int `struct:"int32"`
				_      struct{}
				Array1 []int32 `struct:"sizefrom=Size"`
				Array2 []int32 `struct:"sizefrom=Size"`
				_      struct{}
			}{
				Size:   1,
				Array1: []int32{3},
				Array2: []int32{4},
			},
		},
		{
			data: []byte{
				0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff,
			},
			bitsize: 64,
			value: struct {
				A uint64 `struct:"uint64:12"`
				B uint64 `struct:"uint64:12"`
				C uint64 `struct:"uint64:30"`
				D uint64 `struct:"uint64:1"`
				E uint64 `struct:"uint64:5"`
				F uint64 `struct:"uint64:1"`
				G uint64 `struct:"uint64:3"`
			}{
				A: 0xfff,
				B: 0xfff,
				C: 0x3fffffff,
				D: 0x1,
				E: 0x1f,
				F: 0x1,
				G: 0x7,
			},
		},
		{
			data: []byte{
				// nonvariant/variant 8-bit
				// false, false, true, true
				0x00, 0x00, 0x01, 0xFF,

				// nonvariant/variant 8-bit inverted
				// false, false, true, true
				0x01, 0xFF, 0x00, 0x00,

				// nonvariant/variant 32-bit
				// false, false, true, true
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x01,
				0xFF, 0xFF, 0xFF, 0xFF,

				// nonvariant/variant 32-bit inverted
				// false, false, true, true
				0x00, 0x00, 0x00, 0x01,
				0xFF, 0xFF, 0xFF, 0xFF,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			bitsize: 320,
			value: struct {
				NonVariant8BitFalse          bool `struct:"bool"`
				Variant8BitFalse             bool `struct:"bool,variantbool"`
				NonVariant8BitTrue           bool `struct:"bool"`
				Variant8BitTrue              bool `struct:"bool,variantbool"`
				NonVariant8BitFalseInverted  bool `struct:"bool,invertedbool"`
				Variant8BitFalseInverted     bool `struct:"bool,invertedbool,variantbool"`
				NonVariant8BitTrueInverted   bool `struct:"bool,invertedbool"`
				Variant8BitTrueInverted      bool `struct:"bool,invertedbool,variantbool"`
				NonVariant32BitFalse         bool `struct:"int32"`
				Variant32BitFalse            bool `struct:"uint32,variantbool"`
				NonVariant32BitTrue          bool `struct:"uint32"`
				Variant32BitTrue             bool `struct:"int32,variantbool"`
				NonVariant32BitFalseInverted bool `struct:"uint32,invertedbool"`
				Variant32BitFalseInverted    bool `struct:"int32,invertedbool,variantbool"`
				NonVariant32BitTrueInverted  bool `struct:"int32,invertedbool"`
				Variant32BitTrueInverted     bool `struct:"uint32,invertedbool,variantbool"`
			}{
				NonVariant8BitFalse:          false,
				Variant8BitFalse:             false,
				NonVariant8BitTrue:           true,
				Variant8BitTrue:              true,
				NonVariant8BitFalseInverted:  false,
				Variant8BitFalseInverted:     false,
				NonVariant8BitTrueInverted:   true,
				Variant8BitTrueInverted:      true,
				NonVariant32BitFalse:         false,
				Variant32BitFalse:            false,
				NonVariant32BitTrue:          true,
				Variant32BitTrue:             true,
				NonVariant32BitFalseInverted: false,
				Variant32BitFalseInverted:    false,
				NonVariant32BitTrueInverted:  true,
				Variant32BitTrueInverted:     true,
			},
		},
		{
			data:    []byte{0x80},
			bitsize: 1,
			value: struct {
				Bit bool `struct:"uint8:1"`
			}{
				Bit: true,
			},
		},
		{
			data:    []byte{0x08, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F},
			bitsize: 80,
			value: struct {
				Count uint8 `struct:"uint8,sizeof=List"`
				List  []struct {
					A bool  `struct:"uint8:1,variantbool"`
					B uint8 `struct:"uint8:4"`
					C uint8 `struct:"uint8:4"`
				}
			}{
				Count: 8,
				List: []struct {
					A bool  `struct:"uint8:1,variantbool"`
					B uint8 `struct:"uint8:4"`
					C uint8 `struct:"uint8:4"`
				}{
					{A: false, B: 1, C: 14},
					{A: false, B: 3, C: 12},
					{A: false, B: 7, C: 8},
					{A: false, B: 15, C: 0},
					{A: true, B: 14, C: 1},
					{A: true, B: 12, C: 3},
					{A: true, B: 8, C: 7},
					{A: true, B: 0, C: 15},
				},
			},
		},
	}

	for _, test := range tests {
		v := reflect.New(reflect.TypeOf(test.value))

		// Test unpacking
		err := Unpack(test.data, binary.BigEndian, v.Interface())
		assert.Nil(t, err)
		assert.Equal(t, test.value, v.Elem().Interface())

		// Test packing
		data, err := Pack(binary.BigEndian, v.Interface())
		assert.Nil(t, err)
		assert.Equal(t, test.data, data)

		// Test sizing
		size, err := SizeOf(v.Interface())
		assert.Nil(t, err)
		assert.Equal(t, len(test.data), size)

		// Test bit sizing
		bits, err := BitSize(v.Interface())
		assert.Nil(t, err)
		assert.Equal(t, test.bitsize, bits)
	}
}

func TestUnpackBrokenSizeOf(t *testing.T) {
	data := []byte{
		0x00, 0x02, 0x00,
		0x00, 0x00,
		0x00, 0x22, 0x18,
		0x00, 0x28, 0x12,
	}

	s := struct {
		Length string `struct:"sizeof=Slice,skip=1"`
		Slice  []struct {
			Test int16 `struct:"skip=1"`
		} `struct:"skip=2,lsb"`
	}{
		Length: "no",
		Slice: []struct {
			Test int16 `struct:"skip=1"`
		}{
			{Test: 0x1822},
			{Test: 0x1228},
		},
	}

	// Test unpacking
	err := Unpack(data, binary.BigEndian, &s)
	assert.NotNil(t, err)
	assert.Equal(t, "unsupported size type string: Length", err.Error())

	// Test packing
	_, err = Pack(binary.BigEndian, &s)
	assert.NotNil(t, err)
	assert.Equal(t, "unsupported size type string: Length", err.Error())

	// Test unpacking sizeof to array fails.
	s2 := struct {
		Length int32    `struct:"sizeof=Array,skip=1"`
		Array  [2]int16 `struct:"skip=2,lsb"`
	}{
		Length: 2,
		Array: [2]int16{
			0x1822,
			0x1228,
		},
	}

	err = Unpack(data, binary.BigEndian, &s2)
	assert.NotNil(t, err)
	assert.Equal(t, "unsupported size target [2]int16", err.Error())
}

func TestUnpackBrokenArray(t *testing.T) {
	data := []byte{
		0x02, 0x00,
	}

	s := struct {
		Length int16 `struct:"[2]uint8"`
	}{
		Length: 2,
	}

	// Test unpacking
	err := Unpack(data, binary.BigEndian, &s)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid array cast type: int16", err.Error())

	// Test packing
	_, err = Pack(binary.BigEndian, &s)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid array cast type: int16", err.Error())

	s2 := struct {
		Length int16 `struct:"[]uint8"`
	}{
		Length: 2,
	}

	// Test unpacking
	err = Unpack(data, binary.BigEndian, &s2)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid array cast type: int16", err.Error())

	// Test packing
	_, err = Pack(binary.BigEndian, &s2)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid array cast type: int16", err.Error())
}

func TestUnpackFastPath(t *testing.T) {
	v := struct {
		Size uint8 `struct:"sizeof=Data"`
		Data []byte
	}{}
	Unpack([]byte("\x04Data"), binary.LittleEndian, &v)
	assert.Equal(t, 4, int(v.Size))
	assert.Equal(t, "Data", string(v.Data))
}

func BenchmarkFastPath(b *testing.B) {
	v := struct {
		Size uint8 `struct:"sizeof=Data"`
		Data []byte
	}{}
	data := []byte(" @?>=<;:9876543210/.-,+*)('&%$#\"! ")
	for i := 0; i < b.N; i++ {
		Unpack(data, binary.LittleEndian, &v)
	}
}

// Test custom packing
type CString string

func (s *CString) SizeOf() int {
	return len(*s) + 1
}

func (s *CString) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	for i, l := 0, len(buf); i < l; i++ {
		if buf[i] == 0 {
			*s = CString(buf[:i])
			return buf[i+1:], nil
		}
	}
	return []byte{}, errors.New("unterminated string")
}

func (s *CString) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	l := len(*s)
	for i := 0; i < l; i++ {
		buf[i] = (*s)[i]
	}
	buf[l] = 0
	return buf[l+1:], nil
}

func TestCustomPacking(t *testing.T) {
	x := CString("Test string! テスト。")
	b, err := Pack(binary.LittleEndian, &x)
	assert.Nil(t, err)
	assert.Equal(t, []byte{
		0x54, 0x65, 0x73, 0x74, 0x20, 0x73, 0x74, 0x72,
		0x69, 0x6e, 0x67, 0x21, 0x20, 0xe3, 0x83, 0x86,
		0xe3, 0x82, 0xb9, 0xe3, 0x83, 0x88, 0xe3, 0x80,
		0x82, 0x0,
	}, b)

	y := CString("")
	err = Unpack(b, binary.LittleEndian, &y)
	assert.Nil(t, err)
	assert.Equal(t, "Test string! テスト。", string(y))
}

// Test custom packing with non-pointer receiver
type Custom struct {
	A *int
}

func (s Custom) SizeOf() int {
	return 4
}

func (s Custom) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	*s.A = int(order.Uint32(buf[0:4]))
	return buf[4:], nil
}

func (s Custom) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	order.PutUint32(buf[0:4], uint32(*s.A))
	return buf[4:], nil
}

func TestCustomPackingNonPointer(t *testing.T) {
	c := Custom{new(int)}
	*c.A = 32

	b, err := Pack(binary.LittleEndian, c)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x20, 0x00, 0x00, 0x00}, b)

	d := Custom{new(int)}
	err = Unpack(b, binary.LittleEndian, d)
	assert.Nil(t, err)
	assert.Equal(t, 32, *d.A)
}
