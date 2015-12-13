// +build go1.5

package restruct

import "reflect"

// RegisterArrayType registers an array type for use with struct tags. This
// function is only necessary on Go 1.4 and below, which do not provide
// reflect.ArrayOf. This function is goroutine safe and idempotent (i.e.,
// calling it multiple times with the same value is perfectly safe.) If you
// require Go 1.5 or above, you can safely use any array type in your struct
// tags without using this function.
func RegisterArrayType(array interface{}) {
}

func arrayOf(count int, elem reflect.Type) reflect.Type {
	return reflect.ArrayOf(count, elem)
}
