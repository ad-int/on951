package image_links_parser

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/fs"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func getImageFileName(mimeType string, encoding string, encodedImage string) (string, bool) {
	isValid, extension, _ := validateImage(mimeType, encoding, encodedImage)
	if isValid {
		imageFileNamePrefix := md5.Sum([]byte(mimeType + encoding + encodedImage))
		return hex.EncodeToString(imageFileNamePrefix[:]) + extension, true
	}
	return "", false
}

func grabAllValidImages(text string, imagesDir string) map[string]string {
	imgTagRegexp, err := regexp.Compile(`(?U)<img src="data:([\w/]+);([^"]+),([^"]+)".*>`)
	var foundImages = make(map[string]string)
	var index uint = 0
	if err != nil {
		log.Println(err)
	}
	for _, match := range imgTagRegexp.FindAllStringSubmatch(text, -1) {
		if len(match) != 4 {
			continue
		}
		index = index + 1

		filename, isValid := getImageFileName(match[1], match[2], match[3])
		if isValid {
			if saveImage(filepath.Join(imagesDir, filename), match[1], match[2], match[3]) {
				foundImages[filename] = match[0]
			}
		}

	}
	return foundImages
}

func validateImage(mimeType string, encoding string, encodedImage string) (bool, string, error) {
	parsedMimeType, _, err := mime.ParseMediaType(mimeType)
	if err != nil || parsedMimeType != mimeType {
		log.Printf("Unrecognized mime type %v\n", mimeType)
		return false, "", nil
	}
	if !strings.Contains(parsedMimeType, "image/") {
		log.Printf("Not an image mime type: %v\n", mimeType)
		return false, "", nil
	}
	var extensions []string
	extensions, err = mime.ExtensionsByType(mimeType)
	if err != nil {
		return false, "", errors.New(fmt.Sprintf("No extension for %v\n", mimeType))
	}
	if len(extensions) < 1 {
		return false, "", errors.New(fmt.Sprintf("No extension for %v\n", mimeType))
	}

	decodeImage := decodeContent(encoding, encodedImage)
	switch mimeType {
	case "image/png":
		_, err = png.Decode(bytes.NewReader(decodeImage))
		if err != nil {
			log.Println(err)
			return false, "", nil
		}
		break
	case "image/jpeg":
		_, err = jpeg.Decode(bytes.NewReader(decodeImage))
		if err != nil {
			return false, "", nil
		}
		break
	case "image/gif":
		_, err = gif.Decode(bytes.NewReader(decodeImage))
		if err != nil {
			return false, "", nil
		}
		break
	}
	return true, extensions[0], nil
}
func decodeContent(encoding string, encodedContent string) []byte {

	var decodedContent []byte
	var err error
	switch encoding {
	case "base64":
		decodedContent, err = base64.StdEncoding.DecodeString(encodedContent)
		if err != nil {
			log.Printf("Cannot decode image\n")
			return nil
		}
		break
	case "text":
		decodedContent = []byte(encodedContent)
		break
	default:
		log.Printf("Unkown encoding %v\n", encoding)
		return nil
	}

	return decodedContent
}
func saveImage(imagePath string, mimeType string, encoding string, encodedImage string) bool {

	decodedImage := decodeContent(encoding, encodedImage)
	if _, statErr := os.Stat(imagePath); os.IsExist(statErr) {
		return true
	}
	err := ioutil.WriteFile(imagePath, decodedImage, fs.ModePerm)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func updateImageLinks(text string, links map[string]string, urlPrefix string) string {
	for filename, imgTag := range links {
		updatedImgTag := fmt.Sprintf("<img src=\"%v\" />", string(os.PathSeparator)+filepath.Join(urlPrefix, filename))
		text = strings.Replace(text, imgTag, updatedImgTag, -1)
	}
	return text
}

func Process(text string, imagesDir string, urlPrefix string) (string, bool) {
	log.Println(imagesDir)
	if len(imagesDir) == 0 {
		return text, false
	}
	links := grabAllValidImages(text, imagesDir)
	return updateImageLinks(text, links, urlPrefix), true

}
