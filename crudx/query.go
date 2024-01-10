package crudx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/errx"
	"github.com/qf0129/gox/gormx/daox"
	"github.com/qf0129/gox/respx"
)

func QueryManyHandler[T any](queryBodys ...daox.QueryBody) gin.HandlerFunc {
	return func(c *gin.Context) {
		var queryBody daox.QueryBody
		if len(queryBodys) > 0 {
			queryBody = queryBodys[0]
		} else {
			queryBody = daox.QueryBody{}
		}

		err := c.ShouldBindQuery(&queryBody)
		if err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}
		if queryBody.NoPaging {
			// 不分页
			data, er := daox.QueryAll[T](&queryBody)
			if er != nil {
				respx.Err(c, errx.QueryDataFailed.AddErr(er))
				return
			}
			respx.OK(c, data)
		} else {
			// 分页
			data, er := daox.QueryPage[T](&queryBody)
			if er != nil {
				respx.Err(c, errx.QueryDataFailed.AddErr(er))
				return
			}
			respx.OK(c, data)
		}
	}
}

func QueryManyChildHandler[T any](parentId string, queryBodys ...daox.QueryBody) gin.HandlerFunc {
	return func(c *gin.Context) {
		var queryBody daox.QueryBody
		if len(queryBodys) > 0 {
			queryBody = queryBodys[0]
		} else {
			queryBody = daox.QueryBody{}
		}

		if queryBody.FilterMap == nil {
			queryBody.FilterMap = map[string]any{}
		}
		queryBody.FilterMap[parentId] = c.Param("id")
		queryBody.FilterMap["ctime:ob"] = "desc"

		err := c.ShouldBindQuery(&queryBody)
		if err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}
		if queryBody.NoPaging {
			// 不分页
			data, er := daox.QueryAll[T](&queryBody)
			if er != nil {
				respx.Err(c, errx.QueryDataFailed.AddErr(er))
				return
			}
			respx.OK(c, data)
		} else {
			// 分页
			data, er := daox.QueryPage[T](&queryBody)
			if er != nil {
				respx.Err(c, errx.QueryDataFailed.AddErr(er))
				return
			}
			respx.OK(c, data)
		}
	}
}
