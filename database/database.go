package database

import "gorm.io/gorm"

type IDatabase interface {
	
	connectToDB(dsn string)
	DisconnectDB()
	GetDB()
}

type TDatabase struct {
	Db *gorm.DB
	Config gorm.Config
}

