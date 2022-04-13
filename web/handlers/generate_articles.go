package handlers

import (
	"log"
	"math/rand"
	dbStructure "on951/database/structure"
	"on951/randomizer"
	"on951/state"
	"time"
)

func Generate() {
	db := state.GetApplication().GetArticlesRepo().GetDB()

	rand.Seed(time.Now().UnixNano())

	var ArticleToInsert = dbStructure.Article{
		Title:       randomizer.GetRandomString(16, 1),
		Description: randomizer.GetRandomString(255, 8),
		AuthorId:    uint(rand.Intn(50)) + 1,
	}
	// Create

	result := db.Create(&ArticleToInsert)

	if result.Error != nil {
		log.Println(result.Error)
	}
}
