package database

import (
	"testing"
)

type TestCase = struct {
	actual   string
	expected string
	panics   bool
}

func TestConnectToDB(t *testing.T) {
	//assert := assert.New(t)
	databaseMock := new(TDatabaseMock)

	for _, test := range getTestsDataForConnectToDB() {
		if test.panics {
			databaseMock.On("ConnectToDB", test.actual).Panic("Cannot connect to DB")
			//assert.Panics(func() { database.ConnectToDB(test.actual) })
		} else {
			databaseMock.On("ConnectToDB", test.actual).Return()
			//database.ConnectToDB(test.actual)
			//assert.Equal(test.expected, reflect.TypeOf(database.GetDB()).String())
			//database.DisconnectDB()
		}
		a := App{*databaseMock}
		a.ConnectToDB(test.actual)
		databaseMock.AssertExpectations(t)
	}

}
func getTestsDataForConnectToDB() []TestCase {
	return []TestCase{
		{"test.db", "*gorm.DB", false},
		{"invalid-dsn", "", true},
	}
}
