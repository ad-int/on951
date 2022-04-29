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
	"on951/models"
	"on951/web/handlers"
	"testing"
)

type middlewareTestSuite struct {
	suite.Suite
	context  *gin.Context
	recorder *httptest.ResponseRecorder
}

func (suite *middlewareTestSuite) BeforeTest(suiteName, testName string) {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	body := ioutil.NopCloser(bytes.NewReader([]byte(`{"username":"guest'","password":"not-set"}`)))
	suite.context.Request = httptest.NewRequest("POST", "//token", body)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(middlewareTestSuite))
}

func (suite *middlewareTestSuite) TestApiAuthCheck() {

	app := &application.TApplicationMock{}

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

func (suite *middlewareTestSuite) TestApiAuthCheckWithHighCost() {

	app := &application.TApplicationMock{}

	app.On("GetConfigValue", "AUDIENCE").Return("general")
	app.On("GetConfigValue", "ISSUER").Return("localhost")
	app.On("GetConfigValue", "SECRET").Return("234")
	app.On("GetConfigValue", "BCRYPT_HASH_GENERATION_COST").Return("9999")

	application.SetApplication(app)
	handlers.GetToken(suite.context)

	suite.Equal(suite.context.Writer.Status(), http.StatusInternalServerError)
}
