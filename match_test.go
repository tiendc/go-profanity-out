package profanityout

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Match(t *testing.T) {
	s := &Match{}

	s.WordType = WordTypeSuspect
	assert.Equal(t, WordTypeSuspect, s.WordType)
	assert.Equal(t, true, s.IsSuspect())

	s.WordType = WordTypeProfanity
	assert.Equal(t, WordTypeProfanity, s.WordType)
	assert.Equal(t, true, s.IsProfane())
	assert.Equal(t, false, s.IsSuspect())
	assert.Equal(t, false, s.IsFalsePositive())

	s.WordType = WordTypeFalsePositive
	assert.Equal(t, WordTypeFalsePositive, s.WordType)
	assert.Equal(t, true, s.IsFalsePositive())
	assert.Equal(t, false, s.IsSuspect())
}
