package crudx

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/arrayx"
	"github.com/qf0129/gox/pkg/convertx"
	"github.com/qf0129/gox/pkg/dbx"
	"github.com/qf0129/gox/pkg/errx"
	"github.com/qf0129/gox/pkg/reflectx"
	"github.com/qf0129/gox/pkg/serverx"
)

const (
	CrudMethodCreate = "create"
	CrudMethodDelete = "delete"
	CrudMethodRead   = "read"
	CrudMethodUpdate = "update"
	CrudMethodCount  = "count"
)

var CrudMethods = []string{
	CrudMethodCreate,
	CrudMethodDelete,
	CrudMethodRead,
	CrudMethodUpdate,
	CrudMethodCount,
}

type ReqCrudHandler struct {
	Model   string `binding:"required"`
	Method  string `binding:"required"`
	Select  *[]string
	Filter  *map[string]any
	Data    any
	Order   *string
	Limit   *int
	Offset  *int
	Preload *map[string][]any
}

type CrudModel struct {
	Model   any
	Methods string
}

var CrudPrimaryKey = "Id"

func CrudHandler(models map[string]CrudModel) serverx.HandlerFunc {
	return func(c *gin.Context) (any, errx.Err) {
		var req ReqCrudHandler
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, errx.InvalidParams.AddErr(err)
		}
		m, ok := models[req.Model]
		if !ok || m.Model == nil {
			return nil, errx.InvalidParams.AddMsgf("Model %s does not exist", req.Model)
		}
		if err := checkMethod(m.Methods, req.Method); err != nil {
			return nil, errx.InvalidParams.AddErr(err)
		}
		switch req.Method {
		case CrudMethodCreate:
			return handleCreate(m.Model, &req)
		case CrudMethodRead:
			return handleRead(m.Model, &req)
		case CrudMethodUpdate:
			return handleUpdate(m.Model, &req)
		case CrudMethodDelete:
			return handleDelete(m.Model, &req)
		case CrudMethodCount:
			return handleCount(m.Model, &req)
		default:
			return nil, errx.InvalidParams.AddMsg("Invalid crud method")
		}
	}
}

func handleCreate(model any, req *ReqCrudHandler) (any, errx.Err) {
	if req.Data == nil {
		return nil, errx.InvalidParams.AddMsg("Create data cannot be empty")
	}
	if reflectx.IsSlice(req.Data) {
		sliceData, ok := req.Data.([]any)
		if !ok {
			return nil, errx.InvalidParams.AddMsg("Create data must be a slice of map")
		}

		pks := []any{}
		for _, item := range sliceData {
			pk, err := handleCreateOne(model, item)
			if err != nil {
				return nil, err
			}
			pks = append(pks, pk)
		}
		return pks, nil
	}
	return handleCreateOne(model, req.Data)
}

func handleCreateOne(model any, data any) (any, errx.Err) {
	b, _ := json.Marshal(data)
	itemStruct := convertx.MakeAnyStuct(model)
	if err := json.Unmarshal(b, itemStruct); err != nil {
		return nil, errx.InvalidParams.AddErr(err)
	}
	if err := dbx.DB.Model(model).Create(itemStruct).Error; err != nil {
		return nil, errx.CreateDataFailed.AddErr(err)
	}
	return reflectx.StructGet(itemStruct, CrudPrimaryKey), nil
}

func handleDelete(model any, req *ReqCrudHandler) (any, errx.Err) {
	if req.Filter == nil {
		return nil, errx.InvalidParams.AddMsg("Where cannot be empty")
	}
	err := dbx.DB.Where(*req.Filter).Delete(model).Error
	if err != nil {
		return nil, errx.DeleteDataFailed.AddErr(err)
	}
	return true, nil
}

func handleUpdate(model any, req *ReqCrudHandler) (any, errx.Err) {
	if req.Data == nil {
		return nil, errx.InvalidParams.AddMsg("Update data cannot be empty")
	}
	query := dbx.DB.Model(model)
	if req.Filter != nil {
		query = query.Where(*req.Filter)
	}
	err := query.Updates(req.Data).Error
	if err != nil {
		return nil, errx.UpdateDataFailed.AddErr(err)
	}
	return true, nil
}

func handleRead(model any, req *ReqCrudHandler) (any, errx.Err) {
	query := dbx.DB.Model(model)
	if req.Select != nil && len(*req.Select) > 0 {
		query = query.Select(strings.Join(*req.Select, ","))
	}
	if req.Filter != nil {
		query = query.Where(*req.Filter)
	}
	if req.Order != nil {
		query = query.Order(*req.Order)
	}
	if req.Limit != nil {
		query = query.Limit(*req.Limit)
		if req.Offset != nil {
			query = query.Offset(*req.Offset)
		}
	}
	if req.Preload != nil {
		for k, v := range *req.Preload {
			query = query.Preload(k, v...)
		}
	}
	result := convertx.MakeAnySlice(model)
	if err := query.Find(result).Error; err != nil {
		return nil, errx.QueryDataFailed.AddErr(err)
	}
	return &result, nil
}

func handleCount(model any, req *ReqCrudHandler) (any, errx.Err) {
	query := dbx.DB.Model(model)
	if req.Filter != nil {
		query = query.Where(*req.Filter)
	}
	var result int64
	if err := query.Count(&result).Error; err != nil {
		return nil, errx.QueryDataFailed.AddErr(err)
	}
	return &result, nil
}

func checkMethod(methods string, method string) error {
	if !arrayx.Contains(CrudMethods, method) {
		return fmt.Errorf("invalid method %s", method)
	}
	if method == CrudMethodCount {
		return nil
	}
	if !arrayx.Contains([]rune(methods), rune(method[0])) {
		return fmt.Errorf("method %s not support", method)
	}
	return nil
}
