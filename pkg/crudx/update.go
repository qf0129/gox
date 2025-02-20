package crudx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/dbx"
	"github.com/qf0129/gox/pkg/errx"
	"github.com/qf0129/gox/pkg/ginx"
)

type UpdateHandlerOption struct {
	UpdateFields []string
	AfterHook    func(c *gin.Context, id any)
}

type reqUpdateHandler[T any] struct {
	Id   any `binding:"required"`
	Data T   `binding:"required"`
}

func UpdateHandler[T any](options ...UpdateHandlerOption) ginx.HandlerFunc {
	return func(c *gin.Context) (any, errx.Err) {
		var opt *UpdateHandlerOption
		if len(options) > 0 {
			opt = &options[0]
		}

		req := &reqUpdateHandler[T]{}
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, errx.InvalidJsonParams.AddErr(err)
		}

		target, err := dbx.QueryOneByPk[T](req.Id)
		if err != nil {
			return nil, errx.QueryDataFailed.AddErr(err).AddMsgf("id=%v", req.Id)
		}

		if opt != nil && len(opt.UpdateFields) > 0 {
			if err := dbx.UpdateTargetFileds(target, req.Data, opt.UpdateFields); err != nil {
				return nil, errx.UpdateDataFailed.AddErr(err)
			}
		} else {
			if err := dbx.UpdateTarget(target, req.Data); err != nil {
				return nil, errx.UpdateDataFailed.AddErr(err)
			}
		}

		if opt != nil && opt.AfterHook != nil {
			opt.AfterHook(c, req.Id)
		}
		return req.Id, nil
	}
}
