package crudx

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/arrayx"
	"github.com/qf0129/gox/pkg/convertx"
	"github.com/qf0129/gox/pkg/dbx"
	"github.com/qf0129/gox/pkg/errx"
	"github.com/qf0129/gox/pkg/ptrx"
	"github.com/qf0129/gox/pkg/serverx"
	"github.com/qf0129/gox/reflectx"
	"gorm.io/gorm"
)

type Request struct {
	Model    string            `binding:"required"` // 模型名称
	Action   string            `binding:"required"` // 操作名称
	Fields   *[]string         // 字段列表
	Filters  *map[string]any   // 查询条件
	Preloads *map[string][]any // 关联查询
	Pk       any               // 主键值
	Data     any               // 更新或创建数据
	Order    *string           // 排序字段
	Page     *int              // 页码
	PageSize *int              // 每页数量
}

type HandlerContext struct {
	GinContext  *gin.Context   // gin上下文
	CrudOption  *HandlerOption // crud选项
	ModelOption *ModelOption   // 模型选项
	Request     *Request       // 请求参数
}
type ModelOption struct {
	Model           any                                    // 模型实例
	IgnoreActions   []string                               // 忽略的操作列表
	AfterCreateHook func(c *HandlerContext, row any) error // 创建后钩子函数
	// AfterUpdateHook func(c *HandlerContext, row any) error
	// BeforeDeleteHook func(c *HandlerContext, row any) error
}

type HandlerOption struct {
	Models               map[string]ModelOption // 模型选项映射
	PrimaryKey           string                 // 主键字段名
	QueryAllMax          int                    // 查询所有最大数量
	QueryPageSizeDefault int                    // 查询分页默认数量
}

type ActionHandler func(c *HandlerContext) (any, errx.Err)

var ActionHandlers = map[string]ActionHandler{
	"QueryFirst": handleQueryFirst,
	"QueryMany":  handleQueryMany,
	"QueryPage":  handleQueryPage,
	"Create":     handleCreate,
	"Update":     handleUpdate,
	"Delete":     handleDelete,
	"UpdateByPk": handleUpdateByPk,
	"DeleteByPk": handleDeleteByPk,
	"Count":      handleCount,
}

func CrudHandler(opt *HandlerOption) serverx.HandlerFunc {
	if opt == nil || opt.Models == nil {
		panic("CrudOption.Models is required")
	}
	if opt.PrimaryKey == "" {
		opt.PrimaryKey = "Id"
	}
	return func(ctx *gin.Context) (any, errx.Err) {
		var req Request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			return nil, errx.InvalidParams.AddErr(err)
		}
		mod, ok := opt.Models[req.Model]
		if !ok || mod.Model == nil {
			return nil, errx.InvalidParams.AddMsgf("Unknown model %s", req.Model)
		}
		if arrayx.Contains(mod.IgnoreActions, req.Action) {
			return nil, errx.InvalidParams.AddMsgf("No permission to %s", req.Action)
		}
		handler, ok := ActionHandlers[req.Action]
		if !ok {
			return nil, errx.InvalidParams.AddMsgf("Unknown action %s", req.Action)
		}
		return handler(&HandlerContext{GinContext: ctx, CrudOption: opt, ModelOption: &mod, Request: &req})
	}
}

func handleQueryFirst(c *HandlerContext) (any, errx.Err) {
	query := makeQuery(c)
	row := convertx.MakeAnyStuct(c.ModelOption.Model)
	if err := query.First(row).Error; err != nil {
		return nil, errx.QueryFailed.AddErr(err)
	}
	return &row, nil
}
func handleQueryMany(c *HandlerContext) (any, errx.Err) {
	query := makeQuery(c)
	if c.CrudOption.QueryAllMax > 0 {
		query = query.Limit(c.CrudOption.QueryAllMax)
	}
	rows := convertx.MakeAnySlice(c.ModelOption.Model)
	if err := query.Find(rows).Error; err != nil {
		return nil, errx.QueryFailed.AddErr(err)
	}
	return &rows, nil
}

type PageBody struct {
	List     any
	Page     int
	PageSize int
	Total    int64
}

func handleQueryPage(c *HandlerContext) (any, errx.Err) {
	query := makeQuery(c)
	total := int64(0)
	if err := query.Count(&total).Error; err != nil {
		return nil, errx.QueryFailed.AddErr(err)
	}
	if c.Request.Page == nil {
		c.Request.Page = ptrx.Ptr(1)
	}
	if c.Request.PageSize == nil {
		if c.CrudOption.QueryPageSizeDefault > 0 {
			c.Request.PageSize = ptrx.Ptr(c.CrudOption.QueryPageSizeDefault)
		} else {
			c.Request.PageSize = ptrx.Ptr(10)
		}
	}
	query = query.Limit(*c.Request.PageSize).Offset((*c.Request.Page - 1) * *c.Request.PageSize)
	rows := convertx.MakeAnySlice(c.ModelOption.Model)
	if err := query.Find(rows).Error; err != nil {
		return nil, errx.QueryFailed.AddErr(err)
	}
	result := &PageBody{
		Total:    total,
		Page:     *c.Request.Page,
		PageSize: *c.Request.PageSize,
		List:     rows,
	}
	return result, nil
}

