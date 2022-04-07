package database

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type App struct {
	TDatabaseMock
}

type TDatabaseMock struct {
	mock.Mock
	Db *gorm.DB
	Config gorm.Config
}
func (db *TDatabaseMock)ConnectToDB(dsn string) {
	if db.Db != nil {
		return
	}
	args := db.Called(dsn)

	if dsn == "invalid-dsn" {
		panic(errors.New("Could not connect to DB"))
	}
	args.Bool(0)

	db.Db = &gorm.DB{}

}
func (db *TDatabaseMock)GetDB() *gorm.DB {
	return db.Db
}

func (db *TDatabaseMock)DisconnectDB() {
	db.Db = nil
}
