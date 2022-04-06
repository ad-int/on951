package database

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestCase = struct{
	actual string
	expected string
	panics bool
}

func TestConnectToDB(t *testing.T) {
	assert := assert.New(t)
	for _, test := range getTestsDataForConnectToDB() {
		if test.panics {
			assert.Panics(func() { ConnectToDB(test.actual) } )
		} else {
			gormDB := ConnectToDB(test.actual)
			assert.Equal(test.expected, reflect.TypeOf(gormDB).String())
			Disconnect()
		}
	}

}
func getTestsDataForConnectToDB() []TestCase {
	return []TestCase{
		{"test.db", "*gorm.DB", false},
		{"mysql:////////dlfer'glfds", "", true},
	}
}