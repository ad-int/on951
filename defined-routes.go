package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"on951/router"
	"on951/web/handlers"
	"on951/web/middleware"
)

var DefinedRoutes = map[string]router.TRoutesList{
	http.MethodPost: {
		"token": {
			Name:    "Get token",
			Handler: handlers.GetToken,
		},
	},
	http.MethodGet: {
		"health-check": {
			Name:    "Health Check",
			Handler: handlers.HealthCheck,
		},
		"articles": {
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
}
