package image_links_parser

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"main/state"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func getImageFileName(mimeType string, encoding string, encodedImage string) (string, bool) {
	isValid, extension := validateImage(mimeType, encoding, encodedImage)
	if isValid {
		imageFileNamePrefix := md5.Sum([]byte(mimeType + encoding + encodedImage))
		return hex.EncodeToString(imageFileNamePrefix[:]) + extension, true
	}
	return "", false
}

func grabAllValidImages(text string) map[string]string {
	imgTagRegexp, err := regexp.Compile(`<img src="data:([\w/]+);(.+),(.+)".*>`)
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
			if saveImage(filename, match[1], match[2], match[3]) {
				foundImages[filename] = match[0]
			}
		}
	}
	return foundImages
}

func validateImage(mimeType string, encoding string, encodedImage string) (bool, string) {
	parsedMimeType, _, err := mime.ParseMediaType(mimeType)
	if err != nil || parsedMimeType != mimeType {
		log.Printf("Unrecognized mime type %v\n", mimeType)
		return false, ""
	}
	if !strings.Contains(parsedMimeType, "image/") {
		log.Printf("Not an image mime type: %v\n", mimeType)
		return false, ""
	}
	var extensions []string
	extensions, err = mime.ExtensionsByType(mimeType)
	if len(extensions[0]) < 1 {
		log.Printf("No extension for %v\n", mimeType)
	}
	return true, extensions[0]
}

func saveImage(filename string, mimeType string, encoding string, encodedImage string) bool {

	var decodedImage []byte
	var err error
	switch encoding {
	case "base64":
		decodedImage, err = base64.StdEncoding.DecodeString(encodedImage)
		if err != nil {
			log.Printf("Cannot decode image %v\n", filename)
			return false
		}
		break
	case "text":
		decodedImage = []byte(encodedImage)
		break
	default:
		log.Printf("Unkown encoding %v\n", encoding)
		return false
	}
	err = ioutil.WriteFile(filepath.Join(state.GetImagesDir(), filename), decodedImage, fs.ModePerm)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func storeImageLinks(links []string) map[string]string {
	return map[string]string{}
}

func updateImageLinks(text string, links map[string]string) string {
	for filename, imgTag := range links {
		updatedImgTag := fmt.Sprintf("<img src=\"%v\" />", string(os.PathSeparator)+filepath.Join(state.ImagesDirectory, filename))
		text = strings.Replace(text, imgTag, updatedImgTag, -1)
	}
	return text
}

func Process(text string) (string, bool) {
	links := grabAllValidImages(text)
	return updateImageLinks(text, links), true

}
