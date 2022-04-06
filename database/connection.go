package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

)

var db *gorm.DB

func ConnectToDB(dsn string) *gorm.DB {
	if db != nil {
		return db
	}
	x, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = x
	return db

}
func GetDB() *gorm.DB {
	return db
}

func Disconnect() {
	db = nil
}
