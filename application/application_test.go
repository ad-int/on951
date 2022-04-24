package application

import (
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
	router       router.TAppRouterMock
	articlesRepo *api.TArticlesRepository
}

func TestApplicationTestSuite(t *testing.T) {
	suite.Run(t, new(applicationTestSuite))
}

func (suite *applicationTestSuite) TestGetApplication() {
	suite.IsType(&TApplication{}, GetApplication())
}

func (suite *applicationTestSuite) TestApplicationBootstrap() {

	for _, testCase := range testApplicationData {

		app := &TApplicationMock{}
		app.On("ReadEnvFile").Return(len(app.config) > 0, testCase.config)
		app.On("GetConfigValue", "TRUSTED_PROXIES").Return(testCase.config["TRUSTED_PROXIES"])
		app.On("GetConfigValue", "CORS_ALLOWED_HEADERS").Return(testCase.config["CORS_ALLOWED_HEADERS"])
		app.On("GetConfigValue", "CORS_ALLOW_ALL_ORIGINS").Return(testCase.config["CORS_ALLOW_ALL_ORIGINS"])
		app.On("GetImagesDir").Return(testCase.imagesDir)
		app.On("GetRouter").Return(&app.router)

		app.On("Init").Return(nil)
		appRepo := &TApplicationRepository{application: app}
		log.Println(reflect.TypeOf(appRepo.GetApplication().GetRouter()))

		suite.IsType(&TApplicationMock{}, appRepo.GetApplication())
		suite.Equal(testCase.imagesDir, appRepo.GetApplication().GetImagesDir())

		if len(testCase.imagesDir) == 0 {
			suite.PanicsWithError("Cannot read images directory path", func() {
				appRepo.Bootstrap(&testCase.routes)

			})
			continue
		} else if testCase.configurationFails {
			suite.Panics(func() {
				appRepo.Bootstrap(&testCase.routes)
			})
			continue
		} else {
			appRepo.Bootstrap(&testCase.routes)

		}

		routesFound := 0

		for method, routeDescription := range testCase.routes {
			for path, _ := range *routeDescription {
				for _, h := range appRepo.GetApplication().GetRouter().GetEngine().Routes() {
					if h.Method == method && h.Path == "/"+path {
						routesFound++
						suite.T().Logf("Added route: %v %v\n", method, path)
					}
				}
			}
		}
		suite.Equal(len(appRepo.GetApplication().GetRouter().GetEngine().Routes()), routesFound)
	}
}
