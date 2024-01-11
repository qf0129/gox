package gox

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/constx"
	"github.com/qf0129/gox/gormx"
	"github.com/qf0129/gox/tmplx"
)

type Option struct {
	GinEngine      *gin.Engine
	HttpServer     *http.Server
	GromxOption    *gormx.Option
	EnableTemplate bool
}

func Run(opt *Option) {
	if opt.GinEngine == nil {
		panic("RequiredHttpHandler")
	}
	if opt.EnableTemplate {
		if opt.GromxOption == nil {
			opt.GromxOption = &gormx.Option{}
		}
		if opt.GromxOption.Models == nil {
			opt.GromxOption.Models = []any{}
		}
		opt.GromxOption.Models = append(opt.GromxOption.Models, &tmplx.Api{}, &tmplx.User{}, &tmplx.Role{})
		gormx.Connect(opt.GromxOption)
		tmplx.UploadApis(opt.GinEngine)
		slog.Info("### EnableTemplate true")
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
		opt.HttpServer.Addr = constx.DefaultListenAddr
	}
	if opt.HttpServer.ReadTimeout == 0 {
		opt.HttpServer.ReadTimeout = constx.DefaultReadTimeout * time.Second
	}
	if opt.HttpServer.WriteTimeout == 0 {
		opt.HttpServer.WriteTimeout = constx.DefaultWriteTimeout * time.Second
	}
}
