package application

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"on951/api"
	"on951/router"
	"testing"
)

type applicationTestSuite struct {
	suite.Suite
	config       map[string]string
	router       router.TAppRouter
	articlesRepo *api.TArticlesRepository
}

func TestApplicationTestSuite(t *testing.T) {
	suite.Run(t, new(applicationTestSuite))
}

func (suite *applicationTestSuite) TestApplication() {
	suite.IsType(&TApplication{}, GetApplication())
}

func (suite *applicationTestSuite) TestApplicationBootstrap() {

	Bootstrap(&map[string]router.TRoutesList{
		"GET": {
			"health-check": {
				Name: "Health Check",
				Handler: func(context *gin.Context) {

				},
			},
		},
		"PUT": {
			"comment/:article_id": {
				Name: "Post a comment",
				Handler: func(context *gin.Context) {

				},
			},
		},
	})

	suite.IsType(&TApplication{}, GetApplication())
}
