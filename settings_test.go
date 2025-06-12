package profanityout

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Settings(t *testing.T) {
	s := &DetectorSettings{}

	WithSanitizeSpecialCharacters(true)(s)
	assert.Equal(t, true, s.SanitizeSpecialCharacters)

	WithSanitizeLeetSpeak(true)(s)
	assert.Equal(t, true, s.SanitizeLeetSpeak)

	WithSanitizeAccents(true)(s)
	assert.Equal(t, true, s.SanitizeAccents)

	WithSanitizeRepeatedCharacters(true)(s)
	assert.Equal(t, true, s.SanitizeRepeatedCharacters)

	WithSanitizeWildcardCharacters(true)(s)
	assert.Equal(t, true, s.SanitizeWildcardCharacters)

	WithProcessInputAsHTML(true)(s)
	assert.Equal(t, true, s.ProcessInputAsHTML)

	WithMatchWholeWord(true)(s)
	assert.Equal(t, true, s.MatchWholeWord)

	WithCensorCharacter('%')(s)
	assert.Equal(t, '%', s.CensorCharacter)
}