func makeQuery(c *HandlerContext) *gorm.DB {
	query := dbx.DB.Model(c.ModelOption.Model)
	if c.Request.Fields != nil && len(*c.Request.Fields) > 0 {
		query = query.Select(*c.Request.Fields)
	}
	if c.Request.Filters != nil {
		query = query.Where(*c.Request.Filters)
	}
	if c.Request.Order != nil {
		query = query.Order(*c.Request.Order)
	}
	if c.Request.Preloads != nil {
		for k, args := range *c.Request.Preloads {
			query = query.Preload(k, args...)
		}
	}
	return query
}

func handleCreate(c *HandlerContext) (any, errx.Err) {
	if c.Request.Data == nil {
		return nil, errx.InvalidParams.AddMsg("Create action must have data")
	}
	if reflectx.IsSlice(c.Request.Data) {
		sliceData, ok := c.Request.Data.([]any)
		if !ok {
			return nil, errx.InvalidParams.AddMsg("Create action data must be map or []map")
		}
		pks := []any{}
		for _, item := range sliceData {
			pk, err := handleCreateOne(&HandlerContext{CrudOption: c.CrudOption, ModelOption: c.ModelOption, Request: &Request{Data: item}})
			if err != nil {
				return nil, err
			}
			pks = append(pks, pk)
		}
		return pks, nil
	}
	return handleCreateOne(c)
}

func handleCreateOne(c *HandlerContext) (any, errx.Err) {
	b, _ := json.Marshal(c.Request.Data)
	itemStruct := convertx.MakeAnyStuct(c.ModelOption.Model)
	if err := json.Unmarshal(b, itemStruct); err != nil {
		return nil, errx.InvalidParams.AddErr(err)
	}
	if err := dbx.DB.Model(c.ModelOption.Model).Create(itemStruct).Error; err != nil {
		return nil, errx.CreateFailed.AddErr(err)
	}
	if c.ModelOption.AfterCreateHook != nil {
		if err := c.ModelOption.AfterCreateHook(c, itemStruct); err != nil {
			return nil, errx.CreateFailed.AddErr(err)
		}
	}
	return itemStruct, nil
}

func handleDelete(c *HandlerContext) (any, errx.Err) {
	if c.Request.Filters == nil {
		return nil, errx.InvalidParams.AddMsg("Delete action must have filter")
	}
	err := dbx.DB.Where(*c.Request.Filters).Delete(c.ModelOption.Model).Error
	if err != nil {
		return nil, errx.DeleteFailed.AddErr(err)
	}
	return true, nil
}

func handleDeleteByPk(c *HandlerContext) (any, errx.Err) {
	if c.Request.Pk == nil {
		return nil, errx.InvalidParams.AddMsg("DeleteByPk action must have pk")
	}
	err := dbx.DB.Where(map[string]any{c.CrudOption.PrimaryKey: c.Request.Pk}).Delete(c.ModelOption.Model).Error
	if err != nil {
		return nil, errx.DeleteFailed.AddErr(err)
	}
	return true, nil
}

func handleUpdate(c *HandlerContext) (any, errx.Err) {
	if c.Request.Data == nil {
		return nil, errx.InvalidParams.AddMsg("Update action must have data")
	}
	if c.Request.Filters == nil {
		return nil, errx.InvalidParams.AddMsg("Update action must have filter")
	}
	query := dbx.DB.Model(c.ModelOption.Model).Where(*c.Request.Filters)
	if c.Request.Fields != nil && len(*c.Request.Fields) > 0 {
		query = query.Select(*c.Request.Fields)
	}
	err := query.Updates(c.Request.Data).Error
	if err != nil {
		return nil, errx.UpdateFailed.AddErr(err)
	}
	return true, nil
}
func handleUpdateByPk(c *HandlerContext) (any, errx.Err) {
	if c.Request.Pk == nil {
		return nil, errx.InvalidParams.AddMsg("UpdateByPk action must have pk")
	}
	if c.Request.Data == nil {
		return nil, errx.InvalidParams.AddMsg("Update action must have data")
	}
	query := dbx.DB.Model(c.ModelOption.Model).Where(map[string]any{c.CrudOption.PrimaryKey: c.Request.Pk})
	if c.Request.Fields != nil && len(*c.Request.Fields) > 0 {
		query = query.Select(*c.Request.Fields)
	}
	err := query.Updates(c.Request.Data).Error
	if err != nil {
		return nil, errx.UpdateFailed.AddErr(err)
	}
	return c.Request.Pk, nil
}

func handleCount(c *HandlerContext) (any, errx.Err) {
	query := dbx.DB.Model(c.ModelOption.Model)
	if c.Request.Filters != nil {
		query = query.Where(*c.Request.Filters)
	}
	var result int64
	if err := query.Count(&result).Error; err != nil {
		return nil, errx.QueryFailed.AddErr(err)
	}
	return &result, nil
}
