package dbx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/qf0129/gox/pkg/arrayx"
	"gorm.io/gorm"
)

type PageBody[T any] struct {
	List     []T
	Page     int
	PageSize int
	Total    int64
}

type QueryBody struct {
	Page     int    `json:"page_num" form:"page_num"`   // 页数, 默认1
	PageSize int    `json:"page_size" form:"page_size"` // 每页数量, 默认10
	NoPaging bool   `json:"no_paging" form:"no_paging"` // 关闭分页, 默认false
	Preload  string `json:"preload" form:"preload"`     // 预加载关联表名, 若多个以英文逗号分隔
}

type QueryOption struct {
	Select   []string         // 指定查询字段
	Filter   map[string]any   // 简单查询条件
	Where    map[string]any   // 自定义复杂条件
	Preload  map[string][]any // 预加载关联查询
	OrderBy  string           // 排序, eg: "create_time desc, update_time"
	Limit    int              // QueryAll时限定查询条数
	Page     int              // QueryPage时指定查询页数
	PageSize int              // QueryPage时指定每页数量
}

func NewQuery[T any](options ...*QueryOption) (*gorm.DB, *QueryOption) {
	var opt *QueryOption
	if len(options) > 0 {
		opt = options[0]
	} else {
		opt = &QueryOption{}
	}
	if opt.Page < 1 {
		opt.Page = 1
	}
	if opt.PageSize < 1 {
		opt.PageSize = Option.DefaultPageSize
	}
	query := DB.Model(new(T))
	query = setQuerySelect(query, opt.Select)
	query = setQueryFilter(query, opt.Filter)
	query = setQueryWhere(query, opt.Where)
	query = setQueryPreload(query, opt.Preload)
	query = setQueryOrderBy(query, opt.OrderBy)
	query = setQueryLimit(query, opt.Limit)
	return query, opt
}

func setQuerySelect(query *gorm.DB, fields []string) *gorm.DB {
	if len(fields) > 0 {
		query.Select(fields)
	}
	return query
}

func setQueryFilter(query *gorm.DB, filterMap map[string]any) *gorm.DB {
	for k, v := range filterMap {
		ks := strings.Split(k, ":")
		if len(ks) > 1 {
			query = setQueryFilterOperator(query, ks[0], ks[1], v)
		} else {
			query.Where(k, v)
		}
	}
	return query
}

func setQueryWhere(query *gorm.DB, whereMap map[string]any) *gorm.DB {
	for k, v := range whereMap {
		if a, ok := v.([]any); ok {
			query.Where(k, a...)
		} else {
			query.Where(k, v)
		}
	}
	return query
}

func setQueryPreload(query *gorm.DB, preloadMap map[string][]any) *gorm.DB {
	for k, v := range preloadMap {
		query.Preload(k, v...)
	}
	return query
}

func setQueryOrderBy(query *gorm.DB, orderBy string) *gorm.DB {
	if orderBy != "" {
		query.Order(orderBy)
	}
	return query
}

func setQueryLimit(query *gorm.DB, limit int) *gorm.DB {
	if limit > 0 {
		query.Limit(limit)
	}
	return query
}

func setQueryPage(query *gorm.DB, page int, pageSize int) *gorm.DB {
	if pageSize > 0 {
		query.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	return query
}

func setQueryFilterOperator(q *gorm.DB, key string, operater string, val any) *gorm.DB {
	switch operater {
	case "":
		return q.Where(map[string]any{key: val})
	case "eq":
		return q.Where(fmt.Sprintf("`%s` = ?", key), val)
	case "ne":
		return q.Where(fmt.Sprintf("`%s` != ?", key), val)
	case "gt":
		return q.Where(fmt.Sprintf("`%s` > ?", key), val)
	case "ge":
		return q.Where(fmt.Sprintf("`%s` >= ?", key), val)
	case "lt":
		return q.Where(fmt.Sprintf("`%s` < ?", key), val)
	case "le":
		return q.Where(fmt.Sprintf("`%s` <= ?", key), val)
	case "in":
		vType := reflect.TypeOf(val).Kind()
		if vType == reflect.Ptr {
			vType = reflect.ValueOf(val).Elem().Kind()
		}
		if vType == reflect.String {
			return q.Where(fmt.Sprintf("`%s` in ?", key), strings.Split(val.(string), ","))
		} else if vType == reflect.Slice {
			return q.Where(fmt.Sprintf("`%s` in ?", key), val)
		} else {
			return q.Where(fmt.Sprintf("`%s` in (?)", key), val)
		}
	case "ni":
		if reflect.TypeOf(val).Kind() == reflect.Slice {
			return q.Where(fmt.Sprintf("`%s` not in ?", key), val)
		} else {
			return q.Where(fmt.Sprintf("`%s` not in ?", key), strings.Split(val.(string), ","))
		}
	case "ct":
		return q.Where(fmt.Sprintf("`%s` like '%%%s%%'", key, val))
	case "nc":
		return q.Where(fmt.Sprintf("`%s` not like '%%%s%%'", key, val))
	case "sw":
		return q.Where(fmt.Sprintf("`%s` like '%s%%'", key, val))
	case "ew":
		return q.Where(fmt.Sprintf("`%s` like '%%%s'", key, val))
	case "ob":
		if !arrayx.Contains([]string{"asc", "desc", "ASC", "DESC"}, val.(string)) {
			val = "asc"
		}
		return q.Order(fmt.Sprintf("%s %s", key, val))
	default:
		return q.Where(map[string]any{key: val})
	}
}
