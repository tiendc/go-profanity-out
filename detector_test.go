package profanityout

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tiendc/go-profanity-out/data/en"
)

const (
	sanitizeAccents            = false
	sanitizeLeetSpeak          = true
	sanitizeSpecialCharacters  = true
	sanitizeSpaces             = true
	sanitizeRepeatedCharacters = true
	sanitizeWildcardCharacters = false
	processInputAsHTML         = false
	matchWholeWord             = true
)

func newDetectorEN() *ProfanityDetector {
	return NewProfanityDetector().
		WithProfaneWords(en.DefaultProfanities).
		WithFalsePositiveWords(en.DefaultFalsePositives).
		WithSuspectWords(en.DefaultSuspects).
		WithLeetSpeakCharacters(en.LeetSpeakCharacters).
		WithSpecialCharacters(en.SpecialCharacters).
		WithWildcardCharacters(en.WildcardCharacters).
		WithSanitizeLeetSpeak(sanitizeLeetSpeak).
		WithSanitizeSpecialCharacters(sanitizeSpecialCharacters).
		WithSanitizeSpaces(sanitizeSpaces).
		WithSanitizeAccents(sanitizeAccents).
		WithSanitizeRepeatedCharacters(sanitizeRepeatedCharacters).
		WithSanitizeWildcardCharacters(sanitizeWildcardCharacters).
		WithProcessInputAsHTML(processInputAsHTML).
		WithMatchWholeWord(matchWholeWord)
}

func toCmp(m *Match) *Match {
	return &Match{
		Word:          m.Word,
		WordType:      m.WordType,
		Start:         m.Start,
		End:           m.End,
		LeadingSpace:  m.LeadingSpace,
		TrailingSpace: m.TrailingSpace,
		Text:          m.Text,
	}
}

