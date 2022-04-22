package randomizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRandomString(t *testing.T) {

	assert.Len(t, GetRandomString(10, 10), 100)
	assert.Len(t, GetRandomString(10, 1), 10)
	assert.Len(t, GetRandomString(5, 2), 10)
}
