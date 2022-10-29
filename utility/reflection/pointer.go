package reflection

import "reflect"

func IsPointer(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Ptr
}
