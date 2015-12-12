package restruct

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

// Unpacker is a type capable of unpacking a binary representation of itself
// into a native representation. The Unpack function is expected to consume
// a number of bytes from the buffer, then return a slice of the remaining
// bytes in the buffer. You may use a pointer receiver even if the type is
// used by value.
type Unpacker interface {
	Unpack(buf []byte, order binary.ByteOrder) ([]byte, error)
}

type decoder struct {
	order   binary.ByteOrder
	buf     []byte
	struc   reflect.Value
	sfields []field
}

func (d *decoder) read8() uint8 {
	x := d.buf[0]
	d.buf = d.buf[1:]
	return x
}

func (d *decoder) read16() uint16 {
	x := d.order.Uint16(d.buf[0:2])
	d.buf = d.buf[2:]
	return x
}

func (d *decoder) read32() uint32 {
	x := d.order.Uint32(d.buf[0:4])
	d.buf = d.buf[4:]
	return x
}

func (d *decoder) read64() uint64 {
	x := d.order.Uint64(d.buf[0:8])
	d.buf = d.buf[8:]
	return x
}

func (d *decoder) readS8() int8 { return int8(d.read8()) }

func (d *decoder) readS16() int16 { return int16(d.read16()) }

func (d *decoder) readS32() int32 { return int32(d.read32()) }

func (d *decoder) readS64() int64 { return int64(d.read64()) }

func (d *decoder) readn(count int) []byte {
	x := d.buf[0:count]
	d.buf = d.buf[count:]
	return x
}

func (d *decoder) skipn(count int) {
	d.buf = d.buf[count:]
}

func (d *decoder) skip(f field, v reflect.Value) {
	d.skipn(f.SizeOf(v))
}

func (d *decoder) unpacker(v reflect.Value) (Unpacker, bool) {
	if s, ok := v.Interface().(Unpacker); ok {
		return s, true
	}

	if !v.CanAddr() {
		return nil, false
	}

	if s, ok := v.Addr().Interface().(Unpacker); ok {
		return s, true
	}

	return nil, false
}

func (d *decoder) read(f field, v reflect.Value) {
	if f.Name != "_" {
		if s, ok := d.unpacker(v); ok {
			var err error
			d.buf, err = s.Unpack(d.buf, d.order)
			if err != nil {
				panic(err)
			}
			return
		}
	} else {
		d.skipn(f.SizeOf(v))
		return
	}

	struc := d.struc
	sfields := d.sfields
	order := d.order

	if f.Order != nil {
		d.order = f.Order
		defer func() { d.order = order }()
	}

	if f.Skip != 0 {
		d.skipn(f.Skip)
	}

	switch f.Type.Kind() {
	case reflect.Array:
		l := f.Type.Len()

		// If the underlying value is a slice, initialize it.
		if f.DefType.Kind() == reflect.Slice {
			v.Set(reflect.MakeSlice(reflect.SliceOf(f.Type.Elem()), l, l))
		}

		switch f.DefType.Kind() {
		case reflect.String:
			v.SetString(string(d.readn(f.SizeOf(v))))
		case reflect.Slice, reflect.Array:
			ef := f.Elem()
			for i := 0; i < l; i++ {
				d.read(ef, v.Index(i))
			}
		default:
			panic(fmt.Errorf("invalid array cast type: %s", f.DefType.String()))
		}

	case reflect.Struct:
		d.struc = v
		d.sfields = cachedFieldsFromStruct(f.Type)
		l := len(d.sfields)
		for i := 0; i < l; i++ {
			f := d.sfields[i]
			v := v.Field(f.Index)
			if v.CanSet() {
				d.read(f, v)
			} else {
				d.skip(f, v)
			}
		}
		d.sfields = sfields
		d.struc = struc

	case reflect.Slice, reflect.String:
		switch f.DefType.Kind() {
		case reflect.String:
			l := v.Len()
			v.SetString(string(d.readn(l)))
		case reflect.Slice, reflect.Array:
			switch f.DefType.Elem().Kind() {
			case reflect.Uint8:
				v.SetBytes(d.readn(f.SizeOf(v)))
			default:
				l := v.Len()
				ef := f.Elem()
				for i := 0; i < l; i++ {
					d.read(ef, v.Index(i))
				}
			}
		default:
			panic(fmt.Errorf("invalid array cast type: %s", f.DefType.String()))
		}

	case reflect.Int8:
		v.SetInt(int64(d.readS8()))
	case reflect.Int16:
		v.SetInt(int64(d.readS16()))
	case reflect.Int32:
		v.SetInt(int64(d.readS32()))
	case reflect.Int64:
		v.SetInt(d.readS64())

	case reflect.Uint8:
		v.SetUint(uint64(d.read8()))
	case reflect.Uint16:
		v.SetUint(uint64(d.read16()))
	case reflect.Uint32:
		v.SetUint(uint64(d.read32()))
	case reflect.Uint64:
		v.SetUint(d.read64())

	case reflect.Float32:
		v.SetFloat(float64(math.Float32frombits(d.read32())))
	case reflect.Float64:
		v.SetFloat(math.Float64frombits(d.read64()))

	case reflect.Complex64:
		v.SetComplex(complex(
			float64(math.Float32frombits(d.read32())),
			float64(math.Float32frombits(d.read32())),
		))
	case reflect.Complex128:
		v.SetComplex(complex(
			math.Float64frombits(d.read64()),
			math.Float64frombits(d.read64()),
		))
	}

	if f.SIndex != -1 {
		sv := struc.Field(f.SIndex)
		l := len(sfields)
		for i := 0; i < l; i++ {
			if sfields[i].Index != f.SIndex {
				continue
			}

			sf := sfields[i]
			sl := 0

			// Must use different codepath for signed/unsigned.
			switch f.DefType.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				sl = int(v.Int())
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				sl = int(v.Uint())
			default:
				panic(fmt.Errorf("unsupported sizeof type %s", f.DefType.String()))
			}

			// Strings are immutable, but we make a blank one so that we can
			// figure out the size later. It might be better to do something
			// more hackish, like writing the length into the string...
			switch sf.DefType.Kind() {
			case reflect.Slice:
				sv.Set(reflect.MakeSlice(sf.Type, sl, sl))
			case reflect.String:
				sv.SetString(string(make([]byte, sl)))
			default:
				panic(fmt.Errorf("unsupported sizeof target %s", sf.DefType.String()))
			}
		}
	}
}
