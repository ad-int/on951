package handlers

import (
	"log"
	"main/database"
	"main/models"
	"main/randomizer"
)

func Generate() {
	db := database.GetDB()

	var ArticleToInsert = models.Article{Title: randomizer.GetRandomString(16, 1), Description: randomizer.GetRandomString(255, 8)}
	// Create

	result := db.Create(&ArticleToInsert)

	if result.Error != nil {
		log.Println(result.Error)
	}
}
