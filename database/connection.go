package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (db *TDatabase) ConnectToDB(dsn string) bool {
	var err error
	if db.Db != nil {
		return db.Db.Error == nil
	}
	db.Db, err = gorm.Open(sqlite.Open(dsn), &db.Config)
	if err != nil {
		panic(err)
	}
	return db.Db.Error == nil
}
func (db *TDatabase) GetDB() *gorm.DB {
	return db.Db
}

func (db *TDatabase) DisconnectDB() {
	db.Db = nil
}
