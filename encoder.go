package restruct

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

// Packer is a type capable of packing a native value into a binary
// representation. The Pack function is expected to overwrite a number of
// bytes in buf then return a slice of the remaining buffer. Note that you
// must also implement SizeOf, and returning an incorrect SizeOf will cause
// the encoder to crash. The SizeOf should be equal to the number of bytes
// consumed from the buffer slice in Pack. You may use a pointer receiver even
// if the type is used by value.
type Packer interface {
	Sizer
	Pack(buf []byte, order binary.ByteOrder) ([]byte, error)
}

type encoder struct {
	order   binary.ByteOrder
	buf     []byte
	struc   reflect.Value
	sfields []field
}

func (e *encoder) write8(x uint8) {
	e.buf[0] = x
	e.buf = e.buf[1:]
}

func (e *encoder) write16(x uint16) {
	e.order.PutUint16(e.buf[0:2], x)
	e.buf = e.buf[2:]
}

func (e *encoder) write32(x uint32) {
	e.order.PutUint32(e.buf[0:4], x)
	e.buf = e.buf[4:]
}

func (e *encoder) write64(x uint64) {
	e.order.PutUint64(e.buf[0:8], x)
	e.buf = e.buf[8:]
}

func (e *encoder) writeS8(x int8) { e.write8(uint8(x)) }

func (e *encoder) writeS16(x int16) { e.write16(uint16(x)) }

func (e *encoder) writeS32(x int32) { e.write32(uint32(x)) }

func (e *encoder) writeS64(x int64) { e.write64(uint64(x)) }

func (e *encoder) skipn(count int) {
	e.buf = e.buf[count:]
}

func (e *encoder) skip(f field, v reflect.Value) {
	e.skipn(f.SizeOf(v))
}

func (e *encoder) packer(v reflect.Value) (Packer, bool) {
	if s, ok := v.Interface().(Packer); ok {
		return s, true
	}

	if !v.CanAddr() {
		return nil, false
	}

	if s, ok := v.Addr().Interface().(Packer); ok {
		return s, true
	}

	return nil, false
}

func (e *encoder) write(f field, v reflect.Value) {
	if f.Name != "_" {
		if s, ok := e.packer(v); ok {
			var err error
			e.buf, err = s.Pack(e.buf, e.order)
			if err != nil {
				panic(err)
			}
			return
		}
	} else {
		e.skipn(f.SizeOf(v))
		return
	}

	struc := e.struc
	sfields := e.sfields
	order := e.order

	if f.Order != nil {
		e.order = f.Order
		defer func() { e.order = order }()
	}

	if f.Skip != 0 {
		e.skipn(f.Skip)
	}

	// If this is a sizeof field, pull the current slice length into it.
	if f.SIndex != -1 {
		sv := struc.Field(f.SIndex)

		switch f.DefType.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(int64(sv.Len()))
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(uint64(sv.Len()))
		default:
			panic(fmt.Errorf("unsupported sizeof type %s", f.DefType.String()))
		}
	}

	switch f.Type.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		switch f.DefType.Kind() {
		case reflect.Array, reflect.Slice, reflect.String:
			ef := f.Elem()
			l := v.Len()
			for i := 0; i < l; i++ {
				e.write(ef, v.Index(i))
			}
		default:
			panic(fmt.Errorf("invalid array cast type: %s", f.DefType.String()))
		}

	case reflect.Struct:
		e.struc = v
		e.sfields = cachedFieldsFromStruct(f.Type)
		l := len(e.sfields)
		for i := 0; i < l; i++ {
			f := e.sfields[i]
			v := v.Field(f.Index)
			if v.CanSet() {
				e.write(f, v)
			} else {
				e.skip(f, v)
			}
		}
		e.sfields = sfields
		e.struc = struc

	case reflect.Int8:
		e.writeS8(int8(v.Int()))
	case reflect.Int16:
		e.writeS16(int16(v.Int()))
	case reflect.Int32:
		e.writeS32(int32(v.Int()))
	case reflect.Int64:
		e.writeS64(int64(v.Int()))

	case reflect.Uint8:
		e.write8(uint8(v.Uint()))
	case reflect.Uint16:
		e.write16(uint16(v.Uint()))
	case reflect.Uint32:
		e.write32(uint32(v.Uint()))
	case reflect.Uint64:
		e.write64(uint64(v.Uint()))

	case reflect.Float32:
		e.write32(math.Float32bits(float32(v.Float())))
	case reflect.Float64:
		e.write64(math.Float64bits(float64(v.Float())))

	case reflect.Complex64:
		x := v.Complex()
		e.write32(math.Float32bits(float32(real(x))))
		e.write32(math.Float32bits(float32(imag(x))))
	case reflect.Complex128:
		x := v.Complex()
		e.write64(math.Float64bits(float64(real(x))))
		e.write64(math.Float64bits(float64(imag(x))))
	}
}
