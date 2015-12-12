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
		data  []byte
		value interface{}
	}{
		{
			data: []byte{
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x02,
				0x03, 0x00, 0x00, 0x00,
			},
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
			value: struct {
				String string `struct:"[4]byte"`
			}{
				String: "👌",
			},
		},
		{
			data: []byte{
				0x00, 0x02, 0x00,
				0x00, 0x00,
				0x00, 0x22, 0x18,
				0x00, 0x28, 0x12,
			},
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
			value: struct {
				Ints []uint16 `struct:"[3]uint16"`
			}{
				Ints: []uint16{1, 2, 3},
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
	assert.Equal(t, "unsupported sizeof type string", err.Error())

	// Test packing
	_, err = Pack(binary.BigEndian, &s)
	assert.NotNil(t, err)
	assert.Equal(t, "unsupported sizeof type string", err.Error())

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
	assert.Equal(t, "unsupported sizeof target [2]int16", err.Error())
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

// Test custom packing with non-pointer reciever
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
