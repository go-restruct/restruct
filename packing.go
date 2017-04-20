/*
Package restruct implements packing and unpacking of raw binary formats.

Structures can be created with struct tags annotating the on-disk or in-memory
layout of the structure, using the "struct" struct tag, like so:

	struct {
		Length int `struct:"int32,sizeof=Packets"`
		Packets []struct{
			Source    string    `struct:"[16]byte"`
			Timestamp int       `struct:"int32,big"`
			Data      [256]byte `struct:"skip=8"`
		}
	}

To unpack data in memory to this structure, simply use Unpack with a byte slice:

	msg := Message{}
	restruct.Unpack(data, binary.LittleEndian, &msg)
*/
package restruct

import (
	"encoding/binary"
	"reflect"
)

/*
Unpack reads data from a byteslice into a value.

Two types of values are directly supported here: Unpackers and structs. You can
pass them by value or by pointer, although it is an error if Restruct is
unable to set a value because it is unaddressable.

For structs, each field will be read sequentially based on a straightforward
interpretation of the type. For example, an int32 will be read as a 32-bit
signed integer, taking 4 bytes of memory. Structures and arrays are laid out
flat with no padding or metadata.

Unexported fields are ignored, except for fields named _ - those fields will
be treated purely as padding. Padding will not be preserved through packing
and unpacking.

The behavior of deserialization can be customized using struct tags. The
following struct tag syntax is supported:

	`struct:"[flags...]"`

Flags are comma-separated keys. The following are available:

	type            A bare type name, e.g. int32 or []string.

	sizeof=[Field]  Specifies that the field should be treated as a count of
	                the number of elements in Field.

	skip=[Count]    Skips Count bytes before the field. You can use this to
	                e.g. emulate C structure alignment.

	big,msb         Specifies big endian byte order. When applied to structs,
	                this will apply to all fields under the struct.

	little,lsb      Specifies little endian byte order. When applied to structs,
	                this will apply to all fields under the struct.
*/
func Unpack(data []byte, order binary.ByteOrder, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	d := decoder{order: order, buf: data}
	d.read(fieldFromType(val.Type()), val)

	return
}

/*
Pack writes data from a datastructure into a byteslice.

Two types of values are directly supported here: Packers and structs. You can
pass them by value or by pointer.

Each structure is serialized in the same way it would be deserialized with
Unpack. See Unpack documentation for the struct tag format.
*/
func Pack(order binary.ByteOrder, v interface{}) (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			data = nil
			err = r.(error)
		}
	}()

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	f := fieldFromType(val.Type())
	data = make([]byte, f.SizeOf(val))

	e := encoder{order: order}
	e.assignBuffer(data)
	e.write(f, val)

	data = e.initBuf

	return
}
