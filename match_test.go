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

func Test_Matches(t *testing.T) {
	s := Matches{&Match{Word: "bass", WordType: WordTypeFalsePositive},
		&Match{Word: "ass", WordType: WordTypeProfanity}}

	assert.Equal(t, true, s.HasProfaneMatch())
	assert.Equal(t, "ass", s.GetFirstProfaneMatch().Word)
	assert.Equal(t, 1, len(s.GetProfaneMatches()))

	assert.Equal(t, false, s.HasSuspectMatch())
	assert.Equal(t, 0, len(s.GetSuspectMatches()))

	assert.Equal(t, true, s.HasFalsePositiveMatch())
	assert.Equal(t, 1, len(s.GetFalsePositiveMatches()))
}
