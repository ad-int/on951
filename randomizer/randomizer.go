package randomizer

import (
	"math/rand"
	"strings"
	"time"
)

var characters = "abcdefghjklmnoprstuvwxyz0123456789`!@#$%^&*()_-+="

func GetRandomString(length int, size int) string  {
var charactersSlice []string

	rand.Seed(time.Now().UnixNano())

	charactersSlice = strings.Split(characters,"")

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
