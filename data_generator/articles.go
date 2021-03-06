package data_generator

import (
	"math/rand"
	"on951/application"
	dbStructure "on951/database/structure"
	"on951/randomizer"
	"time"
)

func GenerateArticle() {
	db := application.GetApplication().GetDatabase().GetDB()

	rand.Seed(time.Now().UnixNano())

	var ArticleToInsert = dbStructure.Article{
		Title:       randomizer.GetRandomString(16, 1),
		Description: randomizer.GetRandomString(255, 8),
		AuthorId:    uint(rand.Intn(50)) + 1,
	}
	_ = db.Create(&ArticleToInsert)
}
