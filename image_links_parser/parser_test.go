package image_links_parser

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"on951/application"
	"on951/randomizer"
	"os"
	"path/filepath"
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
		tempDir, err := ioutil.TempDir(os.TempDir(), "*")
		suite.Nil(err, "creating temp images dir")

		isValid = saveImage(filepath.Join(tempDir, testCase.fileName+testCase.extension), testCase.mimyType, testCase.encoding, testCase.encodedImage)
		isValid = saveImage(filepath.Join(tempDir, testCase.fileName+testCase.extension), testCase.mimyType, testCase.encoding, testCase.encodedImage)

		suite.Equal(testCase.valid, isValid)
	}
}
func (suite *imageLinksParserTestSuite) TestSaveImage2() {
	for _, testCase := range testImagesData {
		isValid := false
		tempDir := "#############"

		if testCase.valid {
			isValid = saveImage(filepath.Join(tempDir, testCase.fileName+testCase.extension), testCase.mimyType, testCase.encoding, testCase.encodedImage)
		}
		suite.False(isValid)
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
		tempDir, err := ioutil.TempDir(os.TempDir(), "*")
		suite.Nil(err, "creating temp images dir")
		suite.Equal(testCase.images, grabAllValidImages(testCase.text, tempDir))
	}
}
func (suite *imageLinksParserTestSuite) TestUpdateImageLinks() {
	for _, testCase := range testCommentsWithFixedLinksData {
		suite.Equal(testCase.textAfter, updateImageLinks(testCase.textBefore, testCase.images, "images"))
	}
}

func (suite *imageLinksParserTestSuite) TestDecodeContent() {
	for _, testCase := range testDecodeContentData {
		decodedContent, _ := decodeContent(testCase.encoding, testCase.content)
		suite.Equal(testCase.decodedContent, decodedContent)
	}
}
func (suite *imageLinksParserTestSuite) TestProcess() {

	for _, testCase := range testProcessData {
		d := testCase.imagesDir
		if len(testCase.imagesDir) > 0 {
			d, _ = ioutil.TempDir(os.TempDir(), testCase.imagesDir)
		}
		parsedText, success := Process(testCase.content, d, testCase.urlPrefix)
		suite.Equal(testCase.success, success)
		suite.Equal(testCase.parsedContent, parsedText)

	}

}
func (suite *imageLinksParserTestSuite) TestHeavyProcess() {

	d, _ := ioutil.TempDir(os.TempDir(), "images-*")
	urlPrefix := "images"
	var parsedContent, content, randomText, img string
	for i := 0; i < 1024; i++ {
		for _, mimeType := range []string{"image/gif", "image/png", "image/jpeg"} {
			mimeType, img, _ = randomizer.GetRandomBase64Image("image/png", 10, 10)
			randomText = randomizer.GetRandomString(16, 2)
			content += randomText + ` <img src="data:` + mimeType + `;base64,` + img + `"/>`
			fName, _ := getImageFileName(mimeType, "base64", img)
			parsedContent += randomText +
				` <img src="` + string(os.PathSeparator) + `images` + string(os.PathSeparator) + fName + `" />`
		}
	}
	parsedText, success := Process(content, d, urlPrefix)
	suite.True(success)
	suite.Equal(parsedContent, parsedText)
}
