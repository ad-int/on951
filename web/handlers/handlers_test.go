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
	"net/url"
	"on951/api"
	"on951/application"
	"on951/database"
	dbStructure "on951/database/structure"
	"on951/image_links_parser"
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
	articlesRepo api.TArticlesRepository
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(handlersTestSuite))
}

func (suite *handlersTestSuite) getNewAuthToken(cost int) string {
	app := &application.TApplicationMock{}

	app.On("GetConfigValue", "AUDIENCE").Return("general")
	app.On("GetConfigValue", "ISSUER").Return("localhost")
	app.On("GetConfigValue", "SECRET").Return("234")
	app.On("GetConfigValue", "BCRYPT_HASH_GENERATION_COST").Return(strconv.Itoa(cost))

	oldApp := application.GetApplication()
	application.SetApplication(app)
	aRecorder := httptest.NewRecorder()
	aContext, _ := gin.CreateTestContext(aRecorder)
	aContext.Request = &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme:      "",
			Opaque:      "",
			User:        nil,
			Host:        "",
			Path:        "",
			RawPath:     "",
			ForceQuery:  false,
			RawQuery:    "user=guest&password=p123",
			Fragment:    "",
			RawFragment: "",
		},
	}
	GetToken(aContext)
	application.SetApplication(oldApp)
	return aRecorder.Body.String()
}

func (suite *handlersTestSuite) prepare(testCase TestCase, requiresLogin bool, handlerFunc ...func(ctx *gin.Context)) {

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
	app.SetArticlesRepo(&suite.articlesRepo)
	app.On("GetImagesDir").Return("images")
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

		app.On("GetConfigValue", "AUDIENCE").Return("general")
		app.On("GetConfigValue", "ISSUER").Return("localhost")
		app.On("GetConfigValue", "SECRET").Return("234")

		token := suite.getNewAuthToken(testCase.bcryptHashGenerationCost)

		suite.context.Request.Header = http.Header{}
		suite.context.Request.Header.Set("Authorization", "Bearer "+token)
	}
	for _, hFunc := range handlerFunc {
		hFunc(suite.context)
	}

	jsonBytes, err := json.MarshalIndent(testCase.response, "", "    ")
	suite.Nil(err)
	suite.Equal(string(jsonBytes)+"\r\n", suite.recorder.Body.String())
	application.GetApplication().GetArticlesRepo().DisconnectDB()
	suite.Nil(suite.articlesRepo.GetDB(), "disconnecting in-memory DB")

}

func (suite *handlersTestSuite) TestGenerateArticles() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)

	suite.db.On("ConnectToDB", "dummy").Return(true)
	suite.Assert().True(suite.db.ConnectToDB("dummy"), "Connecting to in-memory DB")
	suite.db.AutoMigrate()
	suite.articlesRepo = api.TArticlesRepository{IDatabase: &suite.db}
	app := &application.TApplicationMock{}
	app.SetArticlesRepo(&suite.articlesRepo)
	application.SetApplication(app)
	suite.IsType(&application.TApplicationMock{}, application.GetApplication())
	Generate()
	articles := application.GetApplication().GetArticlesRepo().GetArticles(1, 20)
	suite.Len(articles, 1)
	application.GetApplication().GetArticlesRepo().DisconnectDB()
	suite.Nil(suite.articlesRepo.GetDB(), "disconnecting in-memory DB")

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
func (suite *handlersTestSuite) TestPutComment() {
	for _, testCase := range testHandlersData.PutComment {
		suite.prepare(testCase, true, PutComment)
		if testCase.check != nil {
			suite.prepare(testCase.check.(TestCase), true, GetComments)
		}

	}
}
