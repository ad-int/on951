package main

import (
	"on951/application"
	"on951/web/handlers"
)

func main() {
	app := application.GetApplicationRepository()
	app.Bootstrap(&DefinedRoutes, func() {
		for t := 0; t < 20; t++ { // Generating random articles
			handlers.Generate()
		}
	})
}
