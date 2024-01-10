package daox

import (
	"strings"
)

const (
	FixedKeyPage        = "Page"
	FixedKeyPageSize    = "PageSize"
	FixedKeyPreload     = "Preload"
	FixedKeyClosePaging = "ClosePaging"
)

var FIXED_KEYS = []string{FixedKeyPage, FixedKeyPageSize, FixedKeyPreload, FixedKeyClosePaging}

type PageBody[T any] struct {
	List     []T   `json:"list"`
	PageNum  int   `json:"page_num"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

// 查询结构体
type QueryBody struct {
	PageNum   int            `json:"page_num" form:"page_num"`     // 页数, 默认1
	PageSize  int            `json:"page_size" form:"page_size"`   // 每页数量, 默认10
	Preload   string         `json:"preload" form:"preload"`       // 预加载关联表名, 若多个以英文逗号分隔
	NoPaging  bool           `json:"no_paging" form:"no_paging"`   // 关闭分页, 默认false
	Filter    string         `json:"filter" form:"filter"`         // 过滤条件, 'key1:value1|key2:value2'
	FilterMap map[string]any `json:"filter_map" form:"filter_map"` // 过滤条件map
	Fields    []string       `json:"fields" form:"fields"`         // 查询字段
}

func (query *QueryBody) ParseFilterToMap() {
	if query.FilterMap == nil {
		query.FilterMap = map[string]any{}
	}
	if query.Filter != "" {
		filterList := strings.Split(query.Filter, "|")
		for _, filter := range filterList {
			items := strings.Split(filter, ":")
			if len(items) == 2 {
				query.FilterMap[items[0]] = items[1]
			} else if len(items) >= 3 {
				query.FilterMap[items[0]+":"+items[1]] = items[2]
			}
		}
	}
}
