package application

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"on951/router"
	"strings"
)

type TApplicationMock struct {
	mock.Mock
	TApplication
}

func (app *TApplicationMock) ReadEnvFile() (bool, map[string]string) {
	args := app.Called()
	app.config = args.Get(1).(map[string]string)
	return args.Bool(0), args.Get(1).(map[string]string)
}

func (app *TApplicationMock) GetImagesDir() string {
	args := app.Called()
	return args.String(0)
}

func (app *TApplicationMock) GetConfigValue(key string) string {
	args := app.Called(key)
	return args.String(0)
}

func (app *TApplicationMock) Init(routes *map[string]router.TRoutesList) error {
	args := app.Called()
	if len(app.GetImagesDir()) < 1 {
		panic(errors.New("Cannot read images directory path"))
	}

	err := app.router.Configure(
		strings.Fields(app.GetConfigValue("TRUSTED_PROXIES")),
		strings.Fields(app.GetConfigValue("CORS_ALLOWED_HEADERS")),
		app.GetConfigValue("CORS_ALLOW_ALL_ORIGINS") == "true",
	)
	if err != nil {
		return err
	}
	app.router.InitRoutes(routes)
	return args.Error(0)
}
func (app *TApplicationMock) GetRouter() router.AppRouter {
	args := app.Called()
	return args.Get(0).(router.AppRouter)
}
