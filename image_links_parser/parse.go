package image_links_parser

import (
	"log"
	"mime"
	"regexp"
)

func findAllLinks(text string) []string {
	imgTagRegexp, err := regexp.Compile(`<img src="data:([\w/]+);(.+?),(.+?)"?>`)
	if err != nil {
		log.Println(err)
	}
	for _, match := range imgTagRegexp.FindAllStringSubmatch(text, -1) {
		if len(match) != 4 {
			continue
		}
		_ = validateImage(match[1], match[2], match[3])
	}
	return []string{}
}

func validateImage(mimeType string, encoding string, encodedImage string) bool {
	parsedMimeType, _, err := mime.ParseMediaType(mimeType)
	if err != nil || parsedMimeType != mimeType {
		log.Printf("Unrecognized mime type %v\n", mimeType)
		return false
	}
	return false
}

func storeImageLinks(links []string) map[string]string {
	return map[string]string{}
}

func updateImageLinks(text string, links map[string]string) string {
	return text
}

func Process(text string) (string, bool) {
	links := findAllLinks(text)
	storedLinks := storeImageLinks(links)
	return updateImageLinks(text, storedLinks), true

}
