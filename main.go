package main

import (
	"on951/application"
	"on951/data_generator"
	"strconv"
)

func main() {
	app := application.GetApplicationRepository()
	app.Bootstrap(&DefinedRoutes, func() {
		for t := 0; t < 20; t++ { // Generating random articles
			data_generator.GenerateArticle()
		}
		cost, _ := strconv.Atoi(application.GetApplication().GetConfigValue("BCRYPT_HASH_GENERATION_COST"))
		_ = data_generator.GenerateUser("user1", "password1", cost)
	})
}
