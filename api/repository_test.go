package api

import (
	"github.com/stretchr/testify/suite"
	"on951/database"
	dbStructure "on951/database/structure"
	"strconv"
	"testing"
)

type repositoryTestSuite struct {
	suite.Suite
	db    database.TDatabaseMock
	aRepo TArticlesRepository
}

func (suite *repositoryTestSuite) SetupSuite() {

	suite.db.On("ConnectToDB", "dummy").Return(true)
	suite.Assert().True(suite.db.ConnectToDB("dummy"), "Connecting to in-memory DB")
	suite.db.AutoMigrate()
	suite.aRepo = TArticlesRepository{IDatabase: &suite.db}

}
func (suite *repositoryTestSuite) TearDownSuite() {

	suite.db.DisconnectDB()
	suite.Nil(suite.db.GetDB(), "Disconnected in-memory db")
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(repositoryTestSuite))
}

func (suite *repositoryTestSuite) TestGetArticles() {

	for _, testCase := range testGetArticlesData {

		for i := 1; i <= testCase.articlesInDb; i++ {
			record := dbStructure.ArticleBriefInfo{Title: "article " + strconv.Itoa(i)}
			suite.db.GetDB().Table(dbStructure.TableArticles).Create(&record)
		}

		articles := suite.aRepo.GetArticles(testCase.page, testCase.pageSize)
		suite.Equal(len(articles), testCase.articlesCount)
		suite.db.GetDB().Unscoped().AllowGlobalUpdate = true
		tx := suite.db.GetDB().Debug().Unscoped().Delete(&dbStructure.Article{})
		suite.Equal(int64(testCase.articlesInDb), tx.RowsAffected, "Deleted test articles")

	}
}
