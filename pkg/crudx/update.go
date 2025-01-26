package crudx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/dbx"
	"github.com/qf0129/gox/pkg/errx"
	"github.com/qf0129/gox/pkg/ginx"
)

type UpdateOption struct {
	UpdateFields []string
}

type reqUpdateHandler[T any] struct {
	Id   string `binding:"required"`
	Data T
}

func UpdateHandler[T any](options ...UpdateOption) ginx.HandlerFunc {
	return func(c *gin.Context) (any, errx.Err) {
		var opt UpdateOption
		if len(options) > 0 {
			opt = options[0]
		} else {
			opt = UpdateOption{}
		}

		req := &reqUpdateHandler[T]{}
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, errx.InvalidJsonParams.AddErr(err)
		}

		_, err := dbx.QueryOneByPk[T](req.Id)
		if err != nil {
			return nil, errx.QueryDataFailed.AddErr(err)
		}

		if len(opt.UpdateFields) > 0 {
			if err := dbx.UpdateOneFiledsByPk(req.Id, req.Data, opt.UpdateFields); err != nil {
				return nil, errx.UpdateDataFailed.AddErr(err)
			}
		} else {
			if err := dbx.UpdateOneByPk(req.Id, req.Data); err != nil {
				return nil, errx.UpdateDataFailed.AddErr(err)
			}
		}
		return req.Id, nil
	}
}
