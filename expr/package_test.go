package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	pkg := NewPackage(map[string]Value{
		"A": ValueOf(1),
		"B": ValueOf(int32(2)),
	})
	assert.Equal(t, 1, pkg.Symbol("A").RawValue())
	assert.Equal(t, int32(2), pkg.Symbol("B").RawValue())
	assert.Equal(t, nil, pkg.Symbol("C"))
	assert.Equal(t, Int, pkg.Type().(*PackageType).Symbol("A").Kind())
	assert.Equal(t, Int32, pkg.Type().(*PackageType).Symbol("B").Kind())
	assert.Equal(t, nil, pkg.Type().(*PackageType).Symbol("C"))
}
