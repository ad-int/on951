package randomizer

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"strings"
	"time"
)

var characters = "abcdefghjklmnoprstuvwxyz0123456789`!@#$%^&*()_-+="

func GetRandomString(length int, size int) string {
	var charactersSlice []string

	rand.Seed(time.Now().UnixNano())

	charactersSlice = strings.Split(characters, "")

	var randomStr []string
	var element string
	for i := 0; i < length; i++ {
		for j := 0; j < size; j++ {
			randomIndex := rand.Intn(len(charactersSlice))
			element = strings.Join(charactersSlice[randomIndex:randomIndex+1], "")
			randomStr = append(randomStr, element)
		}

	}
	return strings.Join(randomStr, "")
}

func randomBase64Image(mimeType string, width int, height int) (string, string, error) {
	var err error
	rand.Seed(time.Now().UnixNano())
	r := image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(width, height)}
	img := image.NewNRGBA(r)
	buf := bytes.NewBuffer([]byte{})
	switch mimeType {
	case "image/png":
		err = png.Encode(buf, img)
		break
	case "image/jpeg":
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 100})
		break
	case "image/gif":
		err = gif.Encode(buf, img, &gif.Options{NumColors: 256})
		break
	default:
		err = errors.New("unknown image mime type")
		mimeType = ""
		break
	}

	return mimeType, base64.StdEncoding.EncodeToString(buf.Bytes()), err
}
