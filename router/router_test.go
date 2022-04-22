package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type routerTestSuite struct {
	suite.Suite
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(routerTestSuite))
}

func (suite *routerTestSuite) TestInitRoutes() {
	for _, testCase := range testRouterData {
		router := &TAppRouter{}

		configurationError := router.Configure(
			strings.Fields(testCase.config["TRUSTED_PROXIES"]),
			strings.Fields(testCase.config["CORS_ALLOWED_HEADERS"]),
			testCase.config["CORS_ALLOW_ALL_ORIGINS"] == "true",
		)
		if testCase.configurationFails {
			suite.NotNil(configurationError)
			continue
		} else {
			suite.Nil(configurationError)
		}
		router.InitRoutes(&testCase.routes)
		routesFound := 0

		for method, routeDescription := range testCase.routes {
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
}

func (suite *routerTestSuite) TestGetEngine() {
	for _, testCase := range testRouterData {
		router := &TAppRouter{}
		configurationError := router.Configure(
			strings.Fields(testCase.config["TRUSTED_PROXIES"]),
			strings.Fields(testCase.config["CORS_ALLOWED_HEADERS"]),
			testCase.config["CORS_ALLOW_ALL_ORIGINS"] == "true",
		)
		if testCase.configurationFails {
			suite.NotNil(configurationError)
			continue
		} else {
			suite.Nil(configurationError)
		}
		suite.IsType(&gin.Engine{}, router.GetEngine())
	}

}
