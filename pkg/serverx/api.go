package serverx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/errx"
)

type ApiInfo struct {
	Method      string
	Path        string
	Handler     HandlerFunc
	Module      string
	Name        string
	Type        string
	Description string
	GinHandler  gin.HandlerFunc
}
type HandlerFunc func(c *gin.Context) (interface{}, errx.Err)

func Api(method string, path string, handler HandlerFunc) *ApiInfo {
	return &ApiInfo{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}

func (api *ApiInfo) loadDefaut() *ApiInfo {
	if api.Method == "" {
		api.Method = "GET"
	}
	if api.Path == "" {
		api.Path = api.Name
	}
	if api.Type == "" {
		api.Type = api.Method
	}
	return api
}

func (api *ApiInfo) handle(app *App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqId := ""
		if app.Config.EnableRequestId {
			reqId = ctx.GetString(KeyOfRequestId)
		}
		rsp, err := api.Handler(ctx)
		if err != nil {
			Response(ctx, http.StatusOK, err.Code(), err.Msg(), nil, reqId)
		} else {
			Response(ctx, http.StatusOK, errx.Success.Code(), errx.Success.Msg(), rsp, reqId)
		}
	}
}
