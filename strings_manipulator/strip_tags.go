package strings_manipulator

import (
	"github.com/microcosm-cc/bluemonday"
)

func StripTags(text string) string {
	p := bluemonday.StripTagsPolicy()
	p.AllowAttrs("src", "alt").OnElements("img")
	return p.Sanitize(text)
}
