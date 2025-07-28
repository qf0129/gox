package serverx

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func NewApp(cfgs ...*Config) *App {
	var cfg *Config
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}
	cfg = loadDefaultConfig(cfg)
	gin.SetMode(cfg.GinMode)
	ginEngine := gin.Default()
	if cfg.EnableCheckHealthApi {
		ginEngine.GET(cfg.CheckHealthApiPath, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}
	app := &App{
		Engine:    ginEngine,
		Config:    cfg,
		ApiGroups: []*ApiGroup{{Path: "/", Name: "default"}},
	}
	return app
}

type App struct {
	Engine      *gin.Engine
	Config      *Config
	ApiGroups   []*ApiGroup
	Middlewares []gin.HandlerFunc
}

func (app *App) UseCORS() *App {
	return app.Use(MiddlewareCORS())
}

func (app *App) AddGroup(groups ...*ApiGroup) *App {
	return app.AddGroups(groups)
}

func (app *App) AddGroups(groups []*ApiGroup) *App {
	app.ApiGroups = append(app.ApiGroups, groups...)
	return app
}

func (app *App) AddApi(apis ...*ApiInfo) *App {
	if len(app.ApiGroups) == 0 {
		app.addDefaultGroup()
	}
	app.ApiGroups[0].Add(apis...)
	return app
}

func (app *App) Use(middlewares ...gin.HandlerFunc) *App {
	if len(app.ApiGroups) == 0 {
		app.addDefaultGroup()
	}
	app.Middlewares = append(app.Middlewares, middlewares...)
	return app
}

func (app *App) Run() {
	for _, apiGroup := range app.ApiGroups {
		if len(apiGroup.Apis) > 0 {
			ginGroup := app.Engine.Group(apiGroup.Path)
			if app.Config.EnableRequestId {
				ginGroup.Use(func(c *gin.Context) { c.Set(KeyOfRequestId, xid.New().String()) })
			}
			// ginGroup.Use(app.Middlewares...)
			ginGroup.Use(apiGroup.Middlewares...)
			for _, api := range apiGroup.Apis {
				api.loadDefaut()
				if api.GinHandler != nil {
					ginGroup.Handle(api.Method, api.Path, api.GinHandler)
				} else {
					if api.Handler == nil {
						log.Fatalf("api \"%s\" handler is nil", api.Path)
					}
					ginGroup.Handle(api.Method, api.Path, api.handle(app))
				}
			}
		}
	}
	app.runServer()
}

func (app *App) runServer() {
	slog.Info("### Server listening on " + app.Config.Addr)
	app.Engine.Use(app.Middlewares...)
	server := &http.Server{
		Handler:      app.Engine,
		Addr:         app.Config.Addr,
		ReadTimeout:  time.Duration(app.Config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.Config.WriteTimeout) * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		slog.Error(err.Error())
	}
}

func (app *App) addDefaultGroup() {
	app.ApiGroups = append(app.ApiGroups, &ApiGroup{Path: "/", Name: "default"})
}
