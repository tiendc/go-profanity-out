package profanityout

import (
	"unicode"
)

type scanner struct {
	settings            *DetectorSettings
	specialCharacters   map[rune]rune
	leetSpeakCharacters map[rune]rune
	wildcardCharacters  map[rune]rune
	tree                *tree

	inputOrig []rune
	input     []rune
	inputLen  int
}

func (s *scanner) scanAll(input string) (matches []*Match) {
	// Sanitizes accents if configured
	if s.settings.SanitizeAccents {
		s.inputOrig = []rune(input)
		s.input = []rune(removeAccents(input))
	} else {
		s.inputOrig = []rune(normalizeAsNFC(input))
		s.input = s.inputOrig
	}

	s.inputLen = len(s.input)
	pos := 0
	for pos < s.inputLen {
		nextPos := s.skipWhitespaces(pos)
		if nextPos >= s.inputLen {
			break
		}
		hasLeadingSpace := pos == 0 || nextPos != pos
		pos = nextPos

		match := &Match{Start: pos, LeadingSpace: hasLeadingSpace, Settings: s.settings}
		s.scanOne(pos, 0, s.tree.root, match)
		if match.WordType != 0 {
			if !s.settings.ConfidenceCalculator(match) {
				goto ScanNextPos
			}
			matches = append(matches, match)
			if match.WordType == WordTypeProfanity && !s.settings.findAllProfanityMatches {
				return matches
			}
			pos = match.End
			continue
		}

	ScanNextPos:
		// No match found
		if s.settings.MatchWholeWord {
			nextPos = s.skipUntilWhitespace(pos)
			if nextPos != pos {
				pos = nextPos
				continue
			}
		}
		pos++
	}

	return matches
}

//nolint:gocognit,gocyclo
func (s *scanner) scanOne(pos int, prevCh rune, currentNode *node, match *Match) {
	// Logging:
	// fmt.Printf("------------\n")
	// fmt.Printf("Scanning start: pos=%d text='%s'\n", pos, string(s.input[pos:]))

	wildcardPos := -1
	var wildcardNode *node

	for pos < s.inputLen {
		ch, nextPos := s.nextCharAt(pos)
		if ch == 0 {
			break
		}

		ch = unicode.ToLower(ch)
		nextNode := currentNode.Next(ch)
		if nextNode == nil { //nolint:nestif
			if ch == ' ' && s.settings.SanitizeSpaces {
				pos = nextPos
				prevCh = ch
				continue
			}

			if s.settings.SanitizeWildcardCharacters && wildcardPos == -1 {
				if _, exists := s.wildcardCharacters[ch]; exists {
					// Stores the pos we may start a new scan from when no matching found
					wildcardPos, wildcardNode = pos, currentNode
				}
			}

			if s.settings.SanitizeLeetSpeak {
				if lsCh, exists := s.leetSpeakCharacters[ch]; exists {
					if lsNode := currentNode.Next(lsCh); lsNode != nil {
						if lsNode.wordType != 0 { // match found at the current node
							s.updateMatchWithFoundNode(match, nextPos, lsNode)
						}
						s.scanOne(nextPos, lsCh, lsNode, match)  // scan deeper
						if match.WordType == WordTypeProfanity { // found a profanity, return
							break
						}
					}
				}
			}

			if s.settings.SanitizeSpecialCharacters {
				if specialCh, exists := s.specialCharacters[ch]; exists {
					ch = specialCh
					nextNode = currentNode.Next(ch)
					if nextNode != nil {
						goto HandleNodeFound
					}
					if ch == ' ' && s.settings.SanitizeSpaces {
						pos = nextPos
						prevCh = ch
						continue
					}
				}
			}

			if s.settings.SanitizeRepeatedCharacters && s.isCharRepeatedAt(pos, prevCh) {
				pos = nextPos
				prevCh = ch
				continue
			}

			break
		}

	HandleNodeFound:
		// At this point: nextNode != nil
		pos = nextPos
		prevCh = ch
		currentNode = nextNode

		// If there is a matching detected
		if currentNode.wordType != 0 {
			s.updateMatchWithFoundNode(match, nextPos, currentNode)
		}
	}

	// After all scans and no matching found, we may start a new scan for wildcard matching
	if match.WordType == 0 && wildcardPos >= 0 {
		for currCh, currNode := range wildcardNode.children {
			if currNode.wordType != 0 { // match found at the current node
				s.updateMatchWithFoundNode(match, wildcardPos+1, currNode)
			}
			s.scanOne(wildcardPos+1, currCh, currNode, match) // scan deeper
			if match.WordType == WordTypeProfanity {          // found a profanity, return
				break
			}
		}
	}

	// Logging:
	// fmt.Printf("Match Data: type=%d word='%s'\n", match.WordType, match.Word)
}

