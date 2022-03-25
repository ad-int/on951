package state

import (
	"github.com/joho/godotenv"
	"log"
	"main/router"
)
var application *TApplication


type TApplication struct {
	config map[string]string
	router router.TAppRouter
}

func GetApplication() *TApplication {
	if application == nil {
		application = &TApplication{}
		return application
	}
	return application
}


func (app *TApplication) GetConfigValue(key string) string {
	return app.config[key]
}

func (app *TApplication) Init(routes *map[string]router.TRoutesList) {

	var err error
	app.config, err = godotenv.Read(".env")
	if err != nil {
		log.Println(err)
		log.Fatalln("Unable to read .env file!")
	}

	log.Println(app.config)

	app.router.Configure()
	app.router.InitRoutes(routes)
	app.router.Run()
}