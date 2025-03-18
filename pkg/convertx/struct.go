package convertx

import (
	"github.com/qf0129/gox/pkg/jsonx"
)

func StructToMap(data any) (map[string]any, error) {
	jsonByte, err := jsonx.Marshal(data)
	if err != nil {
		return nil, err
	}
	var target map[string]any
	err = jsonx.Unmarshal(jsonByte, &target)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func MapToStruct[T any](data map[string]any) (*T, error) {
	jsonByte, err := jsonx.Marshal(data)
	if err != nil {
		return nil, err
	}

	var target *T
	err = jsonx.Unmarshal(jsonByte, &target)
	if err != nil {
		return nil, err
	}

	return target, nil
}
