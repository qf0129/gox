package tmplx

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/gormx"
	"github.com/qf0129/gox/gormx/daox"
)

func UploadApis(engine *gin.Engine) {
	apis := []Api{}
	for _, route := range engine.Routes() {
		key := ConcatApiKey(&route)
		api, _ := daox.QueryOneByMap[Api](map[string]any{"key": key})
		if api.Id == "" {
			apis = append(apis, Api{Key: key, Method: route.Method, Path: route.Path})
		}
	}

	if len(apis) > 0 {
		gormx.DB.Save(apis)
	}
	slog.Info("### UpdateApis", slog.Int("len", len(apis)))
}

func ConcatApiKey(route *gin.RouteInfo) string {
	return route.Method + "|" + route.Path
}
