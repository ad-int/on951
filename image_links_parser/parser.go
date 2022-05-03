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
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var mx sync.Mutex
var wg sync.WaitGroup

func getImageFileName(mimeType string, encoding string, encodedImage string) (string, bool) {
	isValid, extension, _ := validateImage(mimeType, encoding, encodedImage)
	if isValid {
		imageFileNamePrefix := md5.Sum([]byte(mimeType + encoding + encodedImage))
		return hex.EncodeToString(imageFileNamePrefix[:]) + extension, true
	}
	return "", false
}

func grabAllValidImages(text string, imagesDir string) map[string]string {
	imgTagRegexp, _ := regexp.Compile(`(?U)<img src="data:([\w/]+);([^"]+),([^"]+)".*>`)
	var foundImages = make(map[string]string)
	allMatches := imgTagRegexp.FindAllStringSubmatch(text, -1)
	wg.Add(len(allMatches))
	for _, match := range allMatches {
		go func(match []string) {
			filename, isValid := getImageFileName(match[1], match[2], match[3])
			if isValid {
				if saved, _ := saveImage(filepath.Join(imagesDir, filename), match[1], match[2], match[3]); saved {
					mx.Lock()
					foundImages[filename] = match[0]
					mx.Unlock()
				}
			}
			wg.Done()
		}(match)
	}
	wg.Wait()
	return foundImages
}

func validateImage(mimeType string, encoding string, encodedImage string) (bool, string, error) {
	parsedMimeType, _, _ := mime.ParseMediaType(mimeType)

	if !strings.Contains(parsedMimeType, "image/") {
		return false, "", errors.New(fmt.Sprintf("Not an image mime type: %v\n", mimeType))
	}
	var err error
	var extensions []string
	extensions, _ = mime.ExtensionsByType(mimeType)

	if len(extensions) < 1 {
		return false, "", errors.New(fmt.Sprintf("No extension for %v\n", mimeType))
	}

	decodeImage, _ := decodeContent(encoding, encodedImage)
	switch mimeType {
	case "image/png":
		_, err = png.Decode(bytes.NewReader(decodeImage))
		if err != nil {
			return false, "", err
		}
		break
	case "image/jpeg":
		_, err = jpeg.Decode(bytes.NewReader(decodeImage))
		extensions[0] = ".jpg" // overwrite strange extensions such as '.jfif' or '.jpe'
		if err != nil {
			return false, "", err
		}
		break
	case "image/gif":
		_, err = gif.Decode(bytes.NewReader(decodeImage))
		if err != nil {
			return false, "", err
		}
		break
	}
	return true, extensions[0], nil
}
func decodeContent(encoding string, encodedContent string) ([]byte, error) {

	var decodedContent []byte
	var err error
	switch encoding {
	case "base64":
		decodedContent, err = base64.StdEncoding.DecodeString(encodedContent)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Cannot decode image\n"))
		}
		break
	case "text":
		decodedContent = []byte(encodedContent)
		break
	default:
		return nil, errors.New(fmt.Sprintf("Unkown encoding %v\n", encoding))
	}

	return decodedContent, nil
}
func saveImage(imagePath string, mimeType string, encoding string, encodedImage string) (bool, error) {
	var err error
	decodedImage, err := decodeContent(encoding, encodedImage)
	if err != nil {
		return false, err
	}
	if _, statErr := os.Stat(imagePath); statErr == nil {
		return true, errors.New(fmt.Sprintln(imagePath, "already exists"))
	}
	err = ioutil.WriteFile(imagePath, decodedImage, fs.ModePerm)
	if err != nil {
		return false, err
	}
	return true, nil
}

func updateImageLinks(text string, links map[string]string, urlPrefix string) string {
	for filename, imgTag := range links {
		updatedImgTag := fmt.Sprintf("<img src=\"%v\" />", string(os.PathSeparator)+filepath.Join(urlPrefix, filename))
		text = strings.Replace(text, imgTag, updatedImgTag, -1)
	}
	return text
}

func Process(text string, imagesDir string, urlPrefix string) (string, bool) {
	if len(imagesDir) == 0 {
		return text, false
	}
	links := grabAllValidImages(text, imagesDir)
	return updateImageLinks(text, links, urlPrefix), true

}
