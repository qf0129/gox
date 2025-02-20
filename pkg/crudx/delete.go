package crudx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/dbx"
	"github.com/qf0129/gox/pkg/errx"
	"github.com/qf0129/gox/pkg/ginx"
)

type reqDeleteHandler struct {
	Ids []any `binding:"required"`
}

type DeleteHandlerOption struct {
	AfterHook func(c *gin.Context, id any)
}

func DeleteHandler[T any](options ...DeleteHandlerOption) ginx.HandlerFunc {
	return func(c *gin.Context) (any, errx.Err) {
		var req reqDeleteHandler
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, errx.InvalidJsonParams.AddErr(err)
		}

		var opt *DeleteHandlerOption
		if len(options) > 0 {
			opt = &options[0]
		}

		if len(req.Ids) == 0 {
			return nil, errx.InvalidParams.AddMsg("Ids不能为空")
		}

		deletedIds := []any{}
		for _, id := range req.Ids {
			if exists, _ := dbx.ExistByPk[T](id); !exists {
				return nil, errx.TargetNotExists.AddMsgf("id=%v", id)
			}
			if err := dbx.DeleteByPk[T](id); err != nil {
				return nil, errx.DeleteDataFailed.AddErr(err).AddMsgf("id=%v", id)
			}
			deletedIds = append(deletedIds, id)
			if opt != nil && opt.AfterHook != nil {
				opt.AfterHook(c, id)
			}
		}
		return deletedIds, nil
	}
}
