package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

)

func (db *TDatabase)ConnectToDB(dsn string) {
	var err error
	if db.Db != nil {
		return
	}
	db.Db, err = gorm.Open(sqlite.Open(dsn), &db.Config)
	if err != nil {
		panic(err)
	}

}
func (db *TDatabase)GetDB() *gorm.DB {
	return db.Db
}

func (db *TDatabase)DisconnectDB() {
	db.Db = nil
}
