package api

import (
	"github.com/stretchr/testify/suite"
	"on951/database"
	dbStructure "on951/database/structure"
	"strconv"
	"testing"
)

const defaultNumberOfArticleInTestDb = 20

type repositoryTestSuite struct {
	suite.Suite
	db    database.TDatabaseMock
	aRepo TArticlesRepository
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(repositoryTestSuite))
}

func (suite *repositoryTestSuite) BeforeTest(suiteName, testName string) {
	suite.db.On("ConnectToDB", "dummy").Return(true)
	suite.Assert().True(suite.db.ConnectToDB("dummy"), "Connecting to in-memory DB")
	suite.db.AutoMigrate()
	suite.db.GetDB().Unscoped().AllowGlobalUpdate = true
	suite.aRepo = TArticlesRepository{IDatabase: &suite.db}

	if testName != "TestGetArticles" {
		suite.createTestArticles(defaultNumberOfArticleInTestDb)
	}
}

func (suite *repositoryTestSuite) AfterTest(suiteName, testName string) {
	suite.db.DisconnectDB()
	suite.Nil(suite.db.GetDB(), "Disconnected in-memory db")
}

func (suite *repositoryTestSuite) createTestArticles(total int) {
	for i := 1; i <= total; i++ {
		record := dbStructure.ArticleBriefInfo{Title: "article " + strconv.Itoa(i)}
		suite.db.GetDB().Table(dbStructure.TableArticles).Create(&record)
	}
}

func (suite *repositoryTestSuite) TestGetArticles() {

	for _, testCase := range testGetArticlesData {

		suite.createTestArticles(testCase.articlesInDb)

		articles := suite.aRepo.GetArticles(testCase.page, testCase.pageSize)
		suite.Equal(len(articles), testCase.articlesCount)
		tx := suite.db.GetDB().Debug().Unscoped().Delete(&dbStructure.Article{})
		suite.Equal(int64(testCase.articlesInDb), tx.RowsAffected, "Deleted test articles")

	}
}

func (suite *repositoryTestSuite) TestGetArticle() {
	var testArticle = dbStructure.Article{
		Id:          999,
		Title:       "Test title",
		Description: "dfsfdfgcxfd",
		AuthorId:    1455,
	}

	suite.db.GetDB().Create(&testArticle)
	article, found := suite.aRepo.GetArticle(int(testArticle.Id))
	suite.True(found, "Article found")
	suite.Equal(testArticle, article, "Article data match")
}
func (suite *repositoryTestSuite) TestPutComment() {
	suite.TestGetArticle()
	var testComment = dbStructure.Comment{
		Id:        5,
		ArticleId: 999,
		Content:   "abc",
		UserId:    1455,
	}

	var testArticle = dbStructure.ArticleWithComments{
		Id:          999,
		AuthorId:    1455,
		Title:       "Test title",
		Description: "dfsfdfgcxfd",
		Comments: []dbStructure.Comment{
			testComment,
		},
	}

	suite.db.GetDB().Create(&testComment)
	article, found := suite.aRepo.GetArticleWithComments(999, 1, 1)
	suite.True(found, "Article found")
	suite.Equal(testArticle, article, "Article data match")
}
