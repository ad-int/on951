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
		a := App{databaseMock}
		if test.isConnectionEstablished {
			connOk := a.ConnectToDB(test.dsn)
			a.DisconnectDB()
			assert.Equal(t, test.isConnectionEstablished, connOk)
		} else {
			assert.PanicsWithError(t, DbConnectionError, func() { a.ConnectToDB(test.dsn) })
		}
	}
	databaseMock.AssertExpectations(t)

}
func getTestsDataForConnectToDB() []TestCase {
	return []TestCase{
		{dsn: "test.db", isConnectionEstablished: true},
		{dsn: "invalid-dsn", isConnectionEstablished: false},
	}
}
