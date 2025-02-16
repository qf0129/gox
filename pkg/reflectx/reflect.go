package reflectx

import "reflect"

func StructHas(s any, field string) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}
	return v.FieldByName(field).IsValid()
}

func StructGet(s any, field string) any {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	if v.FieldByName(field).IsValid() {
		return v.FieldByName(field).Interface()
	}
	return nil
}

func StructSet(s any, field string, val any) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}
	f := v.FieldByName(field)
	if !f.IsValid() || !f.CanSet() {
		return false
	}

	fType := f.Type()
	vVal := reflect.ValueOf(val)
	if vVal.Type().ConvertibleTo(fType) {
		f.Set(vVal.Convert(fType))
		return true
	}
	return false
}
