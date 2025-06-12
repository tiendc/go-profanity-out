package profanityout

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_skipHTMLTag(t *testing.T) {
	assert.Equal(t, 0, skipHTMLTag([]rune(""), 0))
	assert.Equal(t, 5, skipHTMLTag([]rune("<tag>"), 0))
	assert.Equal(t, 5, skipHTMLTag([]rune("<tag> x"), 0))
	assert.Equal(t, 0, skipHTMLTag([]rune("<tag"), 0))
}

func Test_decodeHTMLEntityAt(t *testing.T) {
	var ch rune
	var i int

	ch, i = decodeHTMLEntityAt([]rune("x&gt;"), 0)
	assert.True(t, ch == 0 && i == 0)

	ch, i = decodeHTMLEntityAt([]rune("x&gt;"), 1)
	assert.True(t, ch == '>' && i == 5)

	ch, i = decodeHTMLEntityAt([]rune("x&gt"), 1)
	assert.True(t, ch == 0 && i == 1)

	ch, i = decodeHTMLEntityAt([]rune("x&xxxxxx;"), 1)
	assert.True(t, ch == 0 && i == 1)

	ch, i = decodeHTMLEntityAt([]rune("x&#97;"), 1)
	assert.True(t, ch == 'a' && i == 6)

	ch, i = decodeHTMLEntityAt([]rune("x&#xxxx;"), 1)
	assert.True(t, ch == 0 && i == 1)
}
