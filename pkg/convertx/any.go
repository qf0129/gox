package convertx

import (
	"fmt"
	"strconv"
	"strings"
)

func AnyToInt64(input any) (int64, error) {
	switch v := input.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("不支持将 %T 转换为int64", input)
	}
}

func AnyToFloat64(input any) (float64, error) {
	switch v := input.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("不支持将 %T 转换为float64", input)
	}
}

func AnyToBool(input any) (bool, error) {
	switch v := input.(type) {
	case bool:
		return v, nil
	case int:
		if v > 0 {
			return true, nil
		}
		return false, nil
	case string:
		if strings.ToLower(v) == "true" {
			return true, nil
		} else if strings.ToLower(v) == "false" {
			return false, nil
		}
	}
	return false, fmt.Errorf("不支持将 %v 转换为bool", input)
}

func AnyToSlice(input any) ([]any, error) {
	switch v := input.(type) {
	case []any:
		return v, nil
	default:
		return nil, fmt.Errorf("不支持将 %T 转换为切片类型", input)
	}
}

func AnyToString(input any) string {
	switch v := input.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func AnyToBytes(input any) []byte {
	return []byte(AnyToString(input))
}

func AnyMapToStrMap(anyMap map[string]any) map[string]string {
	dstMap := make(map[string]string, len(anyMap))
	for k, v := range anyMap {
		dstMap[k] = AnyToString(v)
	}
	return dstMap
}

func AnySliceToStrSlice(anySlice []any) ([]string, error) {
	strSlice := make([]string, 0, len(anySlice))
	for _, v := range anySlice {
		strSlice = append(strSlice, AnyToString(v))
	}
	return strSlice, nil
}
