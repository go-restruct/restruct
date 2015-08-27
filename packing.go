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

// Unpack reads data from a byteslice into a structure.
func Unpack(data []byte, order binary.ByteOrder, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	val := reflect.ValueOf(v).Elem()
	d := decoder{order: order, buf: data}
	d.read(FieldFromType(val.Type()), val)

	return
}

// Pack writes data from a datastructure into a byteslice.
func Pack(order binary.ByteOrder, v interface{}) (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			data = nil
			err = r.(error)
		}
	}()

	val := reflect.ValueOf(v).Elem()
	f := FieldFromType(val.Type())
	data = make([]byte, f.SizeOf(val))
	e := encoder{order: order, buf: data}
	e.write(f, val)

	return
}
