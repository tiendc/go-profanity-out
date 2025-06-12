package profanityout

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	firstRuneSupported = ' '
	lastRuneSupported  = '~'
)

var (
	removeAccentsTransformer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
)

// removeAccents strips all accents from characters
func removeAccents(s string) string {
	for _, character := range s {
		// If there's a character outside the range of supported runes, there might be some accented words
		if character < firstRuneSupported || character > lastRuneSupported {
			ss, _, err := transform.String(removeAccentsTransformer, s)
			if err == nil {
				s = ss
			}
			break
		}
	}
	return s
}

func normalizeAsNFC(s string) string {
	for _, character := range s {
		// If there's a character outside the range of supported runes, there might be some accented words
		if character < firstRuneSupported || character > lastRuneSupported {
			ss, _, err := transform.String(norm.NFC, s)
			if err == nil {
				s = ss
			}
			break
		}
	}
	return s
}
