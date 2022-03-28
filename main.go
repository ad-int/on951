package main

import (
	"github.com/gin-gonic/gin"
	"main/database"
	"main/router"
	"main/state"
	"main/web/handlers"
	"main/web/middleware"
	"net/http"
)

func main() {

	app := state.GetApplication()

	database.ConnectToDB()
	database.AutoMigrate()
	handlers.Generate()
	app.Init(&map[string]router.TRoutesList{
		http.MethodGet: {
			"health-check": {
				Name:    "Health Check",
				Handler: handlers.HealthCheck,
			},
			"token": {
				Name:    "Get token",
				Handler: handlers.GetToken,
			},
			"articles": {
				Name:        "List articles",
				Handler:     handlers.GetArticles,
				Middlewares: []gin.HandlerFunc{middleware.ApiAuthCheck},
			},
			"article": {
				Name:        "Get specific article",
				Handler:     handlers.GetArticle,
				Middlewares: []gin.HandlerFunc{middleware.ApiAuthCheck},
			},
		},
	})
}
