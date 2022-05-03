package database

import "gorm.io/gorm"

type IDatabase interface {
	ConnectToDB(dsn string) bool
	DisconnectDB()
	GetDB() *gorm.DB
	AutoMigrate()
}
