package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"testing"
)

type routerTestSuite struct {
	suite.Suite
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(routerTestSuite))
}

func (suite *routerTestSuite) TestInitRoutes() {

	middlewareFunc := func(context *gin.Context) {

	}

	handlerFunc1 := func(context *gin.Context) {
	}
	handlerFunc2 := func(context *gin.Context) {

	}
	router := &TAppRouter{}
	routesList := &map[string]TRoutesList{
		"GET": {
			"health-check": {
				Name:    "Health Check",
				Handler: handlerFunc1,
			},
		},
		"PUT": {
			"comment/:article_id": {
				Name:        "Post a comment",
				Handler:     handlerFunc2,
				Middlewares: []gin.HandlerFunc{middlewareFunc},
			},
		},
	}
	suite.Nil(router.Configure())
	router.InitRoutes(routesList)
	routesFound := 0

	for method, routeDescription := range *routesList {
		for path, _ := range *routeDescription {
			for _, h := range router.GetEngine().Routes() {
				if h.Method == method && h.Path == "/"+path {

					routesFound++
				}
			}
		}
	}
	suite.Equal(len(router.GetEngine().Routes()), routesFound)
}

func (suite *routerTestSuite) TestGetEngine() {
	router := &TAppRouter{}
	suite.Nil(router.Configure())
	suite.IsType(&gin.Engine{}, router.GetEngine())
}
