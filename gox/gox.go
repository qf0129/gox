package gox

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/gormx"
	"github.com/qf0129/gox/gormx/daox"
	"github.com/qf0129/gox/modelx"
)

type Option struct {
	GinEngine        *gin.Engine
	HttpServer       *http.Server
	GromxOption      *gormx.Option
	EnablePermission bool
}

func Run(opt *Option) {
	if opt.GinEngine == nil {
		panic("RequiredHttpHandler")
	}

	if opt.EnablePermission {
		if opt.GromxOption == nil {
			opt.GromxOption = &gormx.Option{}
		}
		if opt.GromxOption.Models == nil {
			opt.GromxOption.Models = []any{}
		}
		opt.GromxOption.Models = append(opt.GromxOption.Models, &modelx.Api{}, &modelx.User{}, &modelx.Role{})
		gormx.Connect(opt.GromxOption)
		uploadApis(opt.GinEngine)
		slog.Info("### EnablePermission: true")
	} else {
		if opt.GromxOption != nil {
			gormx.Connect(opt.GromxOption)
		}
	}

	initHttpServer(opt)
	slog.Info("### server is listening " + opt.HttpServer.Addr)
	opt.HttpServer.ListenAndServe()
}

func initHttpServer(opt *Option) {
	if opt.HttpServer == nil {
		opt.HttpServer = &http.Server{}
	}
	opt.HttpServer.Handler = opt.GinEngine
	if opt.HttpServer.Addr == "" {
		opt.HttpServer.Addr = ":8080"
	}
	if opt.HttpServer.ReadTimeout == 0 {
		opt.HttpServer.ReadTimeout = 60 * time.Second
	}
	if opt.HttpServer.WriteTimeout == 0 {
		opt.HttpServer.WriteTimeout = 60 * time.Second
	}
}

func uploadApis(engine *gin.Engine) {
	apis := []modelx.Api{}
	for _, route := range engine.Routes() {
		key := concatApiKey(&route)
		api, _ := daox.QueryOneByMap[modelx.Api](map[string]any{"key": key})
		if api.Id == "" {
			apis = append(apis, modelx.Api{Key: key, Method: route.Method, Path: route.Path})
		}
	}

	if len(apis) > 0 {
		gormx.DB.Save(apis)
	}
	slog.Info("### UpdateApis", slog.Int("len", len(apis)))
}

func concatApiKey(route *gin.RouteInfo) string {
	return route.Method + "|" + route.Path
}
