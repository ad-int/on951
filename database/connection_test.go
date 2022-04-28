package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestCase = struct {
	dsn                     string
	isConnectionEstablished bool
	driverName              string
	migrationPanics         bool
	panicMsg                string
}

func TestConnectToDB(t *testing.T) {
	for _, testCase := range getTestsDataForConnectToDB() {
		databaseMock := new(TDatabaseMock)

		databaseMock.On("ConnectToDB", testCase.dsn).Return(testCase.isConnectionEstablished)
		if !testCase.isConnectionEstablished || testCase.panicMsg != "" {
			databaseMock.On("ConnectToDB", testCase.dsn).Panic(testCase.panicMsg)
		}
		application := App{TDatabaseMock: databaseMock}
		if testCase.isConnectionEstablished {
			connOk := application.ConnectToDB(testCase.dsn)
			assert.Equal(t, testCase.driverName, application.GetDB().Dialector.Name())
			if testCase.migrationPanics {
				assert.Panics(t, func() {
					application.AutoMigrate()
				})
			} else {
				application.AutoMigrate()
			}
			application.DisconnectDB()
			assert.Equal(t, testCase.isConnectionEstablished, connOk)
			assert.Nil(t, application.GetDB())
		} else if testCase.panicMsg != "" {
			assert.PanicsWithError(t, testCase.panicMsg, func() { application.ConnectToDB(testCase.dsn) })
		}

		databaseMock.AssertExpectations(t)
	}

}
func getTestsDataForConnectToDB() []TestCase {
	return []TestCase{
		{dsn: "will-connect-to-valid-in-memory-db", isConnectionEstablished: true, driverName: "sqlite"},
		{dsn: "dummy-db", isConnectionEstablished: true, driverName: "dummy", migrationPanics: true},
		{dsn: "dummy-db-init-fails", isConnectionEstablished: false, driverName: "dummy", panicMsg: "dummy-db-init-fails"},
		{dsn: "invalid-dsn", isConnectionEstablished: false, panicMsg: DbConnectionError},
	}
}
func TestConnectToMemoryDB(t *testing.T) {
	db := &TDatabase{}
	db.ConnectToDB("file:test-mem-db?mode=memory")
	db.ConnectToDB("file:test-mem-db?mode=memory")
	db.AutoMigrate()
	assert.NotNil(t, db.GetDB())
	db.DisconnectDB()
	assert.Nil(t, db.GetDB())
}

func TestConnectToMemoryDBFailsBecauseOfInvalidMode(t *testing.T) {
	db := &TDatabase{}
	assert.Panics(t, func() {
		db.ConnectToDB("file:test-mem-db?mode=sadaszxfddfreqr43eqs")
	})

}
