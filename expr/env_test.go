package expr

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructResolver(t *testing.T) {
	v := struct {
		A int
		B int32
	}{
		A: 1,
		B: 2,
	}
	resolver := NewStructResolver(reflect.ValueOf(v))
	typeResolver := NewStructTypeResolver(v)
	assert.Equal(t, 1, resolver.Resolve("A").RawValue())
	assert.Equal(t, int32(2), resolver.Resolve("B").RawValue())
	assert.Equal(t, nil, resolver.Resolve("C"))
	assert.Equal(t, Int, typeResolver.TypeResolve("A").Kind())
	assert.Equal(t, Int32, typeResolver.TypeResolve("B").Kind())
	assert.Equal(t, nil, typeResolver.TypeResolve("C"))
}

func TestMapResolver(t *testing.T) {
	resolver := NewMapResolver(map[string]Value{
		"A": ValueOf(1),
		"B": ValueOf(int32(2)),
	})
	typeResolver := NewMapTypeResolver(map[string]Type{
		"A": TypeOf(1),
		"B": TypeOf(int32(2)),
	})
	assert.Equal(t, 1, resolver.Resolve("A").RawValue())
	assert.Equal(t, int32(2), resolver.Resolve("B").RawValue())
	assert.Equal(t, nil, resolver.Resolve("C"))
	assert.Equal(t, Int, typeResolver.TypeResolve("A").Kind())
	assert.Equal(t, Int32, typeResolver.TypeResolve("B").Kind())
	assert.Equal(t, nil, typeResolver.TypeResolve("C"))
}

func TestTypeResolverAdapter(t *testing.T) {
	typeResolver := NewTypeResolverAdapter(NewMapResolver(map[string]Value{
		"A": ValueOf(1),
		"B": ValueOf(int32(2)),
	}))
	assert.Equal(t, Int, typeResolver.TypeResolve("A").Kind())
	assert.Equal(t, Int32, typeResolver.TypeResolve("B").Kind())
	assert.Equal(t, nil, typeResolver.TypeResolve("C"))
}

func TestMetaResolver(t *testing.T) {
	v := struct{ A int }{A: 1}
	resolver := NewMetaResolver()
	typeResolver := NewMetaTypeResolver()
	resolver.AddResolver(NewStructResolver(reflect.ValueOf(v)))
	resolver.AddResolver(NewMapResolver(map[string]Value{"B": ValueOf(int32(2))}))
	typeResolver.AddResolver(NewStructTypeResolver(v))
	typeResolver.AddResolver(NewMapTypeResolver(map[string]Type{"B": TypeOf(int32(2))}))
	assert.Equal(t, 1, resolver.Resolve("A").RawValue())
	assert.Equal(t, int32(2), resolver.Resolve("B").RawValue())
	assert.Equal(t, nil, resolver.Resolve("C"))
	assert.Equal(t, Int, typeResolver.TypeResolve("A").Kind())
	assert.Equal(t, Int32, typeResolver.TypeResolve("B").Kind())
	assert.Equal(t, nil, typeResolver.TypeResolve("C"))
}
