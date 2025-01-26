package ginx

import "github.com/gin-gonic/gin"

type ApiGroup struct {
	Name        string
	Path        string
	Apis        []*Api
	Middlewares []gin.HandlerFunc
}

func NewApiGroup(name, path string, apis ...*Api) *ApiGroup {
	g := &ApiGroup{
		Name: name,
		Path: path,
	}
	for _, api := range apis {
		g.AddApi(api)
	}
	return g
}

func (g *ApiGroup) Use(middlewares ...gin.HandlerFunc) *ApiGroup {
	g.Middlewares = append(g.Middlewares, middlewares...)
	return g
}

func (g *ApiGroup) AddApi(apis ...*Api) *ApiGroup {
	g.Apis = append(g.Apis, apis...)
	return g
}
