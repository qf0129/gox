package crudx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/dbx"
	"github.com/qf0129/gox/errx"
	"github.com/qf0129/gox/jsonx"
	"github.com/qf0129/gox/respx"
)

type UpdateOption struct {
	IgnoreFields []string
}

func UpdateHandler[T any](options ...UpdateOption) gin.HandlerFunc {

	return func(c *gin.Context) {

		var opt UpdateOption
		if len(options) > 0 {
			opt = options[0]
		} else {
			opt = UpdateOption{}
		}

		if opt.IgnoreFields == nil {
			opt.IgnoreFields = []string{}
		}

		opt.IgnoreFields = append(opt.IgnoreFields, dbx.Opt.ModelPrimaryKey)

		// 请求json转map
		postMap := map[string]any{}
		if err := c.ShouldBindJSON(&postMap); err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}

		// 若包含主键则删除
		for k := range postMap {
			for _, f := range opt.IgnoreFields {
				if k == f {
					delete(postMap, k)
				}
			}
		}

		// 获取路径中主键
		pk := c.Param(dbx.Opt.ModelPrimaryKey)
		if pk == "" {
			respx.Err(c, errx.InvalidPathParams)
			return
		}
		// 查询目标对象
		obj, err := dbx.QueryOneByPk[T](pk)
		if err != nil {
			respx.Err(c, errx.TargetNotExists)
			return
		}

		// map转json
		jsonByte, err := jsonx.Marshal(postMap)
		if err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}

		// json赋值给对象
		err = jsonx.Unmarshal(jsonByte, &obj)
		if err != nil {
			respx.Err(c, errx.InvalidParams.AddErr(err))
			return
		}

		// 更新数据
		if err = dbx.DB.Save(&obj).Error; err != nil {
			respx.Err(c, errx.UpdateDataFailed.AddErr(err))
			return
		}
		respx.OK(c, true)
	}
}
