package main

import (
	"log"
	"on951/application"
	"on951/web/handlers"
)

func main() {

	//ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	//ctx.Set("article_id1", 1)
	//handlers.GetComments(ctx)
	//image_links_parser.Process(`ikljghi  <img src="data:application/json;base64," /> 232refwdsf <img src="data:image/png;Es!a, QQQ" />grd`)
	//return

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
