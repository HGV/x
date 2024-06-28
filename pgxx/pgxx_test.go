package pgxx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLikeContains(t *testing.T) {
	assert.Equal(t, LikeBegins("abc"), "abc%")
	assert.Equal(t, LikeEnds("abc"), "%abc")
	assert.Equal(t, LikeContains("abc"), "%abc%")
}
