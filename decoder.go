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
	order      binary.ByteOrder
	buf        []byte
	struc      reflect.Value
	sfields    []field
	bitCounter uint8
}

func (d *decoder) readBits(f field, outputLength uint8) (output []byte) {

	output = make([]byte, outputLength)

	if f.BitSize == 0 {
		// Having problems with complex64 type ... so we asume we want to read all
		// f.BitSize = uint8(f.Type.Bits())
		f.BitSize = 8 * outputLength
	}

	// originPos: Original position of the first bit in the first byte
	var originPos uint8 = 8 - d.bitCounter

	// destPos: Destination position ( in the result ) of the first bit in the first byte
	var destPos uint8 = f.BitSize % 8
	if destPos == 0 {
		destPos = 8
	}

	// numBytes: number of complete bytes to hold the result
	var numBytes uint8 = f.BitSize / 8

	// numBits: number of remaining bits in the first non-complete byte of the result
	var numBits uint8 = f.BitSize % 8

	// number of positions we have to shift the bytes to get the result
	var shift uint8 = (uint8(math.Abs(float64(originPos - destPos)))) % 8

	var outputInitialIdx uint8 = outputLength - numBytes
	if numBits > 0 {
		outputInitialIdx = outputInitialIdx - 1
	}

	if originPos < destPos { // shift left
		var idx uint8 = 0
		for outIdx := outputInitialIdx; outIdx < outputLength; outIdx++ {
			// TODO: Control the number of bytes of d.buf ... we need to read ahead
			var carry uint8 = d.buf[idx+1] >> (8 - shift)
			output[outIdx] = (d.buf[idx] << shift) | carry
			idx++
		}
	} else { // originPos >= destPos => shift right
		var idx uint8 = 0

		// carry : is a little bit tricky in this case because of the first case
		// when idx == 0 and there is no carry at all
		carry := func(idx uint8) uint8 {
			if idx == 0 {
				return 0x00
			}
			return (d.buf[idx-1] << (8 - shift))
		}

		for outIdx := outputInitialIdx; outIdx < outputLength; outIdx++ {
			output[outIdx] = (d.buf[idx] >> shift) | carry(idx)
			idx++
		}
	}

	// here the output is calculated ... but the first byte may have some extra bits
	// therefore we apply a mask to erase those unaddressable bits
	output[outputInitialIdx] &= ((0x01 << destPos) - 1)

	// now we need to update the head of the incoming buffer and the bitCounter
	d.bitCounter = (d.bitCounter + f.BitSize) % 8

	// move the head to the next non-complete byte used
	headerUpdate := func() uint8 {
		if (d.bitCounter == 0) && ((f.BitSize % 8) != 0) {
			return (numBytes + 1)
		}
		return numBytes
	}

	d.buf = d.buf[headerUpdate():]

	return
}

func (d *decoder) read8(f field) uint8 {
	rawdata := d.readBits(f, 1)
	return uint8(rawdata[0])
}

func (d *decoder) read16(f field) uint16 {
	rawdata := d.readBits(f, 2)
	return d.order.Uint16(rawdata)
}

func (d *decoder) read32(f field) uint32 {
	rawdata := d.readBits(f, 4)
	return d.order.Uint32(rawdata)
}

func (d *decoder) read64(f field) uint64 {
	rawdata := d.readBits(f, 8)
	return d.order.Uint64(rawdata)
}

func (d *decoder) readS8(f field) int8 { return int8(d.read8(f)) }

func (d *decoder) readS16(f field) int16 { return int16(d.read16(f)) }

func (d *decoder) readS32(f field) int32 { return int32(d.read32(f)) }

func (d *decoder) readS64(f field) int64 { return int64(d.read64(f)) }

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
		v.SetInt(int64(d.readS8(f)))
	case reflect.Int16:
		v.SetInt(int64(d.readS16(f)))
	case reflect.Int32:
		v.SetInt(int64(d.readS32(f)))
	case reflect.Int64:
		v.SetInt(d.readS64(f))

	case reflect.Uint8:
		v.SetUint(uint64(d.read8(f)))
	case reflect.Uint16:
		v.SetUint(uint64(d.read16(f)))
	case reflect.Uint32:
		v.SetUint(uint64(d.read32(f)))
	case reflect.Uint64:
		v.SetUint(d.read64(f))

	case reflect.Float32:
		v.SetFloat(float64(math.Float32frombits(d.read32(f))))
	case reflect.Float64:
		v.SetFloat(math.Float64frombits(d.read64(f)))

	case reflect.Complex64:
		v.SetComplex(complex(
			float64(math.Float32frombits(d.read32(f))),
			float64(math.Float32frombits(d.read32(f))),
		))
	case reflect.Complex128:
		v.SetComplex(complex(
			math.Float64frombits(d.read64(f)),
			math.Float64frombits(d.read64(f)),
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
