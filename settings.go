package profanityout

type DetectorSettings struct {
	SanitizeSpecialCharacters  bool
	SanitizeLeetSpeak          bool
	SanitizeAccents            bool
	SanitizeSpaces             bool
	SanitizeRepeatedCharacters bool
	SanitizeWildcardCharacters bool
	ProcessInputAsHTML         bool
	MatchWholeWord             bool

	ConfidenceCalculator ConfidenceCalculator
	CensorCharacter      rune

	// this is used internally
	findAllProfanityMatches bool
}

type DetectorOption func(*DetectorSettings)
type ConfidenceCalculator func(*Match) bool

func WithSanitizeSpecialCharacters(flag bool) DetectorOption {
	return func(settings *DetectorSettings) {
		settings.SanitizeSpecialCharacters = flag
	}
}

func WithSanitizeLeetSpeak(flag bool) DetectorOption {
	return func(settings *DetectorSettings) {
		settings.SanitizeLeetSpeak = flag
	}
}

func WithSanitizeAccents(flag bool) DetectorOption {
	return func(settings *DetectorSettings) {
		settings.SanitizeAccents = flag
	}
}

func WithSanitizeRepeatedCharacters(flag bool) DetectorOption {
	return func(settings *DetectorSettings) {
		settings.SanitizeRepeatedCharacters = flag
	}
}

func WithSanitizeWildcardCharacters(flag bool) DetectorOption {
	return func(settings *DetectorSettings) {
		settings.SanitizeWildcardCharacters = flag
	}
}

func WithProcessInputAsHTML(flag bool) DetectorOption {
	return func(settings *DetectorSettings) {
		settings.ProcessInputAsHTML = flag
	}
}

func WithMatchWholeWord(flag bool) DetectorOption {
	return func(settings *DetectorSettings) {
		settings.MatchWholeWord = flag
	}
}

func WithConfidenceCalculator(fn ConfidenceCalculator) DetectorOption {
	return func(settings *DetectorSettings) {
		settings.ConfidenceCalculator = fn
	}
}

func WithCensorCharacter(ch rune) DetectorOption {
	return func(settings *DetectorSettings) {
		settings.CensorCharacter = ch
	}
}
