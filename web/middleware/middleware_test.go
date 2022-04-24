package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"on951/application"
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
	suite.context.Request = &http.Request{
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
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(middlewareTestSuite))
}

func (suite *middlewareTestSuite) TestApiAuthCheck() {

	app := &application.TApplicationMock{}

	app.On("GetConfigValue", "AUDIENCE").Return("general")
	app.On("GetConfigValue", "ISSUER").Return("localhost")
	app.On("GetConfigValue", "SECRET").Return("234")

	application.SetApplication(app)
	handlers.GetToken(suite.context)

	token := suite.recorder.Body.String()

	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme:      "",
			Opaque:      "",
			User:        nil,
			Host:        "",
			Path:        "",
			RawPath:     "",
			ForceQuery:  false,
			RawQuery:    "",
			Fragment:    "",
			RawFragment: "",
		},
	}
	context.Request.Header = http.Header{}
	context.Request.Header.Set("Authorization", "Bearer "+token)
	ApiAuthCheck(context)
	suite.NotEqual(context.Writer.Status(), http.StatusUnauthorized)
}
