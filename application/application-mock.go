package application

import (
	"github.com/stretchr/testify/mock"
	"log"
	"on951/router"
)

type TApplicationRepositoryMock struct {
	mock.Mock
	Application IApplication
}

type TApplicationMock struct {
	mock.Mock
	TApplication
}

func (app *TApplicationRepositoryMock) GetApplication() IApplication {
	args := app.Called()
	return args.Get(0).(IApplication)
}

func (app *TApplicationMock) ReadEnvFile() bool {
	args := app.Called()
	returnValue := args.Bool(0)
	if returnValue {
		app.config = map[string]string{
			"DSN":      "",
			"SECRET":   "213243fdessf",
			"ISSUER":   "localhost",
			"AUDIENCE": "general",
		}
	}
	log.Println(app.config)
	return returnValue
}

func (app *TApplicationMock) GetImagesDir() string {
	args := app.Called()
	return args.String(0)
}

func (app *TApplicationMock) Init(routes *map[string]router.TRoutesList) error {
	args := app.Called()
	if len(app.GetImagesDir()) < 1 {
		log.Fatalln("Cannot read images directory path")
	}

	app.router.Configure()
	app.router.InitRoutes(routes)
	return args.Error(0)
}
