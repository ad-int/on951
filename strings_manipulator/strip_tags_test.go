package strings_manipulator

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestStripTags(t *testing.T) {
	r := StripTags(
		`Hello <STYLE>.XSS{background-image:url("javascript:alert('XSS')");}</STYLE>my <A CLASS=XSS></A>World <img src="url"/>`,
	)
	assert.Equal(t, r, `Hello my World <img src="url"/>`)

	r = StripTags(
		`Hello <STYLE>.XSS{background-image:url("javascript:alert('XSS<img src=\"url\"/>')");}</STYLE>my <A CLASS=XSS></A>World`,
	)
	assert.Equal(t, r, `Hello my World`)

	r = StripTags(`My comment <m>jhh <img src="` + string(os.PathSeparator) + `images` + string(os.PathSeparator) + `51f967ab257ca1deb59556c321795702.gif" /><img src="data:application/json;base64,4554" /> test <img src="` + string(os.PathSeparator) + `images` + string(os.PathSeparator) + `e2bbf175e27a09a101eed75a59adcb24.png" />`)
	assert.Equal(t,
		`My comment jhh <img src="`+string(os.PathSeparator)+`images`+string(os.PathSeparator)+`51f967ab257ca1deb59556c321795702.gif"/><img src="data:application/json;base64,4554"/> test <img src="`+string(os.PathSeparator)+`images`+string(os.PathSeparator)+`e2bbf175e27a09a101eed75a59adcb24.png"/>`,
		r,
	)

	r = StripTags(
		`Hello <strong>test</strong><img alt="323" src="http://url/1"><img alt="nnn" onclick="javascript:alert('XSS')" src="http://url/2">`,
	)
	assert.Equal(t, r, `Hello test<img alt="323" src="http://url/1"><img alt="nnn" src="http://url/2">`)
}
