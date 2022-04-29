package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"on951/api"
	"on951/application"
	"on951/database"
	dbStructure "on951/database/structure"
	"on951/image_links_parser"
	"on951/models"
	"os"
	"strconv"
	"testing"
)

type MalformedBody struct {
	io.ReadCloser
}

func (rc *MalformedBody) Write(p []byte) (n int, err error) {
	return 0, errors.New("Failed!")
}

func (rc *MalformedBody) Read(p []byte) (n int, err error) {
	return 0, errors.New("Failed!")
}

type handlersTestSuite struct {
	suite.Suite
	context      *gin.Context
	db           database.TDatabaseMock
	recorder     *httptest.ResponseRecorder
	articlesRepo api.ArticlesRepository
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(handlersTestSuite))
}

func (suite *handlersTestSuite) getNewAuthToken(cost int) string {
	app := &application.TApplicationMock{}

	suite.db.On("ConnectToDB", "dummy").Return(true)
	suite.Assert().True(suite.db.ConnectToDB("dummy"), "Connecting to in-memory DB")

	suite.db.AutoMigrate()
	app.SetDB(&suite.db)

	app.On("GetConfigValue", "AUDIENCE").Return("general")
	app.On("GetConfigValue", "ISSUER").Return("localhost")
	app.On("GetConfigValue", "SECRET").Return("234")
	app.On("GetConfigValue", "BCRYPT_HASH_GENERATION_COST").Return(strconv.Itoa(cost))

	oldApp := application.GetApplication()
	application.SetApplication(app)
	aRecorder := httptest.NewRecorder()
	aContext, _ := gin.CreateTestContext(aRecorder)
	body := ioutil.NopCloser(bytes.NewReader([]byte(`{"username":"guest'","password":"not-set"}`)))
	aContext.Request = httptest.NewRequest("GET", "//token", body)
	GetToken(aContext)
	application.SetApplication(oldApp)
	authToken := models.AuthTokenResponse{}
	_ = json.Unmarshal(aRecorder.Body.Bytes(), &authToken)
	return authToken.GetAuthorizationString()
}

func (suite *handlersTestSuite) prepare(testCase TestCase, requiresLogin bool, handlerFunc ...func(ctx *gin.Context)) {

	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)

	suite.db.On("ConnectToDB", "dummy").Return(true)
	suite.Assert().True(suite.db.ConnectToDB("dummy"), "Connecting to in-memory DB")
	if !testCase.runHandlerTwice {
		_ = suite.db.GetDB().Migrator().DropTable(dbStructure.TableArticles, "Comments")
	}
	suite.db.AutoMigrate()
	if testCase.runHandlerTwice {
		suite.articlesRepo = &api.TArticlesRepositoryMock{IDatabase: &suite.db}
	} else {
		suite.articlesRepo = &api.TArticlesRepository{IDatabase: &suite.db}
	}

	for i := 1; i <= testCase.totalArticlesInDb; i++ {
		record := dbStructure.ArticleBriefInfo{Title: "article " + strconv.Itoa(i)}
		suite.db.GetDB().Table(dbStructure.TableArticles).Create(&record)
		for j := 1; j <= testCase.commentsPerArticle; j++ {
			commentText := "Comment " + strconv.Itoa(j)
			userId := 1
			if testCase.commentsToInsert[j] != nil {
				commentText = testCase.commentsToInsert[j].Content
				userId = int(testCase.commentsToInsert[j].UserId)
			}
			tempDir, err := ioutil.TempDir(os.TempDir(), "*")
			suite.Nil(err, "creating temp images dir")
			commentText, areImageLinksProcessed := image_links_parser.Process(commentText, tempDir, application.ImagesDirectory)
			comment := dbStructure.Comment{

				Content:   commentText,
				ArticleId: record.Id,
				UserId:    uint(userId),
			}

			suite.db.GetDB().Create(&comment)
			suite.True(areImageLinksProcessed)
		}
	}

	app := &application.TApplicationMock{}
	app.SetDB(&suite.db)
	app.SetArticlesRepo(suite.articlesRepo)
	app.On("GetImagesDir").Return(testCase.imagesDir)

	if testCase.appConfig != nil {
		for key, value := range testCase.appConfig {
			app.On("GetConfigValue", key).Return(value)
		}
	}

	application.SetApplication(app)
	suite.IsType(&application.TApplicationMock{}, application.GetApplication())
	suite.context.Params = testCase.params

	var body io.ReadCloser
	if testCase.incorrectBody {
		body = &MalformedBody{}
	} else {
		body = ioutil.NopCloser(bytes.NewReader([]byte(testCase.body)))
	}
	suite.context.Request = httptest.NewRequest(testCase.method, testCase.requestURI, body)
	if requiresLogin {

		token := suite.getNewAuthToken(testCase.bcryptHashGenerationCost)

		app.On("GetConfigValue", "AUDIENCE").Return("general")
		app.On("GetConfigValue", "ISSUER").Return("localhost")
		app.On("GetConfigValue", "SECRET").Return("234")

		suite.context.Request.Header = http.Header{}
		suite.context.Request.Header.Set("Authorization", token)
	}
	for _, hFunc := range handlerFunc {
		hFunc(suite.context)
		if testCase.runHandlerTwice {
			suite.recorder.Body.Truncate(0)
			hFunc(suite.context)
		}
	}

	suite.Equal(testCase.statusCode, suite.context.Writer.Status())
	if testCase.response != nil {
		jsonBytes, err := json.MarshalIndent(testCase.response, "", "    ")
		suite.Nil(err)
		suite.Equal(string(jsonBytes)+"\r\n", suite.recorder.Body.String())
	}
	suite.db.DisconnectDB()
	suite.Nil(application.GetApplication().GetDatabase().GetDB(), "disconnecting in-memory DB")

}

