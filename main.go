package main

import (
	"log"
	"on951/application"
	"on951/web/handlers"
)

func main() {
	app := application.GetApplicationRepository()
	if len(app.GetApplication().GetImagesDir()) < 1 {
		log.Fatalln("Cannot read images directory path")
	}
	app.Bootstrap(&DefinedRoutes, func() {
		for t := 0; t < 20; t++ { // Generating random articles
			handlers.Generate()
		}
	})
}
