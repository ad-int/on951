package application

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"log"
	"on951/api"
	"on951/router"
	"reflect"
	"testing"
)

const fakeImagesDir = "path/fake-images-dir"

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

	app := &TApplicationMock{}

	app.On("ReadEnvFile").Return(true)
	app.On("GetImagesDir").Return(fakeImagesDir)
	app.On("Init").Return(nil)
	appRepo := &TApplicationRepository{application: app}
	log.Println(reflect.TypeOf(appRepo.GetApplication()))
	routesList := &map[string]router.TRoutesList{
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
	}
	appRepo.Bootstrap(routesList)
	routesFound := 0
	suite.IsType(&TApplicationMock{}, appRepo.GetApplication())

	for method, routeDescription := range *routesList {
		for path, _ := range *routeDescription {
			for _, h := range appRepo.GetApplication().GetRouter().GetEngine().Routes() {
				if h.Method == method && h.Path == "/"+path {
					routesFound++
				}
			}
		}
	}
	suite.Equal(fakeImagesDir, appRepo.GetApplication().GetImagesDir())
	suite.Equal(len(appRepo.GetApplication().GetRouter().GetEngine().Routes()), routesFound)
}
