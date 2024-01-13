package crudx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/qf0129/gox/dbx"
	"github.com/qf0129/gox/errx"
	"github.com/qf0129/gox/respx"
)

func DeleteHandler[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		pk := c.Param(dbx.Opt.ModelPrimaryKey)
		if pk == "" {
			respx.Err(c, errx.InvalidParams.AddMsg("主键为空"))
			return
		}
		if er := dbx.DeleteOneByPk[T](pk); er != nil {
			if errMySQL, ok := er.(*mysql.MySQLError); ok {
				switch errMySQL.Number {
				case 1451:
					respx.Err(c, errx.DeleteDataFailed.AddMsg("无法删除有关联数据的项"))
					return
				}
			} else {
				respx.Err(c, errx.DeleteDataFailed.AddErr(er))
				return
			}
		}
		respx.OK(c, pk)
	}
}
