// +build !go1.5

package restruct

import (
	"fmt"
	"reflect"
	"sync"
)

type arrayTypeKey struct {
	count int
	elem  reflect.Type
}

var arrayTypeMap = map[arrayTypeKey]reflect.Type{}
var arrayTypeMutex = sync.RWMutex{}

// RegisterArrayType registers an array type for use with struct tags. This
// function is only necessary on Go 1.4 and below, which do not provide
// reflect.ArrayOf. This function is goroutine safe and idempotent (i.e.,
// calling it multiple times with the same value is perfectly safe.) If you
// require Go 1.5 or above, you can safely use any array type in your struct
// tags without using this function.
func RegisterArrayType(array interface{}) {
	typ := reflect.TypeOf(array)
	key := arrayTypeKey{
		count: typ.Len(),
		elem:  typ.Elem(),
	}

	arrayTypeMutex.Lock()
	arrayTypeMap[key] = typ
	arrayTypeMutex.Unlock()
}

func arrayOf(count int, elem reflect.Type) reflect.Type {
	key := arrayTypeKey{count, elem}

	arrayTypeMutex.RLock()
	defer arrayTypeMutex.RUnlock()
	if typ, ok := arrayTypeMap[key]; ok {
		return typ
	}
	panic(fmt.Errorf("unregistered array type [%d]%s", count, elem.String()))
}
