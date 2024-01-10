package crudx

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/errx"
	"github.com/qf0129/gox/gormx/daox"
	"github.com/qf0129/gox/respx"
	"github.com/qf0129/gox/strx"
)

func CreateModelHandler[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var model T
		if er := c.ShouldBindJSON(&model); er != nil {
			respx.Err(c, errx.InvalidParams.AddErr(er))
			return
		}

		er := daox.CreateOne[T](&model)
		if er != nil {
			respx.Err(c, errx.CreateDataFailed.AddErr(er))
			return
		}
		respx.OK(c, model)
	}
}

func CreateModelChildHandler[T any](parentId string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var model T
		if er := c.ShouldBindJSON(&model); er != nil {
			respx.Err(c, errx.InvalidParams.AddErr(er))
			return
		}

		parentIdField := reflect.ValueOf(&model).Elem().FieldByName(strx.SnakeToCamelCase(parentId))
		if parentIdField.CanSet() {
			parentIdField.SetString(c.Param("id"))
		}

		er := daox.CreateOne[T](&model)
		if er != nil {
			respx.Err(c, errx.CreateDataFailed.AddErr(er))
			return
		}
		respx.OK(c, model)
	}
}
