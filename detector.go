package profanityout

type ProfanityDetector struct {
	settings            DetectorSettings
	specialCharacters   map[rune]rune
	leetSpeakCharacters map[rune]rune
	wildcardCharacters  map[rune]rune
	tree                *tree
}

func NewProfanityDetector() *ProfanityDetector {
	return &ProfanityDetector{
		settings: DetectorSettings{
			SanitizeSpecialCharacters:  true,
			SanitizeLeetSpeak:          true,
			SanitizeSpaces:             true,
			SanitizeRepeatedCharacters: true,
			SanitizeWildcardCharacters: true,
			SanitizeAccents:            true,
			ConfidenceCalculator:       confidenceCalculator,
			CensorCharacter:            '*',
		},
		tree: newTree(),
	}
}

// WithProfaneWords sets profane words
func (d *ProfanityDetector) WithProfaneWords(profaneWords []string) *ProfanityDetector {
	for _, word := range profaneWords {
		d.tree.Add(word, WordTypeProfanity)
	}
	return d
}

// WithSuspectWords sets suspect words
func (d *ProfanityDetector) WithSuspectWords(suspectWords []string) *ProfanityDetector {
	for _, word := range suspectWords {
		d.tree.Add(word, WordTypeSuspect)
	}
	return d
}

// WithFalsePositiveWords sets false positive words
func (d *ProfanityDetector) WithFalsePositiveWords(falsePositives []string) *ProfanityDetector {
	for _, word := range falsePositives {
		d.tree.Add(word, WordTypeFalsePositive)
	}
	return d
}

// WithLeetSpeakCharacters sets leet speak character map
func (d *ProfanityDetector) WithLeetSpeakCharacters(leetSpeakChars map[rune]rune) *ProfanityDetector {
	d.leetSpeakCharacters = leetSpeakChars
	return d
}

// WithSpecialCharacters sets special character map
func (d *ProfanityDetector) WithSpecialCharacters(specialChars map[rune]rune) *ProfanityDetector {
	d.specialCharacters = specialChars
	return d
}

// WithWildcardCharacters sets wildcard character map
func (d *ProfanityDetector) WithWildcardCharacters(wildcardChars map[rune]rune) *ProfanityDetector {
	d.wildcardCharacters = wildcardChars
	return d
}

// WithSanitizeLeetSpeak allows configuring whether the sanitization process should also
// take into account leetspeak.
//
// For instance, '4' is replaced by 'a' and '3' is replaced by 'e', which means that "4sshol3" would be
// sanitized to "asshole", which would be detected as a profanity.
func (d *ProfanityDetector) WithSanitizeLeetSpeak(sanitize bool) *ProfanityDetector {
	d.settings.SanitizeLeetSpeak = sanitize
	return d
}

// WithSanitizeSpecialCharacters allows configuring whether the sanitization process should
// also take into account special characters.
//
// For instance, "fu_ck" might be sanitized to "fuck", which would be detected as a profanity.
func (d *ProfanityDetector) WithSanitizeSpecialCharacters(sanitize bool) *ProfanityDetector {
	d.settings.SanitizeSpecialCharacters = sanitize
	return d
}

func (d *ProfanityDetector) WithSanitizeSpaces(sanitize bool) *ProfanityDetector {
	d.settings.SanitizeSpaces = sanitize
	return d
}

// WithSanitizeAccents allows configuring of whether the sanitization process should also
// take into account accents.
//
// For instance, "fúck" might be sanitized to "fuck", which would be detected as a profanity.
func (d *ProfanityDetector) WithSanitizeAccents(sanitize bool) *ProfanityDetector {
	d.settings.SanitizeAccents = sanitize
	return d
}

// WithSanitizeRepeatedCharacters allows configuring of whether the sanitization process should also take
// into account repeated characters.
//
// For instance, "fuuck" might be sanitized to "fuck", which would be detected as a profanity.
func (d *ProfanityDetector) WithSanitizeRepeatedCharacters(sanitize bool) *ProfanityDetector {
	d.settings.SanitizeRepeatedCharacters = sanitize
	return d
}

// WithSanitizeWildcardCharacters allows configuring of whether the sanitization process should also take
// into account wildcard characters.
//
// For instance, "f**k" might be sanitized to "fuck", which would be detected as a profanity.
func (d *ProfanityDetector) WithSanitizeWildcardCharacters(sanitize bool) *ProfanityDetector {
	d.settings.SanitizeWildcardCharacters = sanitize
	return d
}

// WithProcessInputAsHTML allows configuring of whether the sanitization process should also take
// into account HTML content.
//
// For instance, all HTML tags in the input will be removed and all HTMl entities will be replaced
// by real characters (for example, &gt; will be replaced with '>').
func (d *ProfanityDetector) WithProcessInputAsHTML(asHTML bool) *ProfanityDetector {
	d.settings.ProcessInputAsHTML = asHTML
	return d
}

// WithConfidenceCalculator sets custom confidence calculator function
func (d *ProfanityDetector) WithConfidenceCalculator(calculator ConfidenceCalculator) *ProfanityDetector {
	d.settings.ConfidenceCalculator = calculator
	return d
}

// WithCensorCharacter sets custom censor character (default: *)
func (d *ProfanityDetector) WithCensorCharacter(censorCharacter rune) *ProfanityDetector {
	d.settings.CensorCharacter = censorCharacter
	return d
}

// IsProfane checks a string containing profanity or not
func (d *ProfanityDetector) IsProfane(s string, options ...DetectorOption) bool {
	return d.ScanProfanity(s, options...).HasProfaneMatch()
}

// ScanProfanity scans for the first profanity
func (d *ProfanityDetector) ScanProfanity(s string, options ...DetectorOption) Matches {
	return d.newScanner(false, options...).scan(s)
}

// ScanAllProfanities scans for all profanities
func (d *ProfanityDetector) ScanAllProfanities(s string, options ...DetectorOption) (matches Matches) {
	return d.newScanner(true, options...).scan(s)
}

// Censor scans for all profanities and censors all of them if found
func (d *ProfanityDetector) Censor(s string, options ...DetectorOption) (string, Matches) {
	scanner := d.newScanner(true, options...)
	matches := scanner.scan(s)
	if len(matches) == 0 {
		return s, nil
	}

	content := scanner.inputOrig
	for _, match := range matches.GetProfaneMatches() {
		for i := match.Start; i < match.End; i++ {
			if content[i] != ' ' {
				content[i] = scanner.settings.CensorCharacter
			}
		}
	}
	return string(content), matches
}

func (d *ProfanityDetector) newScanner(findAllMatches bool, options ...DetectorOption) *scanner {
	settings := d.settings
	settings.findAllProfanityMatches = findAllMatches
	for _, opt := range options {
		opt(&settings)
	}
	return &scanner{
		settings:            &settings,
		specialCharacters:   d.specialCharacters,
		leetSpeakCharacters: d.leetSpeakCharacters,
		wildcardCharacters:  d.wildcardCharacters,
		tree:                d.tree,
	}
}

func confidenceCalculator(match *Match) bool {
	return true
}
