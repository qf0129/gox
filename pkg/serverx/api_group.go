package serverx

import (
	"github.com/gin-gonic/gin"
)

type ApiGroup struct {
	Path        string
	Name        string
	Apis        []*ApiInfo
	Middlewares []gin.HandlerFunc
}

func Group(path string, apis ...*ApiInfo) *ApiGroup {
	g := &ApiGroup{Path: path}
	for _, api := range apis {
		g.Add(api)
	}
	return g
}

func (g *ApiGroup) Use(middlewares ...gin.HandlerFunc) *ApiGroup {
	g.Middlewares = append(g.Middlewares, middlewares...)
	return g
}

func (g *ApiGroup) Add(apis ...*ApiInfo) *ApiGroup {
	g.Apis = append(g.Apis, apis...)
	return g
}

func (g *ApiGroup) AddMap(method string, m map[string]HandlerFunc) *ApiGroup {
	for path, api := range m {
		g.Add(Api(method, path, api))
	}
	return g
}
