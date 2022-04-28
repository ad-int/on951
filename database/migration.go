package database

import (
	"log"
	dbStructure "on951/database/structure"
)

func (db *TDatabase) AutoMigrate() {
	err := db.Db.AutoMigrate(&dbStructure.User{}, &dbStructure.Article{}, &dbStructure.Comment{})
	if err != nil {
		log.Panicf("Error occurred during DB migration %v\n", err)
	}
}
