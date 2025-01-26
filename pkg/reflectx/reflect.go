package reflectx

import "reflect"

func GetVal(obj any, key string) any {
	v := reflect.ValueOf(obj)
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		v = reflect.ValueOf(obj).Elem()
	}
	if v.FieldByName(key).IsValid() {
		return v.FieldByName(key).Interface()
	}
	return nil
}
