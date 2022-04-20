package api

import (
	assert "github.com/stretchr/testify/assert"
	"on951/database"
	dbStructure "on951/database/structure"
	"strconv"
	"testing"
)

func TestGetArticles(t *testing.T) {
	mAssert := assert.New(t)

	for _, testCase := range testGetArticlesData {

		m := new(database.TDatabaseMock)

		m.On("ConnectToDB", "dummy").Return(true)
		assert.True(t, m.ConnectToDB("dummy"), "Connecting to in-memory DB")
		m.AutoMigrate()

		for i := 1; i <= testCase.articlesInDb; i++ {
			record := dbStructure.ArticleBriefInfo{Title: "article " + strconv.Itoa(i)}
			m.GetDB().Table(dbStructure.TableArticles).Create(&record)
		}

		aRepo := TArticlesRepository{
			IDatabase: m,
		}

		articles := aRepo.GetArticles(testCase.page, testCase.pageSize)
		mAssert.Equal(len(articles), testCase.articlesCount)
		m.DisconnectDB()
	}
}
