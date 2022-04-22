package application

import (
	"github.com/gin-gonic/gin"
	"on951/router"
)

var testApplicationData = []struct {
	config             map[string]string
	imagesDir          string
	routes             map[string]router.TRoutesList
	configurationFails bool
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
		imagesDir:          "path/fake-images-dir",
		configurationFails: false,
		routes: map[string]router.TRoutesList{
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
		imagesDir:          "",
		routes: map[string]router.TRoutesList{
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
	{
		config: map[string]string{
			"DSN":                    "",
			"SECRET":                 "213243fdessf!!!!",
			"ISSUER":                 "localhost",
			"AUDIENCE":               "general",
			"TRUSTED_PROXIES":        "fail",
			"CORS_ALLOWED_HEADERS":   "*",
			"CORS_ALLOW_ALL_ORIGINS": "true",
		},
		configurationFails: true,
		imagesDir:          "image-dir",
		routes: map[string]router.TRoutesList{
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