func (suite *handlersTestSuite) TestGenerateArticles() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)

	suite.db.On("ConnectToDB", "dummy").Return(true)
	suite.Assert().True(suite.db.ConnectToDB("dummy"), "Connecting to in-memory DB")
	suite.db.AutoMigrate()
	suite.articlesRepo = &api.TArticlesRepository{IDatabase: &suite.db}
	app := &application.TApplicationMock{}
	app.SetDB(&suite.db)
	app.SetArticlesRepo(suite.articlesRepo)
	application.SetApplication(app)
	suite.IsType(&application.TApplicationMock{}, application.GetApplication())
	Generate()
	articles := application.GetApplication().GetArticlesRepo().GetArticles(1, 20)
	suite.Len(articles, 1)
	suite.db.DisconnectDB()
	suite.Nil(application.GetApplication().GetDatabase().GetDB(), "disconnecting in-memory DB")

}

func (suite *handlersTestSuite) TestGetArticle() {
	for _, testCase := range testHandlersData.GetArticle {
		suite.prepare(testCase, true, GetArticle)
	}
}

func (suite *handlersTestSuite) TestGetArticles() {
	for _, testCase := range testHandlersData.GetArticles {
		suite.prepare(testCase, true, GetArticles)
	}
}
func (suite *handlersTestSuite) TestGetComments() {
	for _, testCase := range testHandlersData.GetComments {
		suite.prepare(testCase, true, GetComments)
	}
}
func (suite *handlersTestSuite) TestHealthCheck() {
	for _, testCase := range testHandlersData.HealthCheck {
		suite.prepare(testCase, false, HealthCheck)
	}
}
func (suite *handlersTestSuite) TestGetToken() {
	for _, testCase := range testHandlersData.GetToken {
		suite.prepare(testCase, false, GetToken)
	}
}
func (suite *handlersTestSuite) TestPutComment() {
	for _, testCase := range testHandlersData.PutComment {
		suite.prepare(testCase, true, PutComment)
		if testCase.check != nil {
			suite.prepare(testCase.check.(TestCase), true, GetComments)
		}

	}
}
