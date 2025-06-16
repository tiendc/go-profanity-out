package profanityout

type Match struct {
	Word      string
	WordType  WordType
	Start     int
	End       int
	HeadSpace bool
	TailSpace bool
	Text      []rune
	Settings  *DetectorSettings
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

type Matches []*Match

func (ms Matches) HasProfaneMatch() bool {
	for _, m := range ms {
		if m.IsProfane() {
			return true
		}
	}
	return false
}

func (ms Matches) GetProfaneMatches() (resp Matches) {
	for _, m := range ms {
		if m.IsProfane() {
			resp = append(resp, m)
		}
	}
	return resp
}

func (ms Matches) GetFirstProfaneMatch() *Match {
	for _, m := range ms {
		if m.IsProfane() {
			return m
		}
	}
	return nil
}

func (ms Matches) HasSuspectMatch() bool {
	for _, m := range ms {
		if m.IsSuspect() {
			return true
		}
	}
	return false
}

func (ms Matches) GetSuspectMatches() (resp Matches) {
	for _, m := range ms {
		if m.IsSuspect() {
			resp = append(resp, m)
		}
	}
	return resp
}

func (ms Matches) HasFalsePositiveMatch() bool {
	for _, m := range ms {
		if m.IsFalsePositive() {
			return true
		}
	}
	return false
}

func (ms Matches) GetFalsePositiveMatches() (resp Matches) {
	for _, m := range ms {
		if m.IsFalsePositive() {
			resp = append(resp, m)
		}
	}
	return resp
}
