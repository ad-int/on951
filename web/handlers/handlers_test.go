package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http/httptest"
	"on951/api"
	"on951/application"
	"on951/database"
	dbStructure "on951/database/structure"
	"on951/image_links_parser"
	"os"
	"strconv"
	"testing"
)

type handlersTestSuite struct {
	suite.Suite
	context      *gin.Context
	db           database.TDatabaseMock
	recorder     *httptest.ResponseRecorder
	articlesRepo api.TArticlesRepository
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(handlersTestSuite))
}

func (suite *handlersTestSuite) prepare(handlerFunc func(ctx *gin.Context), testCase TestCase) {

	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)

	suite.db.On("ConnectToDB", "dummy").Return(true)
	suite.Assert().True(suite.db.ConnectToDB("dummy"), "Connecting to in-memory DB")
	suite.db.AutoMigrate()
	suite.articlesRepo = api.TArticlesRepository{IDatabase: &suite.db}

	for i := 1; i <= testCase.totalArticlesInDb; i++ {
		record := dbStructure.ArticleBriefInfo{Title: "article " + strconv.Itoa(i)}
		suite.db.GetDB().Table(dbStructure.TableArticles).Create(&record)
		for j := 1; j <= testCase.commentsPerArticle; j++ {
			commentText := "Comment " + strconv.Itoa(j)
			if testCase.commentsToInsert[j] != nil {
				commentText = testCase.commentsToInsert[j].Content
			}
			tempDir, err := ioutil.TempDir(os.TempDir(), "*")
			suite.Nil(err, "creating temp images dir")
			commentText, areImageLinksProcessed := image_links_parser.Process(commentText, tempDir, application.ImagesDirectory)
			comment := dbStructure.Comment{

				Content:   commentText,
				ArticleId: record.Id,
				UserId:    1,
			}

			suite.db.GetDB().Create(&comment)
			suite.True(areImageLinksProcessed)
		}
	}

	app := &application.TApplicationMock{}
	app.SetArticlesRepo(&suite.articlesRepo)
	application.SetApplication(app)
	suite.IsType(&application.TApplicationMock{}, application.GetApplication())
	suite.context.Params = testCase.params
	handlerFunc(suite.context)
	jsonBytes, err := json.MarshalIndent(testCase.response, "", "    ")
	suite.Nil(err)
	suite.Equal(string(jsonBytes)+"\r\n", suite.recorder.Body.String())
	application.GetApplication().GetArticlesRepo().DisconnectDB()
	suite.Nil(suite.articlesRepo.GetDB(), "disconnecting in-memory DB")

}

func (suite *handlersTestSuite) TestGetArticle() {
	for _, testCase := range testHandlersData.GetArticle {
		suite.prepare(GetArticle, testCase)
	}
}

func (suite *handlersTestSuite) TestGetArticles() {
	for _, testCase := range testHandlersData.GetArticles {
		suite.prepare(GetArticles, testCase)
	}
}
func (suite *handlersTestSuite) TestGetComments() {
	for _, testCase := range testHandlersData.GetComments {
		suite.prepare(GetComments, testCase)
	}
}
func (suite *handlersTestSuite) TestPutComment() {
	for _, testCase := range testHandlersData.PutComment {
		suite.prepare(PutComment, testCase)
	}
}
