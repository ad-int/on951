package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func ConnectToDB() *gorm.DB {
	if db != nil {
		return db
	}
	x, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database")
	}
	db = x
	return db
}
func GetDB() *gorm.DB {
	return db
}
