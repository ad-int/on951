package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

)

func (db *TDatabase)ConnectToDB(dsn string) {
	if db.Db != nil {
		return
	}
	x, err := gorm.Open(sqlite.Open(dsn), &db.Config)
	if err != nil {
		panic(err)
	}
	db.Db = x

}
func (db *TDatabase)GetDB() *gorm.DB {
	return db.Db
}

func (db *TDatabase)Disconnect() {
	db.Db = nil
}
