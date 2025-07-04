package profanityout

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tiendc/go-profanity-out/data/en"
)

const (
	sanitizeAccents            = true
	sanitizeLeetSpeak          = true
	sanitizeSpecialCharacters  = true
	sanitizeSpaces             = true
	sanitizeRepeatedCharacters = true
	sanitizeWildcardCharacters = true
	processInputAsHTML         = false
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
		WithProcessInputAsHTML(processInputAsHTML)
}

func toCmp(m *Match) *Match {
	return &Match{
		Word:      m.Word,
		WordType:  m.WordType,
		Start:     m.Start,
		End:       m.End,
		HeadSpace: m.HeadSpace,
		TailSpace: m.TailSpace,
		Text:      m.Text,
	}
}

func Test_Scan_One(t *testing.T) {
	d := newDetectorEN
	var m []*Match

	t.Run("Match whole word tests", func(t *testing.T) {
		m = d().ScanProfanity("x ASs-x")
		assert.Equal(t, &Match{Word: "ass", Start: 2, End: 5, WordType: WordTypeProfanity,
			Text: []rune("ASs"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().ScanProfanity("x aSs")
		assert.Equal(t, &Match{Word: "ass", Start: 2, End: 5, WordType: WordTypeProfanity,
			Text: []rune("aSs"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().ScanProfanity("xAss")
		assert.Nil(t, m)
	})

	t.Run("Sanitize leet speak tests", func(t *testing.T) {
		m = d().WithSanitizeLeetSpeak(true).ScanProfanity("x $ h!t ")
		assert.Equal(t, &Match{Word: "shit", Start: 2, End: 7, WordType: WordTypeProfanity,
			Text: []rune("$ h!t"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeLeetSpeak(true).ScanProfanity("x a$$hol3 ")
		assert.Equal(t, &Match{Word: "asshole", Start: 2, End: 9, WordType: WordTypeProfanity,
			Text: []rune("a$$hol3"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeLeetSpeak(true).ScanProfanity("x_$sh!t_x")
		assert.Equal(t, &Match{Word: "shit", Start: 2, End: 7, WordType: WordTypeProfanity,
			Text: []rune("$sh!t"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeLeetSpeak(true).ScanProfanity("bada$$ a $ $")
		assert.Equal(t, &Match{Word: "ass", Start: 5, End: 12, WordType: WordTypeProfanity,
			Text: []rune("$ a $ $"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))
	})

	t.Run("Sanitize space char tests", func(t *testing.T) {
		m = d().WithSanitizeSpaces(true).ScanProfanity("x A S s")
		assert.Equal(t, &Match{Word: "ass", Start: 2, End: 7, WordType: WordTypeProfanity,
			Text: []rune("A S s"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeSpaces(false).ScanProfanity("x A ss")
		assert.Nil(t, m)
	})

	t.Run("Sanitize special char tests", func(t *testing.T) {
		m = d().WithSanitizeSpecialCharacters(true).ScanProfanity("x-f$u_c%k")
		assert.Equal(t, &Match{Word: "fuck", Start: 2, End: 9, WordType: WordTypeProfanity,
			Text: []rune("f$u_c%k"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeSpecialCharacters(true).ScanProfanity("x_f$u_c%k-")
		assert.Equal(t, &Match{Word: "fuck", Start: 2, End: 9, WordType: WordTypeProfanity,
			Text: []rune("f$u_c%k"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeSpecialCharacters(false).ScanProfanity("x-f$u_c%k")
		assert.Nil(t, m)
	})

	t.Run("Sanitize repeated char tests", func(t *testing.T) {
		m = d().WithSanitizeRepeatedCharacters(true).ScanProfanity("x-fuu_cck")
		assert.Equal(t, &Match{Word: "fuck", Start: 2, End: 9, WordType: WordTypeProfanity,
			Text: []rune("fuu_cck"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeRepeatedCharacters(false).ScanProfanity("x-fuu_cck")
		assert.Nil(t, m)
	})

	t.Run("Sanitize wildcard char tests", func(t *testing.T) {
		// Wildcard in input
		m = d().WithSanitizeWildcardCharacters(true).ScanProfanity("x sh*t")
		assert.Equal(t, &Match{Word: "shit", Start: 2, End: 6, WordType: WordTypeProfanity,
			Text: []rune("sh*t"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeWildcardCharacters(false).ScanProfanity("x sh**t")
		assert.Nil(t, m)

		// Wildcard in profanity dictionary
		m = d().WithProfaneWords([]string{"*blah*"}).WithSanitizeWildcardCharacters(true).
			ScanProfanity("xbl@hx")
		assert.Equal(t, &Match{Word: "blah", Start: 1, End: 5, WordType: WordTypeProfanity,
			Text: []rune("bl@h"), HeadSpace: false, TailSpace: false}, toCmp(m[0]))

		m = d().WithProfaneWords([]string{"*blah*"}).WithSanitizeWildcardCharacters(true).
			ScanProfanity("xxbl@h")
		assert.Equal(t, &Match{Word: "blah", Start: 2, End: 6, WordType: WordTypeProfanity,
			Text: []rune("bl@h"), HeadSpace: false, TailSpace: true}, toCmp(m[0]))

		m = d().WithProfaneWords([]string{"*blah*"}).WithSanitizeWildcardCharacters(true).
			ScanProfanity(" bl@h")
		assert.Equal(t, &Match{Word: "blah", Start: 1, End: 5, WordType: WordTypeProfanity,
			Text: []rune("bl@h"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().WithProfaneWords([]string{"foo*bar"}).WithSanitizeWildcardCharacters(true).
			ScanProfanity("foobar")
		assert.Equal(t, &Match{Word: "foobar", Start: 0, End: 6, WordType: WordTypeProfanity,
			Text: []rune("foobar"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))
	})

	t.Run("Process input as HTML tests", func(t *testing.T) {
		m = d().WithProcessInputAsHTML(true).ScanProfanity("x <a tag> &lt;ock </> ")
		assert.Equal(t, &Match{Word: "cock", Start: 10, End: 17, WordType: WordTypeProfanity,
			Text: []rune("&lt;ock"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().WithProcessInputAsHTML(false).ScanProfanity("x &lt;ock </> ")
		assert.Nil(t, m)
	})

	t.Run("Sanitize accents tests", func(t *testing.T) {
		m = d().WithSanitizeAccents(true).ScanProfanity("aa fúck")
		assert.Equal(t, &Match{Word: "fuck", Start: 3, End: 7, WordType: WordTypeProfanity,
			Text: []rune("fúck"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))

		m = d().WithSanitizeAccents(false).ScanProfanity("aa fúck")
		assert.Nil(t, m)
	})

	t.Run("False positive tests", func(t *testing.T) {
		m = d().WithProfaneWords([]string{"*anal*"}).ScanProfanity("x-analytic-ana7")
		assert.Equal(t, &Match{Word: "analy", Start: 2, End: 7, WordType: WordTypeFalsePositive,
			Text: []rune("analy"), HeadSpace: true, TailSpace: false}, toCmp(m[0]))
		assert.Equal(t, &Match{Word: "anal", Start: 11, End: 15, WordType: WordTypeProfanity,
			Text: []rune("ana7"), HeadSpace: true, TailSpace: true}, toCmp(m[1]))

		m = d().WithFalsePositiveWords([]string{"ASS"}).ScanProfanity("x ASS x")
		assert.Equal(t, &Match{Word: "ASS", Start: 2, End: 5, WordType: WordTypeFalsePositive,
			Text: []rune("ASS"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))
	})

	t.Run("Suspect word tests", func(t *testing.T) {
		m = d().WithSuspectWords([]string{"suspect"}).ScanProfanity("suspect sh!t")
		assert.Equal(t, &Match{Word: "suspect", Start: 0, End: 7, WordType: WordTypeSuspect,
			Text: []rune("suspect"), HeadSpace: true, TailSpace: true}, toCmp(m[0]))
		assert.Equal(t, &Match{Word: "shit", Start: 8, End: 12, WordType: WordTypeProfanity,
			Text: []rune("sh!t"), HeadSpace: true, TailSpace: true}, toCmp(m[1]))
	})

	t.Run("Custom confidence calculator tests", func(t *testing.T) {
		calculator := func(m *Match) bool { return false }

		m = d().WithConfidenceCalculator(calculator).ScanProfanity("fuck")
		assert.Nil(t, m)

		m = d().ScanProfanity("fuck", WithConfidenceCalculator(calculator))
		assert.Nil(t, m)
	})
}

func Test_Scan_All(t *testing.T) {
	// TODO: add tests for this
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

	s, _ = d().WithCensorCharacter('#').Censor("bada$s a $ $")
	assert.Equal(t, "bada$s # # #", s)
}

func Test_IsProfane(t *testing.T) {
	d := newDetectorEN

	assert.Equal(t, true, d().IsProfane("x ass"))
	assert.Equal(t, false, d().IsProfane("xass"))
}

func Test_General(t *testing.T) {
	type TestCase struct {
		phrase    string
		offensive bool
	}

	//nolint:lll
	var generalTestCases = []TestCase{
		{"hi", false},
		{"hello", false},
		{"hello my name is Bob.", false},
		{"SHIT", true},
		{"shhhhhiiiiter", true},
		{"shhHhhit", true},
		{"lol fuck this", true},
		{"f*u*c*k", true},
		{"$#1t", true},
		{" fučk", true},
		{"ass", true},
		{"glass", false},
		{"ÄšŚ", true},
		{"ĂżŽ", false}, // TODO: unable to detect this case now
		{"sex", true},
		{"hello_world-sex_word", true},
		{"sexy", true},
		{"is extreme", false},
		{"pÓöp", true},
		{"what a bunch of bullsh1t", false},
		{"bitčh", true},
		{"assassin", false},
		{"push it", false},
		{"carcass", false},
		{"retarded", true},
		{"βιτ⊂η", true}, // greek letters
		{"ⓅɄȿⓢⓨ", true},
		{"I had called upon my friend, Mr. Sherlock Holmes, one day in the autumn of last year and found him in deep conversation with a very stout, florid-faced, elderly gentleman with fiery red hair.", false},
		{"With an apology for my intrusion, I was about to withdraw when Holmes pulled me abruptly into the room and closed the door behind me.", false},
		{"You could not possibly have come at a better time, my dear Watson, he said cordially", false},
		{"I was afraid that you were engaged.", false},
		{"So I am. Very much so.", false},
		{"Then I can wait in the next room.", false},
		{"Not at all. This gentleman, Mr. Wilson, has been my partner and helper in many of my most successful cases, and I have no doubt that he will be of the utmost use to me in yours also.", false},
		{"The stout gentleman half rose from his chair and gave a bob of greeting, with a quick little questioning glance from his small fat-encircled eyes", false},
		{"Try the settee, said Holmes, relapsing into his armchair and putting his fingertips together, as was his custom when in judicial moods.", false},
		{"I know, my dear Watson, that you share my love of all that is bizarre and outside the conventions and humdrum routine of everyday life.", false},
		{"You have shown your relish for it by the enthusiasm which has prompted you to chronicle, and, if you will excuse my saying so, somewhat to embellish so many of my own little adventures.", false},
		{"You did, Doctor, but none the less you must come round to my view, for otherwise I shall keep on piling fact upon fact on you until your reason breaks down under them and acknowledges me to be right.", false},
		{"Now, Mr. Jabez Wilson here has been good enough to call upon me this morning, and to begin a narrative which promises to be one of the most singular which I have listened to for some time.", false},
		{"You have heard me remark that the strangest and most unique things are very often connected not with the larger but with the smaller crimes, and occasionally", false},
		{"indeed, where there is room for doubt whether any positive crime has been committed.", false},
		{"As far as I have heard it is impossible for me to say whether the present case is an instance of crime or not, but the course of events is certainly among the most singular that I have ever listened to.", false},
		{"Perhaps, Mr. Wilson, you would have the great kindness to recommence your narrative.", false},
		{"I ask you not merely because my friend Dr. Watson has not heard the opening part but also because the peculiar nature of the story makes me anxious to have every possible detail from your lips.", false},
		{"As a rule, when I have heard some slight indication of the course of events, I am able to guide myself by the thousands of other similar cases which occur to my memory.", false},
		{"In the present instance I am forced to admit that the facts are, to the best of my belief, unique.", false},
		{"We had reached the same crowded thoroughfare in which we had found ourselves in the morning.", false},
		{"Our cabs were dismissed, and, following the guidance of Mr. Merryweather, we passed down a narrow passage and through a side door, which he opened for us", false},
		{"Within there was a small corridor, which ended in a very massive iron gate.", false},
		{"We were seated at breakfast one morning, my wife and I, when the maid brought in a telegram. It was from Sherlock Holmes and ran in this way", false},
	}

	d := newDetectorEN()
	for _, testCase := range generalTestCases {
		inappropriate := d.ScanProfanity(testCase.phrase).HasProfaneMatch()
		if inappropriate != testCase.offensive {
			t.Errorf("phrase=\"%s\" analysis offensive=%v actual offensive=%v",
				testCase.phrase, inappropriate, testCase.offensive)
		}
	}
}