func Test_Scan_One(t *testing.T) {
	d := newDetectorEN
	var m []*Match

	t.Run("Match whole word tests", func(t *testing.T) {
		m = d().WithMatchWholeWord(false).ScanProfanity("xASs")
		assert.Equal(t, &Match{Word: "ass", Start: 1, End: 4, WordType: WordTypeProfanity,
			Text: []rune("ASs"), LeadingSpace: false, TrailingSpace: true}, toCmp(m[0]))

		m = d().WithMatchWholeWord(true).ScanProfanity("x aSs")
		assert.Equal(t, &Match{Word: "ass", Start: 2, End: 5, WordType: WordTypeProfanity,
			Text: []rune("aSs"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))

		m = d().WithMatchWholeWord(true).ScanProfanity("xAss")
		assert.Nil(t, m)
	})

	t.Run("Sanitize leet speak tests", func(t *testing.T) {
		m = d().WithSanitizeLeetSpeak(true).ScanProfanity("x a$$hol3 ")
		assert.Equal(t, &Match{Word: "asshole", Start: 2, End: 9, WordType: WordTypeProfanity,
			Text: []rune("a$$hol3"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeLeetSpeak(true).WithMatchWholeWord(false).ScanProfanity("x$sh!tx")
		assert.Equal(t, &Match{Word: "shit", Start: 1, End: 6, WordType: WordTypeProfanity,
			Text: []rune("$sh!t"), LeadingSpace: false, TrailingSpace: false}, toCmp(m[0]))
	})

	t.Run("Sanitize space char tests", func(t *testing.T) {
		m = d().WithSanitizeSpaces(true).ScanProfanity("x A S s")
		assert.Equal(t, &Match{Word: "ass", Start: 2, End: 7, WordType: WordTypeProfanity,
			Text: []rune("A S s"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeSpaces(false).ScanProfanity("x A ss")
		assert.Nil(t, m)
	})

	t.Run("Sanitize special char tests", func(t *testing.T) {
		m = d().WithSanitizeSpecialCharacters(true).ScanProfanity("x-f$u_c%k")
		assert.Equal(t, &Match{Word: "fuck", Start: 2, End: 9, WordType: WordTypeProfanity,
			Text: []rune("f$u_c%k"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeSpecialCharacters(false).ScanProfanity("x-f$u_c%k")
		assert.Nil(t, m)
	})

	t.Run("Sanitize repeated char tests", func(t *testing.T) {
		m = d().WithSanitizeRepeatedCharacters(true).ScanProfanity("x-fuu_cck")
		assert.Equal(t, &Match{Word: "fuck", Start: 2, End: 9, WordType: WordTypeProfanity,
			Text: []rune("fuu_cck"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeRepeatedCharacters(false).ScanProfanity("x-fuu_cck")
		assert.Nil(t, m)
	})

	t.Run("Sanitize wildcard char tests", func(t *testing.T) {
		m = d().WithSanitizeWildcardCharacters(true).ScanProfanity("x sh**t")
		assert.Equal(t, &Match{Word: "shit", Start: 2, End: 7, WordType: WordTypeProfanity,
			Text: []rune("sh**t"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeWildcardCharacters(false).ScanProfanity("x sh**t")
		assert.Nil(t, m)
	})

	t.Run("Process input as HTML tests", func(t *testing.T) {
		m = d().WithProcessInputAsHTML(true).ScanProfanity("x <a tag> &lt;ock </> ")
		assert.Equal(t, &Match{Word: "cock", Start: 10, End: 17, WordType: WordTypeProfanity,
			Text: []rune("&lt;ock"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))

		m = d().WithProcessInputAsHTML(false).ScanProfanity("x &lt;ock </> ")
		assert.Nil(t, m)
	})

	t.Run("Sanitize accents tests", func(t *testing.T) {
		m = d().WithSanitizeAccents(true).ScanProfanity("aa fúck")
		assert.Equal(t, &Match{Word: "fuck", Start: 3, End: 7, WordType: WordTypeProfanity,
			Text: []rune("fúck"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeAccents(false).ScanProfanity("aa fúck")
		assert.Nil(t, m)
	})

	t.Run("False positive tests", func(t *testing.T) {
		m = d().WithMatchWholeWord(false).ScanProfanity("x bada$$ ")
		assert.Equal(t, &Match{Word: "badass", Start: 2, End: 8, WordType: WordTypeFalsePositive,
			Text: []rune("bada$$"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))

		m = d().WithMatchWholeWord(true).ScanProfanity("x bada$$ ana7")
		assert.Equal(t, &Match{Word: "badass", Start: 2, End: 8, WordType: WordTypeFalsePositive,
			Text: []rune("bada$$"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))
		assert.Equal(t, &Match{Word: "anal", Start: 9, End: 13, WordType: WordTypeProfanity,
			Text: []rune("ana7"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[1]))

		m = d().WithMatchWholeWord(false).ScanProfanity("x aNalytic ")
		assert.Equal(t, &Match{Word: "analytic", Start: 2, End: 10, WordType: WordTypeFalsePositive,
			Text: []rune("aNalytic"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))

		m = d().WithMatchWholeWord(true).ScanProfanity("x analytic ")
		assert.Equal(t, &Match{Word: "analytic", Start: 2, End: 10, WordType: WordTypeFalsePositive,
			Text: []rune("analytic"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))
	})

	t.Run("Suspect word tests", func(t *testing.T) {
		m = d().WithSuspectWords([]string{"suspect"}).ScanProfanity("suspect sh!t")
		assert.Equal(t, &Match{Word: "suspect", Start: 0, End: 7, WordType: WordTypeSuspect,
			Text: []rune("suspect"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[0]))
		assert.Equal(t, &Match{Word: "shit", Start: 8, End: 12, WordType: WordTypeProfanity,
			Text: []rune("sh!t"), LeadingSpace: true, TrailingSpace: true}, toCmp(m[1]))
	})

	t.Run("Custom confidence calculator tests", func(t *testing.T) {
		calculator := func(m *Match) bool { return false }

		m = d().WithConfidenceCalculator(calculator).ScanProfanity("fuck")
		assert.Nil(t, m)

		m = d().ScanProfanity("fuck", WithConfidenceCalculator(calculator))
		assert.Nil(t, m)
	})
}

func Test_Censor(t *testing.T) {
	d := newDetectorEN
	var s string

	s, _ = d().Censor("x ass")
	assert.Equal(t, "x ***", s)

	s, _ = d().Censor("fuck this $h!!t")
	assert.Equal(t, "**** this *****", s)

	s, _ = d().Censor("clean text")
	assert.Equal(t, "clean text", s)

	s, _ = d().WithCensorCharacter('#').Censor("bada$$ a $ $")
	assert.Equal(t, "bada$$ # # #", s)
}

func Test_IsProfane(t *testing.T) {
	d := newDetectorEN

	assert.Equal(t, true, d().IsProfane("x ass"))
	assert.Equal(t, false, d().IsProfane("xass"))
	assert.Equal(t, true, d().IsProfane("xass", WithMatchWholeWord(false)))
}
