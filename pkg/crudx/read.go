package crudx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/dbx"
	"github.com/qf0129/gox/pkg/errx"
	"github.com/qf0129/gox/pkg/ginx"
)

func ReadHandler[T any]() ginx.HandlerFunc {
	return func(c *gin.Context) (any, errx.Err) {
		var req dbx.QueryOption
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, errx.InvalidParams.AddErr(err)
		}
		result, err := dbx.QueryPage[T](&req)
		if err != nil {
			return nil, errx.QueryDataFailed.AddErr(err)
		}
		return result, nil
	}
}
