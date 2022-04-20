package database

import (
	"gorm.io/gorm"
)

const DbConnectionError = "could not connect to DB"

type IDatabase interface {
	ConnectToDB(dsn string) bool
	DisconnectDB()
	GetDB() *gorm.DB
	AutoMigrate()
}

type TDatabase struct {
	Db     *gorm.DB
	Config gorm.Config
}
