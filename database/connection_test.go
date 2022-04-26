package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestCase = struct {
	dsn                     string
	isConnectionEstablished bool
}

func TestConnectToDB(t *testing.T) {
	databaseMock := new(TDatabaseMock)

	for _, testCase := range getTestsDataForConnectToDB() {
		databaseMock.On("ConnectToDB", testCase.dsn).Return(testCase.isConnectionEstablished)
		if !testCase.isConnectionEstablished {
			databaseMock.On("ConnectToDB", testCase.dsn).Panic(DbConnectionError)
		}
	}
	for _, test := range getTestsDataForConnectToDB() {
		application := App{TDatabaseMock: databaseMock}
		if test.isConnectionEstablished {
			connOk := application.ConnectToDB(test.dsn)
			assert.Equal(t, "sqlite", application.GetDB().Dialector.Name())
			application.DisconnectDB()
			assert.Equal(t, test.isConnectionEstablished, connOk)
			assert.Nil(t, application.GetDB())
		} else {
			assert.PanicsWithError(t, DbConnectionError, func() { application.ConnectToDB(test.dsn) })
		}
	}
	databaseMock.AssertExpectations(t)

}
func getTestsDataForConnectToDB() []TestCase {
	return []TestCase{
		{dsn: "will-connect-to-valid-in-memory-db", isConnectionEstablished: true},
		{dsn: "invalid-dsn", isConnectionEstablished: false},
	}
}
