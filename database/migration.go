package database

import (
	"log"
	dbStructure "main/database/structure"
)

func (db *TDatabase)AutoMigrate() {
	err := db.Db.AutoMigrate(&dbStructure.User{}, &dbStructure.Article{}, &dbStructure.Comment{})
	if err != nil {
		log.Fatalf("Error occurred during DB migration %v\n", err)
	}
}
