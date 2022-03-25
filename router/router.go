package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type TAppRouter struct {
	engine *gin.Engine
}

type TRoutesList *map[string]TRouteDescription

type TRouteDescription struct {
	Name string
	Method string
	Handler gin.HandlerFunc
	Group       string
	Middlewares []gin.HandlerFunc
}

type AppRouter interface {
	InitRoutes(routes *map[string]TRoutesList)
	Configure()
	Run()
}

func (appRouter *TAppRouter) Configure() {
	appRouter.engine = gin.Default()
	appRouter.engine.SetTrustedProxies([]string{"127.0.0.1"})
}

func (appRouter *TAppRouter) InitRoutes(routes *map[string]TRoutesList) {
	var handle gin.IRoutes
	for method, routeList := range *routes {
		fmt.Println(method)
		for path, routeDescription := range *routeList {
			if routeDescription.Group != "" {
				handle = appRouter.engine.Group(routeDescription.Group)
			} else {
				handle = appRouter.engine
			}
			if len(routeDescription.Middlewares) > 0 {
				log.Println(routeDescription.Middlewares)
				handle.Use(routeDescription.Middlewares...)
			}
			handle.Handle(method, path, routeDescription.Handler)
		}
	}
}

func (appRouter *TAppRouter) Run() {
	appRouter.engine.Run()
}
