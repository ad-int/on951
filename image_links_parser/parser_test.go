package image_links_parser

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"on951/application"
	"os"
	"strconv"
	"testing"
)

type TApplicationMock struct {
	mock.Mock
	application.TApplication
}

type imageLinksParserTestSuite struct {
	suite.Suite
}

func TestImageLinksParserTestSuite(t *testing.T) {
	suite.Run(t, new(imageLinksParserTestSuite))
}

func (suite *imageLinksParserTestSuite) TestGetImageFileName() {
	for index, testCase := range testImagesData {
		fileName, isValid := getImageFileName(testCase.mimyType, testCase.encoding, testCase.encodedImage)
		suite.Equal(testCase.fileName, fileName, "filename for test case "+strconv.Itoa(index))
		suite.Equal(testCase.valid, isValid)
	}

}

func (suite *imageLinksParserTestSuite) TestSaveImage() {
	for _, testCase := range testImagesData {
		isValid := false
		if testCase.valid {
			isValid = saveImage("test-"+testCase.fileName, testCase.mimyType, testCase.encoding, testCase.encodedImage)
		}
		suite.Equal(testCase.valid, isValid)
	}
}

func (suite *imageLinksParserTestSuite) TearDownSuite() {
	if dir := application.GetApplication().GetImagesDir(); len(dir) > 0 {
		_ = os.RemoveAll(dir)
	}
}
