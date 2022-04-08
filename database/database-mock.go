package database

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type App struct {
	*TDatabaseMock
}

type TDatabaseMock struct {
	mock.Mock
	Db     *gorm.DB
	Config gorm.Config
}

func (db *TDatabaseMock) ConnectToDB(dsn string) bool {
	args := db.Called(dsn)

	if db.Db != nil {
		return db.Db.Error == nil
	}
	if dsn == "invalid-dsn" {
		panic(errors.New(DbConnectionError))
	}

	db.Db = &gorm.DB{}
	return args.Bool(0)

}
func (db *TDatabaseMock) GetDB() *gorm.DB {
	return db.Db
}

func (db *TDatabaseMock) DisconnectDB() {
	db.Db = nil
}
