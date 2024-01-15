package dbx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/qf0129/gox/arrayx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type filteFunc func(q *gorm.DB) *gorm.DB

func paramsMapToFilters(params map[string]any) (funcs []filteFunc) {
	for k, v := range params {
		ks := strings.Split(k, ":")
		if len(ks) > 1 {
			funcs = append(funcs, filteKeyFunc(ks[0], ks[1], v))
		} else {
			funcs = append(funcs, filteKeyFunc(k, "", v))
		}
	}
	return
}

func filteKeyFunc(key string, operater string, val any) filteFunc {
	return func(q *gorm.DB) *gorm.DB {
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
			if reflect.TypeOf(val).Kind() == reflect.Slice {
				return q.Where(fmt.Sprintf("`%s` in ?", key), val)
			} else {
				return q.Where(fmt.Sprintf("`%s` in ?", key), strings.Split(val.(string), ","))
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
			if !arrayx.HasStrItem([]string{"asc", "desc", "ASC", "DESC"}, val) {
				val = "asc"
			}
			return q.Order(fmt.Sprintf("%s %s", key, val))
		default:
			return q.Where(map[string]any{key: val})
		}
	}
}

func FiltePreloadFunc(field string, funcs ...filteFunc) filteFunc {
	return func(q *gorm.DB) *gorm.DB {
		if field == "" {
			return q
		} else if field == "*" {
			return q.Preload(clause.Associations)
		} else {
			return q.Preload(cases.Title(language.Dutch).String(field), func(q *gorm.DB) *gorm.DB {
				for _, fc := range funcs {
					q = fc(q)
				}
				return q
			})
		}
	}
}

func FiltePageFunc(page int, pageSize int) filteFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Limit(pageSize).Offset((page - 1) * pageSize)
	}
}

func FilteWhereFunc(query any, args ...any) filteFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where(query, args...)
	}
}

func FilteKVFunc(k string, v any) filteFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where(map[string]any{k: v})
	}
}
