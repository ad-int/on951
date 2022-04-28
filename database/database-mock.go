package database

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
	"log"
	dbStructure "on951/database/structure"
)

type App struct {
	*TDatabaseMock
}

type DummyDialector struct {
	mock.Mock
	tests.DummyDialector
}

type DummyMigrator struct {
	gorm.Migrator
}

func (dd *DummyDialector) Initialize(db *gorm.DB) error {
	args := dd.Called(db)
	return args.Error(0)
}

func (dm *DummyMigrator) AutoMigrate(dst ...interface{}) error {
	return errors.New("migration failed")
}

func (dd *DummyDialector) Migrator(db *gorm.DB) gorm.Migrator {
	return &DummyMigrator{}
}

type TDatabaseMock struct {
	mock.Mock
	Db     *gorm.DB
	Config gorm.Config
}

func (db *TDatabaseMock) AutoMigrate() {
	err := db.Db.AutoMigrate(&dbStructure.User{}, &dbStructure.Article{}, &dbStructure.Comment{})
	if err != nil {
		log.Panicf("Error occurred during DB migration %v\n", err)
	}
}

func (db *TDatabaseMock) ConnectToDB(dsn string) bool {
	var err error
	args := db.Called(dsn)

	if db.Db != nil {
		return db.Db.Error == nil
	}

	switch dsn {
	case "invalid-dsn":
		panic(errors.New(DbConnectionError))
	case "dummy-db":
		dd := &DummyDialector{}
		dd.On("Initialize", mock.Anything).Return(nil)
		db.Db, err = gorm.Open(dd, &db.Config)
		break
	case "dummy-db-init-fails":
		dd := &DummyDialector{}
		dd.On("Initialize", mock.Anything).Return(errors.New(dsn))
		db.Db, err = gorm.Open(dd, &db.Config)
		break
	default:
		db.Db, err = gorm.Open(sqlite.Open("file:mem-db?mode=memory&cache=shared"), &db.Config)
		break
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
