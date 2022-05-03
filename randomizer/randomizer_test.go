package randomizer

import (
	"encoding/base64"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestCaseType struct {
	mimeType      string
	width         int
	height        int
	triggersError bool
}
type randomizerTestSuite struct {
	suite.Suite
	testCases []TestCaseType
}

func (suite *randomizerTestSuite) SetupSuite() {
	suite.testCases = []TestCaseType{
		{
			mimeType:      "image/png",
			width:         0,
			height:        0,
			triggersError: true,
		},
		{
			mimeType:      "image/jpeg",
			width:         100,
			height:        100,
			triggersError: false,
		},
		{
			mimeType:      "image/jpg",
			width:         100,
			height:        100,
			triggersError: true,
		},
		{
			mimeType:      "image/gif",
			width:         55,
			height:        20,
			triggersError: false,
		},
	}
}

func TestRandomizerTestSuite(t *testing.T) {
	suite.Run(t, new(randomizerTestSuite))
}

func (suite *randomizerTestSuite) TestGetRandomString() {

	suite.Len(GetRandomString(10, 10), 100)
	suite.Len(GetRandomString(10, 1), 10)
	suite.Len(GetRandomString(5, 2), 10)
}

func (suite *randomizerTestSuite) TestGetBase64Image() {
	var mimeType, img string
	var err error
	for _, testCase := range suite.testCases {
		mimeType, img, err = GetRandomBase64Image(testCase.mimeType, testCase.width, testCase.height)
		suite.T().Log("generated image", mimeType, img)
		if testCase.triggersError {
			suite.Empty(img, "encoded image")
			suite.NotNil(err)
		} else {
			decodedImg, _ := base64.StdEncoding.DecodeString(img)
			suite.NotEmpty(mimeType, "mime type")
			suite.NotEmpty(decodedImg, "decoded image")
		}
	}
}
