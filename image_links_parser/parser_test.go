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

func (suite *imageLinksParserTestSuite) TearDownSuite() {
	if dir := application.GetApplication().GetImagesDir(); len(dir) > 0 {
		_ = os.RemoveAll(dir)
	}
}

func (suite *imageLinksParserTestSuite) TestGetImageFileName() {
	for index, testCase := range testImagesData {
		fileName, isValid := getImageFileName(testCase.mimyType, testCase.encoding, testCase.encodedImage)
		suite.Equal(testCase.fileName+testCase.extension, fileName, "filename for test case "+strconv.Itoa(index))
		suite.Equal(testCase.valid, isValid)
	}

}

func (suite *imageLinksParserTestSuite) TestSaveImage() {
	for _, testCase := range testImagesData {
		isValid := false
		if testCase.valid {
			isValid = saveImage("test-"+testCase.fileName+testCase.extension, testCase.mimyType, testCase.encoding, testCase.encodedImage)
		}
		suite.Equal(testCase.valid, isValid)
	}
}

func (suite *imageLinksParserTestSuite) TestValidateImage() {
	for _, testCase := range testImagesData {
		isValid := false
		extension := ""
		if testCase.valid {
			isValid, extension, _ = validateImage(testCase.mimyType, testCase.encoding, testCase.encodedImage)
		}
		suite.Equal(testCase.valid, isValid)
		suite.Equal(testCase.extension, extension)
	}
}

func (suite *imageLinksParserTestSuite) TestGetAllValidImages() {
	for _, testCase := range testCommentsData {
		suite.Equal(testCase.images, grabAllValidImages(testCase.text))
	}
}
func (suite *imageLinksParserTestSuite) TestUpdateImageLinks() {
	for _, testCase := range testCommentsWithFixedLinksData {
		suite.Equal(testCase.textAfter, updateImageLinks(testCase.textBefore, testCase.images))
	}
}
