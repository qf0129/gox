package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/errx"
)

type Api struct {
	Module      string
	Method      string
	Name        string
	Path        string
	Type        string
	Description string
	Handler     HandlerFunc
	GinHandler  gin.HandlerFunc
}
type HandlerFunc func(c *gin.Context) (interface{}, errx.Err)

func NewApi(name string, method string, path string, handler HandlerFunc) Api {
	return Api{
		Name:    name,
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}

func (api *Api) loadDefaut() *Api {
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

func (api *Api) handle(enableRequestId bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqId := ""
		if enableRequestId {
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
