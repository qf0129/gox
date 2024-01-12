package crudx

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/arrayx"
	"github.com/qf0129/gox/errx"
	"github.com/qf0129/gox/gormx/daox"
	"github.com/qf0129/gox/respx"
	"github.com/qf0129/gox/structx"
)

func QueryManyHandler[T any](queryBodys ...daox.QueryBody) gin.HandlerFunc {
	return func(c *gin.Context) {
		var queryBody daox.QueryBody
		if len(queryBodys) > 0 {
			queryBody = queryBodys[0]
		} else {
			queryBody = daox.QueryBody{}
		}
		if queryBody.Params == nil {
			queryBody.Params = map[string]any{}
		}

		if err := c.ShouldBindQuery(queryBody); err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}

		queryBodyKeys := structx.GetJsonFields(queryBody)
		queryParams := c.Request.URL.Query()
		for k, v := range queryParams {
			if arrayx.HasStrItem(queryBodyKeys, k) {
				queryParams.Del(k)
			} else {
				if _, ok := queryBody.Params[k]; !ok {
					if len(v) == 1 {
						queryBody.Params[k] = v[0]
					} else {
						queryBody.Params[k] = v
					}
				}
			}
		}
		fmt.Printf(">>>> %v\n", queryBody)

		var data any
		var err error
		if queryBody.NoPaging {
			data, err = daox.QueryAll[T](&queryBody) // 不分页
		} else {
			data, err = daox.QueryPage[T](&queryBody) // 分页
		}
		if err != nil {
			respx.Err(c, errx.QueryDataFailed.AddErr(err))
			return
		}
		respx.OK(c, data)
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
