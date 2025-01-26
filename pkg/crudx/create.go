package crudx

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/dbx"
	"github.com/qf0129/gox/pkg/errx"
	"github.com/qf0129/gox/pkg/ginx"
	"github.com/qf0129/gox/pkg/reflectx"
	"github.com/qf0129/gox/pkg/structx"
)

type CustomFieldFunc func(c *gin.Context, input map[string]any) any

// 查询选项
type CreateHandlerOption struct {
	CustomFields map[string]CustomFieldFunc
	PrimaryKey   string
	// PathParamsMap map[string]string
}

func CreateHandler[T any](options ...CreateHandlerOption) ginx.HandlerFunc {
	return func(c *gin.Context) (any, errx.Err) {
		var jsonData any
		if err := c.BindJSON(&jsonData); err != nil {
			return nil, errx.InvalidJsonParams.AddErr(err)
		}

		var opt *CreateHandlerOption
		if len(options) > 0 {
			opt = &options[0]
		}

		switch t := jsonData.(type) {
		case map[string]any: // 创建单个
			return createOne[T](c, jsonData, opt)
		case []any: // 创建多个
			ids := make([]string, 0)
			for _, item := range jsonData.([]any) {
				id, err := createOne[T](c, item, opt)
				if err != nil {
					return nil, err
				}
				ids = append(ids, id)
			}
			return ids, nil
		default:
			return nil, errx.InvalidJsonParams.AddMsg(fmt.Sprintf("不支持的类型: %v", reflect.TypeOf(t)))
		}
	}
}

func createOne[T any](c *gin.Context, itemData any, opt *CreateHandlerOption) (string, errx.Err) {
	if reflect.TypeOf(itemData).Kind() != reflect.Map {
		return "", errx.InvalidJsonParams.AddMsg(fmt.Sprintf("不支持的类型: %v", reflect.TypeOf(itemData)))
	}

	itemMap := itemData.(map[string]any)

	if opt != nil {
		// 自定义参数取值
		if opt.CustomFields != nil {
			for field, fieldFunc := range opt.CustomFields {
				itemMap[field] = fieldFunc(c, itemMap)
			}
		}
	}

	target, err := structx.MapToStruct[T](itemMap)
	if err != nil {
		return "", errx.PraseJsonError.AddErr(err)
	}

	if err := dbx.Create(target); err != nil {
		return "", errx.CreateDataFailed.AddErr(err)
	}

	pk := "Id"
	if opt != nil && opt.PrimaryKey != "" {
		pk = opt.PrimaryKey
	}

	idData := reflectx.GetVal(target, pk)
	if idData == nil {
		return "", errx.RequestFailed.AddMsgf("目标对象没有%s字段", pk)
	}
	return idData.(string), nil
}
