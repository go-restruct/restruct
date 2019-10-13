package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueBasicCompare(t *testing.T) {
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(int(10)).Add(ValueOf(int(20))).Equal(ValueOf(int(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(int8(10)).Add(ValueOf(int8(20))).Equal(ValueOf(int8(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(int16(10)).Add(ValueOf(int16(20))).Equal(ValueOf(int16(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(int32(10)).Add(ValueOf(int32(20))).Equal(ValueOf(int32(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(int64(10)).Add(ValueOf(int64(20))).Equal(ValueOf(int64(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(uint(10)).Add(ValueOf(uint(20))).Equal(ValueOf(uint(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(uint8(10)).Add(ValueOf(uint8(20))).Equal(ValueOf(uint8(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(uint16(10)).Add(ValueOf(uint16(20))).Equal(ValueOf(uint16(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(uint32(10)).Add(ValueOf(uint32(20))).Equal(ValueOf(uint32(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(uint64(10)).Add(ValueOf(uint64(20))).Equal(ValueOf(uint64(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(uintptr(10)).Add(ValueOf(uintptr(20))).Equal(ValueOf(uintptr(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(float32(10)).Add(ValueOf(float32(20))).Equal(ValueOf(float32(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(float64(10)).Add(ValueOf(float64(20))).Equal(ValueOf(float64(30))).RawValue())

	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(int(15)).Add(ValueOf(int(20))).Equal(ValueOf(int(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(int8(15)).Add(ValueOf(int8(20))).Equal(ValueOf(int8(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(int16(15)).Add(ValueOf(int16(20))).Equal(ValueOf(int16(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(int32(15)).Add(ValueOf(int32(20))).Equal(ValueOf(int32(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(int64(15)).Add(ValueOf(int64(20))).Equal(ValueOf(int64(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(uint(15)).Add(ValueOf(uint(20))).Equal(ValueOf(uint(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(uint8(15)).Add(ValueOf(uint8(20))).Equal(ValueOf(uint8(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(uint16(15)).Add(ValueOf(uint16(20))).Equal(ValueOf(uint16(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(uint32(15)).Add(ValueOf(uint32(20))).Equal(ValueOf(uint32(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(uint64(15)).Add(ValueOf(uint64(20))).Equal(ValueOf(uint64(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(uintptr(15)).Add(ValueOf(uintptr(20))).Equal(ValueOf(uintptr(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(float32(15)).Add(ValueOf(float32(20))).Equal(ValueOf(float32(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(float64(15)).Add(ValueOf(float64(20))).Equal(ValueOf(float64(30))).RawValue())

	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(int(10)).Add(ValueOf(int(20))).NotEqual(ValueOf(int(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(int8(10)).Add(ValueOf(int8(20))).NotEqual(ValueOf(int8(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(int16(10)).Add(ValueOf(int16(20))).NotEqual(ValueOf(int16(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(int32(10)).Add(ValueOf(int32(20))).NotEqual(ValueOf(int32(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(int64(10)).Add(ValueOf(int64(20))).NotEqual(ValueOf(int64(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(uint(10)).Add(ValueOf(uint(20))).NotEqual(ValueOf(uint(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(uint8(10)).Add(ValueOf(uint8(20))).NotEqual(ValueOf(uint8(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(uint16(10)).Add(ValueOf(uint16(20))).NotEqual(ValueOf(uint16(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(uint32(10)).Add(ValueOf(uint32(20))).NotEqual(ValueOf(uint32(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(uint64(10)).Add(ValueOf(uint64(20))).NotEqual(ValueOf(uint64(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(uintptr(10)).Add(ValueOf(uintptr(20))).NotEqual(ValueOf(uintptr(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(float32(10)).Add(ValueOf(float32(20))).NotEqual(ValueOf(float32(30))).RawValue())
	assert.Equal(t, ValueOf(false).RawValue(), ValueOf(float64(10)).Add(ValueOf(float64(20))).NotEqual(ValueOf(float64(30))).RawValue())

	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(int(15)).Add(ValueOf(int(20))).NotEqual(ValueOf(int(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(int8(15)).Add(ValueOf(int8(20))).NotEqual(ValueOf(int8(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(int16(15)).Add(ValueOf(int16(20))).NotEqual(ValueOf(int16(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(int32(15)).Add(ValueOf(int32(20))).NotEqual(ValueOf(int32(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(int64(15)).Add(ValueOf(int64(20))).NotEqual(ValueOf(int64(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(uint(15)).Add(ValueOf(uint(20))).NotEqual(ValueOf(uint(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(uint8(15)).Add(ValueOf(uint8(20))).NotEqual(ValueOf(uint8(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(uint16(15)).Add(ValueOf(uint16(20))).NotEqual(ValueOf(uint16(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(uint32(15)).Add(ValueOf(uint32(20))).NotEqual(ValueOf(uint32(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(uint64(15)).Add(ValueOf(uint64(20))).NotEqual(ValueOf(uint64(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(uintptr(15)).Add(ValueOf(uintptr(20))).NotEqual(ValueOf(uintptr(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(float32(15)).Add(ValueOf(float32(20))).NotEqual(ValueOf(float32(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), ValueOf(float64(15)).Add(ValueOf(float64(20))).NotEqual(ValueOf(float64(30))).RawValue())
}

func TestCoerceUntyped(t *testing.T) {
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(int(20))).Equal(ValueOf(int(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(int8(20))).Equal(ValueOf(int8(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(int16(20))).Equal(ValueOf(int16(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(int32(20))).Equal(ValueOf(int32(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(int64(20))).Equal(ValueOf(int64(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(uint(20))).Equal(ValueOf(uint(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(uint8(20))).Equal(ValueOf(uint8(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(uint16(20))).Equal(ValueOf(uint16(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(uint32(20))).Equal(ValueOf(uint32(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(uint64(20))).Equal(ValueOf(uint64(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(uintptr(20))).Equal(ValueOf(uintptr(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(float32(20))).Equal(ValueOf(float32(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literalintval(10).Add(ValueOf(float64(20))).Equal(ValueOf(float64(30))).RawValue())

	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(int(20))).Equal(ValueOf(int(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(int8(20))).Equal(ValueOf(int8(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(int16(20))).Equal(ValueOf(int16(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(int32(20))).Equal(ValueOf(int32(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(int64(20))).Equal(ValueOf(int64(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(uint(20))).Equal(ValueOf(uint(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(uint8(20))).Equal(ValueOf(uint8(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(uint16(20))).Equal(ValueOf(uint16(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(uint32(20))).Equal(ValueOf(uint32(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(uint64(20))).Equal(ValueOf(uint64(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(uintptr(20))).Equal(ValueOf(uintptr(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(float32(20))).Equal(ValueOf(float32(30))).RawValue())
	assert.Equal(t, ValueOf(true).RawValue(), literaluintval(10).Add(ValueOf(float64(20))).Equal(ValueOf(float64(30))).RawValue())
}

func TestUnaryOps(t *testing.T) {
	i := 35
	s := struct{ I int }{I: 35}
	assert.Equal(t, ValueOf(-15).RawValue(), ValueOf(15).Negate().RawValue())
	assert.Equal(t, ValueOf(!true).RawValue(), ValueOf(true).Not().RawValue())
	assert.Equal(t, ValueOf(!false).RawValue(), ValueOf(false).Not().RawValue())
	assert.Equal(t, ValueOf(^25).RawValue(), ValueOf(25).BitNot().RawValue())
	assert.Equal(t, ValueOf(35).RawValue(), ValueOf(&i).Deref().RawValue())
	assert.Equal(t, ValueOf(35).RawValue(), ValueOf(&s).Dot("I").Ref().Deref().RawValue())
}
