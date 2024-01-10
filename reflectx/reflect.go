package reflectx

import "reflect"

func GetVal(obj any, key string) any {
	return reflect.ValueOf(obj).FieldByName(key).Interface()
}
