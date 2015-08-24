package restruct

import (
	"encoding/binary"
	"reflect"
)

// Unpack reads data from a byteslice into a structure.
func Unpack(data []byte, order binary.ByteOrder, v interface{}) (err error) {
	/*defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()*/

	val := reflect.ValueOf(v).Elem()
	d := decoder{order: order, buf: data}
	d.read(FieldFromType(val.Type()), val)

	return nil
}
