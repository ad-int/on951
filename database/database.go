package database

import (
	"gorm.io/gorm"
)

const DbConnectionError = "could not connect to DB"

type TDatabase struct {
	Db     *gorm.DB
	Config gorm.Config
}