func (s *scanner) isWhitespaceAt(i int) bool {
	if i < 0 {
		return true
	}
	ch, _ := s.nextCharAt(i)
	if ch == 0 {
		return true
	}
	if s.settings.SanitizeSpecialCharacters {
		if ch2 := s.specialCharacters[ch]; ch2 != 0 {
			ch = ch2
		}
	}
	return ch == ' '
}

func (s *scanner) isCharRepeatedAt(i int, prevCh rune) bool {
	if i <= 0 {
		return false
	}
	ch, _ := s.nextCharAt(i)
	return unicode.ToLower(ch) == prevCh
}

func (s *scanner) skipWhitespaces(i int) int {
	for {
		ch, next := s.nextCharAt(i)
		if ch == 0 {
			return next
		}
		if ch == ' ' {
			i = next
			continue
		}
		if s.settings.SanitizeLeetSpeak {
			if ch2 := s.leetSpeakCharacters[ch]; ch2 != 0 && ch2 != ' ' {
				return i // found a non-whitespace
			}
		}
		if s.settings.SanitizeSpecialCharacters {
			if ch2 := s.specialCharacters[ch]; ch2 != 0 {
				ch = ch2
			}
		}
		if ch != ' ' {
			return i
		}
		i = next
	}
}

func (s *scanner) skipUntilWhitespace(i int) int {
	for {
		ch, next := s.nextCharAt(i)
		if ch == 0 {
			return next
		}
		if ch == ' ' {
			return i
		}
		if s.settings.SanitizeLeetSpeak {
			if ch2 := s.leetSpeakCharacters[ch]; ch2 == ' ' {
				return i // found a non-whitespace
			}
		}
		if s.settings.SanitizeSpecialCharacters {
			if ch2 := s.specialCharacters[ch]; ch2 == ' ' {
				return i
			}
		}
		i = next
	}
}

func (s *scanner) nextCharAt(i int) (rune, int) {
	if i >= s.inputLen {
		return 0, i
	}
	ch := s.input[i]
	if s.settings.ProcessInputAsHTML { //nolint:nestif
		if ch == '<' { // HTML tag opening
			next := skipHTMLTag(s.input, i)
			if next != i {
				return s.nextCharAt(next)
			}
			return ch, i + 1
		}
		if ch == '&' { // HTML entity beginning
			ch2, next := decodeHTMLEntityAt(s.input, i)
			if next == i {
				return ch, i + 1
			}
			return ch2, next
		}
	}
	return ch, i + 1
}

func (s *scanner) updateMatchWithFoundNode(match *Match, end int, node *node) {
	// Only update the match if the target is at equal or higher level.
	// For example, current is `profanity` and the target is `suspect`, just ignore.
	if match.WordType > node.wordType {
		return
	}
	trailingSpace := s.isWhitespaceAt(end)
	if s.settings.MatchWholeWord && (!match.LeadingSpace || !trailingSpace) {
		return
	}

	match.End = end
	match.WordType = node.wordType
	match.Word = node.word
	match.TrailingSpace = trailingSpace
	match.Text = s.inputOrig[match.Start:match.End]
}
