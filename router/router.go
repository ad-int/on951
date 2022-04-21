package router

import (
	"github.com/gin-gonic/gin"
)

type TAppRouter struct {
	engine *gin.Engine
}

type TRoutesList *map[string]TRouteDescription

type TRouteDescription struct {
	Name        string
	Method      string
	Handler     gin.HandlerFunc
	Group       string
	Middlewares []gin.HandlerFunc
}

type AppRouter interface {
	InitRoutes(routes *map[string]TRoutesList)
	GetEngine() *gin.Engine
	Configure()
	Run() error
}

func (appRouter *TAppRouter) Configure() {
	appRouter.engine = gin.Default()
	appRouter.engine.SetTrustedProxies([]string{"127.0.0.1"})
}

func (appRouter *TAppRouter) InitRoutes(routes *map[string]TRoutesList) {
	var handle gin.IRoutes
	for method, routeList := range *routes {
		for path, routeDescription := range *routeList {
			if routeDescription.Group != "" {
				handle = appRouter.engine.Group(routeDescription.Group)
			} else {
				handle = appRouter.engine
			}
			thisHandle := handle.Handle(method, path, routeDescription.Handler)
			if len(routeDescription.Middlewares) > 0 {
				thisHandle.Use(routeDescription.Middlewares...)
			}
		}
	}
}

func (appRouter *TAppRouter) GetEngine() *gin.Engine {
	return appRouter.engine
}

func (appRouter *TAppRouter) Run() error {
	return appRouter.engine.Run()
}
