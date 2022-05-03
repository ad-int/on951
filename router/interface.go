package router

import "github.com/gin-gonic/gin"

type AppRouter interface {
	InitRoutes(routes *map[string]TRoutesList)
	GetEngine() *gin.Engine
	Configure(env string, trustedProxies []string, allowedHeaders []string, allowAllOrigins bool) error
	Run(addr ...string) error
}
