package crudx

import (
	"encoding/json"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/errx"
	"github.com/qf0129/gox/gormx/daox"
	"github.com/qf0129/gox/respx"
)

func UpdateHandler[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := make(map[string]any)
		if err := c.ShouldBindJSON(&body); err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}
		pk := c.Param("id")
		if pk == "" {
			respx.Err(c, errx.InvalidParams.AddMsg("主键为空"))
			return
		}
		if _, er := daox.QueryOneByPk[T](pk); er != nil {
			respx.Err(c, errx.DataNotExists)
			return
		}

		// gorm中updates结构体不支持更新空值，使用map不支持json类型
		// 因此遍历map，将子结构的map或slice转成json字符串
		for k, v := range body {
			valKind := reflect.ValueOf(v).Kind()
			if valKind == reflect.Map || valKind == reflect.Slice {
				bytes, er := json.Marshal(v)
				if er != nil {
					respx.Err(c, errx.InvalidParams.AddErr(er))
					return
				}
				body[k] = string(bytes)
			}
		}

		er := daox.UpdateOneByPk[T](pk, &body)
		if er != nil {
			respx.Err(c, errx.UpdateDataFailed.AddErr(er))
			return
		}

		newModel, er := daox.QueryOneByPk[T](pk)
		if er != nil {
			respx.Err(c, errx.QueryDataFailed.AddErr(er))
			return
		}
		respx.OK(c, newModel)
	}
}
