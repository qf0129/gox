package structx

import (
	"log/slog"
	"reflect"
	"strings"
)

func GetFields(obj any) []string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		slog.Warn("structs: InvalidStruct, " + t.Kind().String())
		return nil
	}

	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

func GetJsonFields(obj any) []string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	keys := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagText := field.Tag.Get("json")
		tags := strings.Split(tagText, ",")
		if len(tags) > 0 && tags[0] != "-" {
			if tags[0] == "" {
				keys = append(keys, field.Name)
			} else {
				keys = append(keys, tags[0])
			}
		}
	}
	return keys
}
