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
