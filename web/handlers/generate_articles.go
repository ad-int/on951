package handlers

import (
	"log"
	dbStructure "main/database/structure"
	"main/randomizer"
	"main/state"
	"math/rand"
	"time"
)

func Generate() {
	db := state.GetApplication().GetDB()

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
