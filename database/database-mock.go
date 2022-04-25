package database

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	dbStructure "on951/database/structure"
)

type App struct {
	*TDatabaseMock
}

type TDatabaseMock struct {
	mock.Mock
	Db     *gorm.DB
	Config gorm.Config
}

func (db *TDatabaseMock) AutoMigrate() {
	err := db.Db.AutoMigrate(&dbStructure.User{}, &dbStructure.Article{}, &dbStructure.Comment{})
	if err != nil {
		log.Fatalf("Error occurred during DB migration %v\n", err)
	}
}

func (db *TDatabaseMock) ConnectToDB(dsn string) bool {
	var err error
	args := db.Called(dsn)

	if db.Db != nil {
		return db.Db.Error == nil
	}

	if dsn == "connection-fails" {
		return false
	}
	if dsn == "invalid-dsn" {
		panic(errors.New(DbConnectionError))
	}

	if dsn == "invalid-driver" {
		db.Db, err = gorm.Open(DummyDialector{}, &db.Config)
	} else {
		db.Db, err = gorm.Open(sqlite.Open(""), &db.Config)
	}

	if err != nil {
		panic(err)
	}
	return args.Bool(0)

}
func (db *TDatabaseMock) GetDB() *gorm.DB {
	return db.Db
}

func (db *TDatabaseMock) DisconnectDB() {
	db.Db = nil
}
