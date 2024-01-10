package daox

import (
	"reflect"
	"strings"

	"github.com/qf0129/gox/gormx"
)

// 查询分页数据
func QueryPage[T any](queryBodys ...*QueryBody) (result PageBody[T], err error) {
	var queryBody *QueryBody
	if len(queryBodys) > 0 {
		queryBody = queryBodys[0]
	} else {
		queryBody = &QueryBody{}
	}
	queryBody.ParseFilterToMap()
	if queryBody.PageNum < 1 {
		queryBody.PageNum = 1
	}
	if queryBody.PageSize < 1 {
		queryBody.PageSize = gormx.Opt.QueryPageSize
	}

	query := gormx.DB.Model(new(T))
	if len(queryBody.Fields) > 0 {
		query = query.Select(queryBody.Fields)
	}
	for _, fc := range ParseFilters(queryBody.FilterMap) {
		query = fc(query)
	}
	if err = query.Count(&result.Total).Error; err != nil {
		return
	}
	if queryBody.Preload != "" {
		for _, preload := range strings.Split(queryBody.Preload, ",") {
			query = FiltePreloadFunc(preload)(query)
		}
	}
	result.PageNum = queryBody.PageNum
	result.PageSize = queryBody.PageSize
	query = FiltePageFunc(result.PageNum, result.PageSize)(query)
	err = query.Find(&result.List).Error
	return
}

func QueryAll[T any](queryBodys ...*QueryBody) (result []T, err error) {
	var queryBody *QueryBody
	if len(queryBodys) > 0 {
		queryBody = queryBodys[0]
	} else {
		queryBody = &QueryBody{}
	}
	queryBody.ParseFilterToMap()

	query := gormx.DB.Model(new(T))
	if len(queryBody.Fields) > 0 {
		query = query.Select(queryBody.Fields)
	}
	for _, fc := range ParseFilters(queryBody.FilterMap) {
		query = fc(query)
	}
	if queryBody.Preload != "" {
		for _, preload := range strings.Split(queryBody.Preload, ",") {
			query = FiltePreloadFunc(preload)(query)
		}
	}
	err = query.Find(&result).Error
	return
}
func QueryAllToMap[T any](queryBodys ...*QueryBody) (result []map[string]any, err error) {
	var queryBody *QueryBody
	if len(queryBodys) > 0 {
		queryBody = queryBodys[0]
	} else {
		queryBody = &QueryBody{}
	}
	queryBody.ParseFilterToMap()

	query := gormx.DB.Model(new(T))
	for _, fc := range ParseFilters(queryBody.FilterMap) {
		query = fc(query)
	}

	if queryBody.Preload != "" {
		for _, preload := range strings.Split(queryBody.Preload, ",") {
			query = FiltePreloadFunc(preload)(query)
		}
	}
	err = query.Find(&result).Error
	return
}
func QueryAllCount[T any](queryBodys ...*QueryBody) (total int64, err error) {
	var queryBody *QueryBody
	if len(queryBodys) > 0 {
		queryBody = queryBodys[0]
	} else {
		queryBody = &QueryBody{}
	}
	queryBody.ParseFilterToMap()

	query := gormx.DB.Model(new(T))
	for _, fc := range ParseFilters(queryBody.FilterMap) {
		query = fc(query)
	}

	err = query.Count(&total).Error
	return
}

func ExistByPk[T any](pk any) (err error) {
	item := new(T)
	return gormx.DB.Model(new(T)).Where(map[string]any{gormx.Opt.ModelPrimaryKey: pk}).First(&item).Error
}

func QueryOneByPk[T any](pk any) (result T, err error) {
	err = gormx.DB.Model(new(T)).Where(map[string]any{gormx.Opt.ModelPrimaryKey: pk}).First(&result).Error
	return
}
func QueryTargetByPk[T any](pk any, tgt any) error {
	return gormx.DB.Model(new(T)).Where(map[string]any{gormx.Opt.ModelPrimaryKey: pk}).First(tgt).Error
}

func QueryOneByPkWithPreload[T any](pk any, preload string) (result T, err error) {
	query := gormx.DB.Model(new(T)).Where(map[string]any{gormx.Opt.ModelPrimaryKey: pk})
	for _, field := range strings.Split(preload, ",") {
		query = FiltePreloadFunc(field)(query)
	}
	err = query.Take(&result).Error
	return
}

func QueryOneByMap[T any](filters map[string]any) (result T, err error) {
	query := gormx.DB.Model(new(T))
	for _, fc := range ParseFilters(filters) {
		query = fc(query)
	}
	err = query.First(&result).Error
	return
}

func QueryOneByMapWithPreload[T any](filters map[string]any, preload string) (result T, err error) {
	query := gormx.DB.Model(new(T))
	for _, fc := range ParseFilters(filters) {
		query = fc(query)
	}

	for _, field := range strings.Split(preload, ",") {
		query = FiltePreloadFunc(field)(query)
	}

	err = query.First(&result).Error
	return
}

func CreateOne[T any](obj any) error {
	return gormx.DB.Model(new(T)).Create(obj).Error
}

func CreateOneWithParentId[T any](obj any, parentIdKey string, parentIdVal string) error {
	types := reflect.TypeOf(obj)
	vals := reflect.ValueOf(obj).Elem()
	for i := 0; i < types.NumField(); i++ {
		if types.Field(i).Name == parentIdKey {
			vals.Field(i).Set(reflect.ValueOf(parentIdVal))
		}
	}
	return gormx.DB.Model(new(T)).Create(&obj).Error
}

func UpdateByMap[T any](filters map[string]any, data any) error {
	return gormx.DB.Model(new(T)).Where(filters).Updates(data).Error
}

func UpdateOneByPk[T any](pk any, data any) error {
	return gormx.DB.Model(new(T)).Where(map[string]any{gormx.Opt.ModelPrimaryKey: pk}).Updates(data).Error
}

func DeleteByMap[T any](filters map[string]any) error {
	return gormx.DB.Where(filters).Delete(new(T)).Error
}

func DeleteOneByPk[T any](pk any) error {
	return gormx.DB.Where(map[string]any{gormx.Opt.ModelPrimaryKey: pk}).Delete(new(T)).Error
}

func HasField[T any](field string) bool {
	return gormx.DB.Model(new(T)).Select(field).Take(new(T)).Error == nil
}

func Save(obj any) error {
	return gormx.DB.Save(obj).Error
}
