package dbx

import (
	"reflect"
	"strings"

	"github.com/qf0129/gox/constx"
	"gorm.io/gorm"
)

type PageBody[T any] struct {
	List     []T   `json:"list"`
	PageNum  int   `json:"page_num"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type QueryBody struct {
	PageNum  int    `json:"page_num" form:"page_num"`   // 页数, 默认1
	PageSize int    `json:"page_size" form:"page_size"` // 每页数量, 默认10
	NoPaging bool   `json:"no_paging" form:"no_paging"` // 关闭分页, 默认false
	Preload  string `json:"preload" form:"preload"`     // 预加载关联表名, 若多个以英文逗号分隔
}

type QueryOption struct {
	// QueryBody

	PageNum  int    // 页数, 默认1
	PageSize int    // 每页数量, 默认10
	NoPaging bool   // 关闭分页, 默认false
	Preload  string // 预加载关联表名, 若多个以英文逗号分隔

	Select []string       // 查询字段
	Where  map[string]any // 条件map
	// 从请求url的路径参数中获取查询条件，
	// {路径参数key: 数据库列名,}
	// 如：GET /api/user/:id/roles，则应为{"id": "user_id"}
	PathParamsMap map[string]string
}

func getQueryOption(options ...*QueryOption) (opt *QueryOption) {
	if len(options) > 0 {
		opt = options[0]
	} else {
		opt = &QueryOption{}
	}
	if opt.PageNum < 1 {
		opt.PageNum = 1
	}
	if opt.PageSize < 1 {
		opt.PageSize = constx.DefaultQueryPageSize
	}
	return
}

func setQuerySelect(query *gorm.DB, selectFields []string) {
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}
}

func setQueryWhere(query *gorm.DB, paramsMap map[string]any) {
	for _, filteFunc := range paramsMapToFilters(paramsMap) {
		query = filteFunc(query)
	}
}

func setQuerySelectAndWhere(query *gorm.DB, selectFields []string, paramsMap map[string]any) {
	setQuerySelect(query, selectFields)
	setQueryWhere(query, paramsMap)
}

func setQueryPreload(query *gorm.DB, preloadStr string) {
	if preloadStr != "" {
		for _, preload := range strings.Split(preloadStr, ",") {
			query = FiltePreloadFunc(preload)(query)
		}
	}
}

// 查询分页
func QueryPage[T any](options ...*QueryOption) (result PageBody[T], err error) {
	opt := getQueryOption(options...)
	query := DB.Model(new(T))
	setQuerySelectAndWhere(query, opt.Select, opt.Where)
	if err = query.Count(&result.Total).Error; err != nil {
		return
	}
	setQueryPreload(query, opt.Preload)
	result.PageNum = opt.PageNum
	result.PageSize = opt.PageSize
	query = FiltePageFunc(result.PageNum, result.PageSize)(query)
	err = query.Find(&result.List).Error
	return
}

// 查询不分页
func QueryAll[T any](options ...*QueryOption) (result []T, err error) {
	opt := getQueryOption(options...)
	query := DB.Model(new(T))
	setQuerySelectAndWhere(query, opt.Select, opt.Where)
	setQueryPreload(query, opt.Preload)
	err = query.Find(&result).Error
	return
}

// 子查询分页
func QueryAssociationPage[T any](model any, field string, options ...*QueryOption) (result PageBody[T], err error) {
	opt := getQueryOption(options...)
	query := DB.Model(model)
	setQuerySelectAndWhere(query, opt.Select, opt.Where)
	result.Total = query.Association(field).Count()
	query = DB.Model(model) // Association查询需要新建对象
	setQuerySelectAndWhere(query, opt.Select, opt.Where)
	result.PageNum = opt.PageNum
	result.PageSize = opt.PageSize
	query = FiltePageFunc(result.PageNum, result.PageSize)(query)
	err = query.Association(field).Find(&result.List)
	return
}

// 子查询不分页
func QueryAssociationAll[T any](model any, field string, options ...*QueryOption) (result PageBody[T], err error) {
	opt := getQueryOption(options...)
	query := DB.Model(model)
	setQuerySelectAndWhere(query, opt.Select, opt.Where)
	err = query.Association(field).Find(&result.List)
	return
}

// 查询所有到指定map
func QueryAllToMap[T any](options ...*QueryOption) (result []map[string]any, err error) {
	opt := getQueryOption(options...)
	query := DB.Model(new(T))
	setQuerySelectAndWhere(query, opt.Select, opt.Where)
	setQueryPreload(query, opt.Preload)
	err = query.Find(&result).Error
	return
}

// 查询数量
func QueryAllCount[T any](options ...*QueryOption) (total int64, err error) {
	opt := getQueryOption(options...)
	query := DB.Model(new(T))
	setQuerySelectAndWhere(query, opt.Select, opt.Where)
	err = query.Count(&total).Error
	return
}

func ExistByPk[T any](pk any) (err error) {
	item := new(T)
	return DB.Model(new(T)).Where(map[string]any{Opt.ModelPrimaryKey: pk}).First(&item).Error
}

func QueryOneByPk[T any](pk any) (result T, err error) {
	err = DB.Model(new(T)).Where(map[string]any{Opt.ModelPrimaryKey: pk}).First(&result).Error
	return
}
func QueryTargetByPk[T any](pk any, tgt any) error {
	return DB.Model(new(T)).Where(map[string]any{Opt.ModelPrimaryKey: pk}).First(tgt).Error
}

func QueryOneByPkWithPreload[T any](pk any, preload string) (result T, err error) {
	query := DB.Model(new(T)).Where(map[string]any{Opt.ModelPrimaryKey: pk})
	setQueryPreload(query, preload)
	err = query.Take(&result).Error
	return
}

func QueryOneByMap[T any](paramsMap map[string]any) (result T, err error) {
	query := DB.Model(new(T))
	setQueryWhere(query, paramsMap)
	err = query.First(&result).Error
	return
}

func QueryOneByMapWithPreload[T any](paramsMap map[string]any, preload string) (result T, err error) {
	query := DB.Model(new(T))
	setQueryWhere(query, paramsMap)
	setQueryPreload(query, preload)
	err = query.First(&result).Error
	return
}

func Create[T any](items any) error {
	return DB.Model(new(T)).Create(items).Error
}

func CreateOneWithParentId[T any](obj any, parentIdKey string, parentIdVal string) error {
	types := reflect.TypeOf(obj)
	vals := reflect.ValueOf(obj).Elem()
	for i := 0; i < types.NumField(); i++ {
		if types.Field(i).Name == parentIdKey {
			vals.Field(i).Set(reflect.ValueOf(parentIdVal))
		}
	}
	return DB.Model(new(T)).Create(&obj).Error
}

func UpdateByMap[T any](filters map[string]any, data any) error {
	return DB.Model(new(T)).Where(filters).Updates(data).Error
}

func UpdateOneByPk[T any](pk any, data any) error {
	return DB.Model(new(T)).Where(map[string]any{Opt.ModelPrimaryKey: pk}).Updates(data).Error
}

func DeleteByMap[T any](filters map[string]any) error {
	return DB.Where(filters).Delete(new(T)).Error
}

func DeleteOneByPk[T any](pk any) error {
	return DB.Where(map[string]any{Opt.ModelPrimaryKey: pk}).Delete(new(T)).Error
}

func HasField[T any](field string) bool {
	return DB.Model(new(T)).Select(field).Take(new(T)).Error == nil
}

func Save(obj any) error {
	return DB.Save(obj).Error
}
