package daox

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/qf0129/gox/arrayx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FilteFunc func(tx *gorm.DB) *gorm.DB

func ParseUrlFilters(vals url.Values) map[string]any {
	filters := map[string]any{}
	for k := range vals {
		if !arrayx.HasStrItem(FIXED_KEYS, k) {
			filters[k] = vals.Get(k)
		}
	}
	return filters
}

func ParseFilters(filters map[string]any) (funcs []FilteFunc) {
	for k, v := range filters {
		ks := strings.Split(k, ":")
		if len(ks) == 2 {
			funcs = append(funcs, FilteKeyFunc(ks[0], ks[1], v))
		} else {
			funcs = append(funcs, FilteKeyFunc(k, "eq", v))
		}
	}
	return
}

func FilteKeyFunc(key string, operater string, val any) FilteFunc {
	return func(tx *gorm.DB) *gorm.DB {
		switch operater {
		case "eq":
			return tx.Where(fmt.Sprintf("`%s` = ?", key), val)
		case "ne":
			return tx.Where(fmt.Sprintf("`%s` != ?", key), val)
		case "gt":
			return tx.Where(fmt.Sprintf("`%s` > ?", key), val)
		case "ge":
			return tx.Where(fmt.Sprintf("`%s` >= ?", key), val)
		case "lt":
			return tx.Where(fmt.Sprintf("`%s` < ?", key), val)
		case "le":
			return tx.Where(fmt.Sprintf("`%s` <= ?", key), val)
		case "in":
			if reflect.TypeOf(val).Kind() == reflect.Slice {
				return tx.Where(fmt.Sprintf("`%s` in ?", key), val)
			} else {
				return tx.Where(fmt.Sprintf("`%s` in ?", key), strings.Split(val.(string), ","))
			}
		case "ni":
			if reflect.TypeOf(val).Kind() == reflect.Slice {
				return tx.Where(fmt.Sprintf("`%s` not in ?", key), val)
			} else {
				return tx.Where(fmt.Sprintf("`%s` not in ?", key), strings.Split(val.(string), ","))
			}
		case "ct":
			return tx.Where(fmt.Sprintf("`%s` like '%%%s%%'", key, val))
		case "nc":
			return tx.Where(fmt.Sprintf("`%s` not like '%%%s%%'", key, val))
		case "sw":
			return tx.Where(fmt.Sprintf("`%s` like '%s%%'", key, val))
		case "ew":
			return tx.Where(fmt.Sprintf("`%s` like '%%%s'", key, val))
		case "ob":
			if !arrayx.HasStrItem([]string{"asc", "desc", "ASC", "DESC"}, val) {
				val = "asc"
			}
			return tx.Order(fmt.Sprintf("%s %s", key, val))
		default:
			return tx.Where("? = '?'", key, val)
		}
	}
}

func FiltePreloadFunc(field string, funcs ...FilteFunc) FilteFunc {
	return func(tx *gorm.DB) *gorm.DB {
		if field == "" {
			return tx
		} else if field == "*" {
			return tx.Preload(clause.Associations)
		} else {
			return tx.Preload(cases.Title(language.Dutch).String(field), func(tx *gorm.DB) *gorm.DB {
				for _, fc := range funcs {
					tx = fc(tx)
				}
				return tx
			})
		}
	}
}

func FiltePageFunc(page int, pageSize int) FilteFunc {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Limit(pageSize).Offset((page - 1) * pageSize)
	}
}

func FilteWhereFunc(query any, args ...any) FilteFunc {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(query, args...)
	}
}

func FilteKVFunc(k string, v any) FilteFunc {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(map[string]any{k: v})
	}
}
