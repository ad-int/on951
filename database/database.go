package database

import "gorm.io/gorm"

const DbConnectionError = "could not connect to DB"

type IDatabase interface {
	connectToDB(dsn string)
	DisconnectDB()
	GetDB()
}

type TDatabase struct {
	Db     *gorm.DB
	Config gorm.Config
}
