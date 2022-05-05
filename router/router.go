package router

import (
	"github.com/gin-contrib/cors"
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

func (appRouter *TAppRouter) Configure(env string, trustedProxies []string, allowedHeaders []string, allowAllOrigins bool) error {
	if env != "" {
		gin.SetMode(env)
	}
	appRouter.engine = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = allowedHeaders
	corsConfig.AllowAllOrigins = allowAllOrigins
	appRouter.engine.Use(cors.New(corsConfig))
	return appRouter.engine.SetTrustedProxies(trustedProxies)
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
			if len(routeDescription.Middlewares) > 0 {
				handle.Use(routeDescription.Middlewares...)
			}
			handle.Handle(method, path, routeDescription.Handler)
		}
	}
}

func (appRouter *TAppRouter) GetEngine() *gin.Engine {
	return appRouter.engine
}

func (appRouter *TAppRouter) Run(addr ...string) error {
	return appRouter.engine.Run(addr...)
}
