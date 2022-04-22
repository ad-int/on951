package web

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"on951/models"
	"testing"
)

type writerTestSuite struct {
	suite.Suite
	context  *gin.Context
	recorder *httptest.ResponseRecorder
}

func (suite *writerTestSuite) BeforeTest(suiteName string, testName string) {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)

}

func TestWriterTestSuite(t *testing.T) {
	suite.Run(t, new(writerTestSuite))
}
func (suite *writerTestSuite) TestWriteNewLine() {
	writeNewLine(suite.context)
	suite.Equal("\r\n", suite.recorder.Body.String())
}

func (suite *writerTestSuite) TestWriteBadRequestError() {
	WriteBadRequestError(suite.context, "blah")
	jsonBytes, err := json.MarshalIndent(models.Response{Code: http.StatusBadRequest, Body: "blah"}, "", "    ")
	suite.Nil(err)
	suite.Equal(string(jsonBytes)+"\r\n", suite.recorder.Body.String())
}

func (suite *writerTestSuite) TestWriteMessage() {
	WriteMessage(suite.context, http.StatusOK, "OK")
	jsonBytes, err := json.MarshalIndent(models.Response{Code: http.StatusOK, Body: "OK"}, "", "    ")
	suite.Nil(err)
	suite.Equal(suite.recorder.Body.String(), string(jsonBytes)+"\r\n")
}
func (suite *writerTestSuite) TestWriteSuccessfullyCreatedMessage() {
	WriteMessage(suite.context, http.StatusCreated, "Created!")
	jsonBytes, err := json.MarshalIndent(models.Response{Code: http.StatusCreated, Body: "Created!"}, "", "    ")
	suite.Nil(err)
	suite.Equal(suite.recorder.Body.String(), string(jsonBytes)+"\r\n")
}
func (suite *writerTestSuite) TestWrite() {
	Write(suite.context, 200, models.Response{Code: 200, Body: "ok"})
	jsonBytes, err := json.MarshalIndent(models.Response{Code: 200, Body: "ok"}, "", "    ")
	suite.Nil(err)
	suite.Equal(suite.recorder.Body.String(), string(jsonBytes)+"\r\n")
}
