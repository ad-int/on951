package router

import (
	"github.com/gin-gonic/gin"
)

var testRouterData = []struct {
	config             map[string]string
	configurationFails bool
	routes             map[string]TRoutesList
}{
	{
		config: map[string]string{
			"DSN":                    "",
			"SECRET":                 "213243fdessf!!!!",
			"ISSUER":                 "localhost",
			"AUDIENCE":               "general",
			"TRUSTED_PROXIES":        "127.0.0.1",
			"CORS_ALLOWED_HEADERS":   "*",
			"CORS_ALLOW_ALL_ORIGINS": "true",
		},
		configurationFails: false,
		routes: map[string]TRoutesList{
			"GET": {
				"health-check": {
					Name: "Health Check",
					Handler: func(context *gin.Context) {

					},
				},
			},
			"PUT": {
				"comment/:article_id": {
					Name: "Post a comment",
					Handler: func(context *gin.Context) {

					},
					Group: "/",
					Middlewares: []gin.HandlerFunc{
						func(context *gin.Context) {

						},
					},
				},
			},
		},
	},
	{
		config: map[string]string{
			"DSN":                    "",
			"SECRET":                 "213243fdessf!!!!",
			"ISSUER":                 "localhost",
			"AUDIENCE":               "general",
			"TRUSTED_PROXIES":        "faijop;l",
			"CORS_ALLOWED_HEADERS":   "*",
			"CORS_ALLOW_ALL_ORIGINS": "true",
		},
		configurationFails: true,
		routes: map[string]TRoutesList{
			"GET": {
				"health-check": {
					Name: "Health Check",
					Handler: func(context *gin.Context) {

					},
				},
			},
			"PUT": {
				"comment/:article_id": {
					Name: "Post a comment",
					Handler: func(context *gin.Context) {

					},
				},
			},
		},
	},
}
