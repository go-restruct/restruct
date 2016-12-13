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
	order      binary.ByteOrder
	buf        []byte
	struc      reflect.Value
	sfields    []field
	bitCounter uint8
	realSize   uint64
	initBuf    []byte
}

func (e *encoder) assignBuffer(in []byte) {
	e.buf = in
	e.initBuf = in
}

func (e *encoder) writeBits(f field, inBuf []byte) {

	var inputLength uint8 = uint8(len(inBuf))

	if f.BitSize == 0 {
		// Having problems with complex64 type ... so we asume we want to read all
		//f.BitSize = uint8(f.Type.Bits())
		f.BitSize = 8 * inputLength
	}
	// destPos: Destination position ( in the result ) of the first bit in the first byte
	var destPos uint8 = 8 - e.bitCounter

	// originPos: Original position of the first bit in the first byte
	var originPos uint8 = f.BitSize % 8
	if originPos == 0 {
		originPos = 8
	}

	// numBytes: number of complete bytes to hold the result
	var numBytes uint8 = f.BitSize / 8

	// numBits: number of remaining bits in the first non-complete byte of the result
	var numBits uint8 = f.BitSize % 8

	// number of positions we have to shift the bytes to get the result
	var shift uint8
	if originPos > destPos {
		shift = originPos - destPos
	} else {
		shift = destPos - originPos
	}
	shift = shift % 8

	var inputInitialIdx uint8 = inputLength - numBytes
	if numBits > 0 {
		inputInitialIdx = inputInitialIdx - 1
	}

	if originPos < destPos {
		// shift left
		carry := func(idx uint8) uint8 {
			if (idx + 1) < inputLength {
				return (inBuf[idx + 1] >> (8 - shift))
			}
			return 0x00

		}
		mask := func(idx uint8) uint8{
			if idx == 0 {
				return (0x01<<destPos)-1
			}
			return 0xFF
		}
		var idx uint8 = 0
		for inIdx := inputInitialIdx; inIdx < inputLength; inIdx++ {
			e.buf[idx] |= ((inBuf[inIdx] << shift) | carry(inIdx) ) & mask(idx)
			idx++
		}

	} else {
		// originPos >= destPos => shift right
		var idx uint8 = 0
		// carry : is a little bit tricky in this case because of the first case
		// when idx == 0 and there is no carry at all
		carry := func(idx uint8) uint8 {
			if idx == 0 {
				return 0x00
			}
			return (inBuf[idx-1] << (8 - shift))
		}
		mask := func(idx uint8) uint8{
			if idx == 0 {
				return (0x01<<destPos)-1
			}
			return 0xFF
		}
		inIdx := inputInitialIdx
		for ; inIdx < inputLength; inIdx++ {
			//note: Should the mask be done BEFORE the OR with carry?
			e.buf[idx] |= ((inBuf[inIdx] >> shift) | carry(inIdx)) & mask(idx)
			idx++
		}
		if ((e.bitCounter + f.BitSize) % 8) > 0 {
			e.buf[idx] |= carry(inIdx)
		}
	}

        //now we should update buffer and bitCounter
	e.bitCounter = (e.bitCounter + f.BitSize) % 8

	// move the head to the next non-complete byte used
	headerUpdate := func() uint8 {
		if (e.bitCounter == 0) && ((f.BitSize % 8) != 0) {
			return (numBytes + 1)
		}
		return numBytes
	}
	e.buf = e.buf[headerUpdate():]
	return
}

func (e *encoder) write8(f field, x uint8) {
	typeSize := uint8(reflect.TypeOf(x).Size())

	b := make([]byte, typeSize)
	b[0] = x
	e.writeBits(f, b)
}

func (e *encoder) write16(f field, x uint16) {
	typeSize := uint8(reflect.TypeOf(x).Size())

	b := make([]byte, typeSize)
	e.order.PutUint16(b[0:typeSize], x)

	e.writeBits(f, b)
}

func (e *encoder) write32(f field, x uint32) {
	typeSize := uint8(reflect.TypeOf(x).Size())

	b := make([]byte, typeSize)
	e.order.PutUint32(b[0:typeSize], x)

	e.writeBits(f, b)
}

func (e *encoder) write64(f field, x uint64) {
	typeSize := uint8(reflect.TypeOf(x).Size())

	b := make([]byte, typeSize)
	e.order.PutUint64(b[0:typeSize], x)

	e.writeBits(f, b)
}

func (e *encoder) writeS8(f field, x int8) { e.write8(f, uint8(x)) }

func (e *encoder) writeS16(f field, x int16) { e.write16(f, uint16(x)) }

func (e *encoder) writeS32(f field, x int32) { e.write32(f, uint32(x)) }

func (e *encoder) writeS64(f field, x int64) { e.write64(f, uint64(x)) }

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
		e.writeS8(f, int8(v.Int()))
	case reflect.Int16:
		e.writeS16(f, int16(v.Int()))
	case reflect.Int32:
		e.writeS32(f, int32(v.Int()))
	case reflect.Int64:
		e.writeS64(f, int64(v.Int()))

	case reflect.Uint8:
		e.write8(f, uint8(v.Uint()))
	case reflect.Uint16:
		e.write16(f, uint16(v.Uint()))
	case reflect.Uint32:
		e.write32(f, uint32(v.Uint()))
	case reflect.Uint64:
		e.write64(f, uint64(v.Uint()))

	case reflect.Float32:
		e.write32(f, math.Float32bits(float32(v.Float())))
	case reflect.Float64:
		e.write64(f, math.Float64bits(float64(v.Float())))

	case reflect.Complex64:
		x := v.Complex()
		e.write32(f, math.Float32bits(float32(real(x))))
		e.write32(f, math.Float32bits(float32(imag(x))))
	case reflect.Complex128:
		x := v.Complex()
		e.write64(f, math.Float64bits(float64(real(x))))
		e.write64(f, math.Float64bits(float64(imag(x))))
	}
}
