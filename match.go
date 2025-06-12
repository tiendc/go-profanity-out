package profanityout

type Match struct {
	Word          string
	WordType      WordType
	Start         int
	End           int
	LeadingSpace  bool
	TrailingSpace bool
	Text          []rune
	Settings      *DetectorSettings
}

type WordType int8

func (m *Match) IsProfane() bool       { return m.WordType == WordTypeProfanity }
func (m *Match) IsSuspect() bool       { return m.WordType == WordTypeSuspect }
func (m *Match) IsFalsePositive() bool { return m.WordType == WordTypeFalsePositive }

const (
	// NOTE: Order matter, Suspects < Profanities < False Positives
	WordTypeSuspect       WordType = 10
	WordTypeProfanity     WordType = 20
	WordTypeFalsePositive WordType = 30
)
