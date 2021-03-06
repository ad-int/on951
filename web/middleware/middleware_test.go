package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"on951/application"
	"on951/data_generator"
	"on951/database"
	"on951/models"
	"on951/web/handlers"
	"testing"
)

type middlewareTestSuite struct {
	suite.Suite
	context  *gin.Context
	db       database.TDatabaseMock
	recorder *httptest.ResponseRecorder
}

func (suite *middlewareTestSuite) SetupSuite() {

	suite.db.On("ConnectToDB", "dummy").Return(true)
	suite.Assert().True(suite.db.ConnectToDB("dummy"), "Connecting to in-memory DB")
	suite.db.AutoMigrate()

	app := &application.TApplicationMock{}
	app.SetDB(&suite.db)
	application.SetApplication(app)
	_ = data_generator.GenerateUser("guest", "not-set", 14)
}

func (suite *middlewareTestSuite) BeforeTest(suiteName, testName string) {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	body := ioutil.NopCloser(bytes.NewReader([]byte(`{"username":"guest","password":"not-set"}`)))
	suite.context.Request = httptest.NewRequest("POST", "//token", body)
	suite.db.On("ConnectToDB", "dummy").Return(true)
	suite.Assert().True(suite.db.ConnectToDB("dummy"), "Connecting to in-memory DB")
	suite.db.AutoMigrate()
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(middlewareTestSuite))
}

func (suite *middlewareTestSuite) TestApiAuthCheck() {

	app := &application.TApplicationMock{}
	app.SetDB(&suite.db)

	app.On("GetConfigValue", "AUDIENCE").Return("general")
	app.On("GetConfigValue", "ISSUER").Return("localhost")
	app.On("GetConfigValue", "SECRET").Return("234")
	app.On("GetConfigValue", "BCRYPT_HASH_GENERATION_COST").Return("14")

	application.SetApplication(app)
	handlers.GetToken(suite.context)

	authToken := models.AuthTokenResponse{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &authToken)

	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = httptest.NewRequest("", "//", nil)
	context.Request.Header = http.Header{}
	context.Request.Header.Set("Authorization", authToken.GetAuthorizationString())
	ApiAuthCheck(context)
	suite.NotEqual(context.Writer.Status(), http.StatusUnauthorized)
}

func (suite *middlewareTestSuite) TestApiAuthCheckThatFails() {

	app := &application.TApplicationMock{}

	app.On("GetConfigValue", "AUDIENCE").Return("general")
	app.On("GetConfigValue", "ISSUER").Return("localhost")
	app.On("GetConfigValue", "SECRET").Return("234")
	app.On("GetConfigValue", "BCRYPT_HASH_GENERATION_COST").Return("14")

	app.SetDB(&suite.db)

	application.SetApplication(app)
	handlers.GetToken(suite.context)

	authToken := models.AuthTokenResponse{}
	_ = json.Unmarshal(suite.recorder.Body.Bytes(), &authToken)

	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = httptest.NewRequest("", "//", nil)

	context.Request.Header = http.Header{}
	context.Request.Header.Set("Authorization", authToken.GetAuthorizationString()+"fail!!!")
	ApiAuthCheck(context)
	suite.Equal(context.Writer.Status(), http.StatusUnauthorized)
}
