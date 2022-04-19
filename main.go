package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/router"
	"main/state"
	"main/web/handlers"
	"main/web/middleware"
	"net/http"
)

func main() {

	app := state.GetApplication()

	//ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	//ctx.Set("article_id1", 1)
	//handlers.GetComments(ctx)
	//image_links_parser.Process(`ikljghi  <img src="data:application/json;base64," /> 232refwdsf <img src="data:image/png;Es!a, QQQ" />grd`)
	//return
	app.ReadEnvFile()
	app.InitDb()

	for t := 0; t < 50; t++ { // Generating random articles
		handlers.Generate()
	}

	err := app.Init(&map[string]router.TRoutesList{
		http.MethodGet: {
			"health-check": {
				Name:    "Health Check",
				Handler: handlers.HealthCheck,
			},
			"token": {
				Name:    "Get token",
				Handler: handlers.GetToken,
			},
			"articles/:page/:page_size": {
				Name:        "List articles",
				Handler:     handlers.GetArticles,
				Middlewares: []gin.HandlerFunc{middleware.ApiAuthCheck},
			},
			"comments": {
				Name:        "List comments",
				Handler:     handlers.GetComments,
				Middlewares: []gin.HandlerFunc{middleware.ApiAuthCheck},
			},
			"article/:article_id": {
				Name:        "Get specific article",
				Handler:     handlers.GetArticle,
				Middlewares: []gin.HandlerFunc{middleware.ApiAuthCheck},
			},
		},
		http.MethodPut: {
			"comment/:article_id": {
				Name:        "Post a comment",
				Handler:     handlers.PutComment,
				Middlewares: []gin.HandlerFunc{middleware.ApiAuthCheck},
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
}
