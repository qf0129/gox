package crudx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/dbx"
	"github.com/qf0129/gox/pkg/errx"
	"github.com/qf0129/gox/pkg/ginx"
)

type reqDeleteHandler struct {
	Ids []string `binding:"required"`
}

func DeleteHandler[T any]() ginx.HandlerFunc {
	return func(c *gin.Context) (any, errx.Err) {
		var req reqDeleteHandler
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, errx.InvalidJsonParams.AddErr(err)
		}

		if len(req.Ids) == 0 {
			return nil, errx.InvalidParams.AddMsg("Ids不能为空")
		}

		deletedIds := []string{}
		for _, id := range req.Ids {
			exists, err := dbx.ExistByPk[T](id)
			if err != nil {
				return nil, errx.QueryDataFailed.AddErr(err).AddMsg("id=" + id)
			}
			if !exists {
				return nil, errx.DeleteDataFailed.AddMsg(id + "不存在")
			}
			if err := dbx.DeleteByPk[T](id); err != nil {
				return nil, errx.DeleteDataFailed.AddErr(err).AddMsg("id=" + id)
			} else {
				deletedIds = append(deletedIds, id)
			}
		}
		return deletedIds, nil
	}
}
