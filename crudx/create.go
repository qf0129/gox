package crudx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/dbx"
	"github.com/qf0129/gox/errx"
	"github.com/qf0129/gox/jsonx"
	"github.com/qf0129/gox/respx"
)

// 查询选项
type CreateManyOption struct {
	ExtraParams   map[string]string
	PathParamsMap map[string]string
}

func CreateManyHandler[T any](options ...CreateManyOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求json转map数组
		var mapList []map[string]any
		if err := c.ShouldBindJSON(&mapList); err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}

		var opt *CreateManyOption
		if len(options) > 0 {
			opt = &options[0]
		}

		if opt != nil {
			// 添加额外参数
			if opt.ExtraParams != nil {
				for _, item := range mapList {
					for k, v := range opt.ExtraParams {
						item[k] = v
					}
				}
			}
			if opt.PathParamsMap != nil {
				// 读取路径参数
				pathParams := map[string]string{}
				for k, v := range opt.PathParamsMap {
					pathParamVal := c.Param(k)
					if pathParamVal != "" {
						pathParams[v] = pathParamVal
					}
				}
				// 添加路径参数
				for _, item := range mapList {
					for k, v := range pathParams {
						item[k] = v
					}
				}
			}
		}

		// map转json
		jsonByte, err := jsonx.Marshal(mapList)
		if err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}

		// json转结构体数组
		var items []T
		err = jsonx.Unmarshal(jsonByte, &items)
		if err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}

		// 存储结构体
		if err := dbx.Create[T](items); err != nil {
			respx.Err(c, errx.CreateDataFailed.AddErr(err))
			return
		}
		respx.OK(c, true)
	}
}
