package database

import (
	"log"
	dbStructure "main/database/structure"
)

func AutoMigrate() {
	err := db.AutoMigrate(&dbStructure.User{}, &dbStructure.Article{}, &dbStructure.Comment{})
	if err != nil {
		log.Fatalf("Error occurred during DB migration %v\n", err)
	}
}
