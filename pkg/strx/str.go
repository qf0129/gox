package strx

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GetModelNameLower(obj any) string {
	reflectType := reflect.TypeOf(obj)
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	return CamelToSnakeCase(reflectType.Name())
}

func GetIdKey(obj any) string {
	reflectType := reflect.TypeOf(obj)
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	return CamelToSnakeCase(reflectType.Name()) + "_id"
}

// AbcDef > abc_def
func CamelToSnakeCase(text string) string {
	temp := regexp.MustCompile(`([A-Z])`).ReplaceAllString(text, "_$1")
	temp = cases.Lower(language.Und).String(temp)
	temp = strings.TrimLeft(temp, "_")
	return temp
}

// abc_def > AbcDef
func SnakeToCamelCase(text string) string {
	temp := strings.ReplaceAll(text, "_", " ")
	temp = cases.Title(language.Und).String(temp)
	return strings.ReplaceAll(temp, " ", "")
}

// 是否为整数或浮点数
func IsNumeric(s string) bool {
	if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		return true // 字符串全是整数
	}
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return true // 字符串全是浮点数
	}
	return false
}
