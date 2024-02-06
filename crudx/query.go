package crudx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/arrayx"
	"github.com/qf0129/gox/dbx"
	"github.com/qf0129/gox/errx"
	"github.com/qf0129/gox/respx"
	"github.com/qf0129/gox/structx"
)

func initOption(c *gin.Context, options ...dbx.QueryOption) (opt dbx.QueryOption, err error) {
	if len(options) > 0 {
		opt = options[0]
	} else {
		opt = dbx.QueryOption{}
	}
	if opt.Where == nil {
		opt.Where = map[string]any{}
	}
	body := &dbx.QueryBody{}
	if err = c.ShouldBindQuery(&body); err != nil {
		respx.Err(c, errx.InvalidParams.AddErr(err))
		return
	}

	if body.NoPaging {
		opt.NoPaging = body.NoPaging
	}
	if body.Preload != "" {
		opt.Preload = body.Preload
	}
	if body.PageNum > 0 {
		opt.PageNum = body.PageNum
	}
	if body.PageSize > 0 {
		opt.PageSize = body.PageSize
	}

	jsonFields := structx.GetJsonFields(body)
	queryParams := c.Request.URL.Query()
	for k, v := range queryParams {
		if arrayx.HasStrItem(jsonFields, k) {
			queryParams.Del(k)
		} else {
			if _, ok := opt.Where[k]; !ok {
				if len(v) == 1 {
					opt.Where[k] = v[0]
				} else {
					opt.Where[k] = v
				}
			}
		}
	}

	// 读取路径参数
	if opt.PathParamsMap != nil {
		for k, v := range opt.PathParamsMap {
			pathParamVal := c.Param(k)
			if pathParamVal != "" {
				opt.Where[v] = pathParamVal
			}
		}
	}
	return
}

func QueryManyHandler[T any](options ...dbx.QueryOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		opt, err := initOption(c, options...)
		if err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}

		var result any
		if opt.NoPaging {
			result, err = dbx.QueryAll[T](&opt) // 不分页
		} else {
			result, err = dbx.QueryPage[T](&opt) // 分页
		}
		if err != nil {
			respx.Err(c, errx.QueryDataFailed.AddErr(err))
			return
		}
		respx.OK(c, result)
	}
}

func QueryOneToManyHandler[P any, T any](relationField string, options ...dbx.QueryOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		opt, err := initOption(c, options...)
		if err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}

		_, err = dbx.QueryOneByPk[P](c.Param(dbx.Opt.ModelPrimaryKey))
		if err != nil {
			respx.Err(c, errx.QueryDataFailed.AddErr(err))
			return
		}
		opt.Where[relationField] = c.Param(dbx.Opt.ModelPrimaryKey)

		var result any
		if opt.NoPaging {
			result, err = dbx.QueryAll[T](&opt) // 不分页
		} else {
			result, err = dbx.QueryPage[T](&opt) // 分页
		}
		if err != nil {
			respx.Err(c, errx.QueryDataFailed.AddErr(err))
			return
		}
		respx.OK(c, result)
	}
}
func QueryAssociationHandler[P any, T any](field string, options ...dbx.QueryOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		opt, err := initOption(c, options...)
		if err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}

		model, err := dbx.QueryOneByPk[P](c.Param(dbx.Opt.ModelPrimaryKey))
		if err != nil {
			respx.Err(c, errx.QueryDataFailed.AddErr(err))
			return
		}

		var result any
		if opt.NoPaging {
			result, err = dbx.QueryAssociationAll[T](model, field, &opt)
		} else {
			result, err = dbx.QueryAssociationPage[T](model, field, &opt)
		}
		if err != nil {
			respx.Err(c, errx.QueryDataFailed.AddErr(err))
			return
		}
		respx.OK(c, result)
	}
}
