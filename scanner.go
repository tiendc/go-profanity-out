package profanityout

import (
	"unicode"
)

type scanner struct {
	settings            *DetectorSettings
	specialCharacters   map[rune]rune
	leetSpeakCharacters map[rune]rune
	wildcardCharacters  map[rune]rune
	profanityTree       *tree
	falsePositiveTree   *tree

	inputOrig []rune
	input     []rune
}

func (s *scanner) scan(input string) (matches Matches) {
	// Sanitizes accents if configured
	if s.settings.SanitizeAccents {
		s.inputOrig = []rune(input)
		s.input = []rune(removeAccents(input))
	} else {
		s.inputOrig = []rune(normalizeAsNFC(input))
		s.input = s.inputOrig
	}

	match := Match{} // declares a match here to reduce the allocations
	hasHeadingWildcard := s.settings.SanitizeWildcardCharacters && s.profanityTree.hasHeadingWildcard
	var prevCh rune
	pos := 0
	for {
		ch, nextPos := s.nextCharAt(pos)
		if ch == 0 {
			break
		}
		if !s.shouldStartScanning(ch) {
			prevCh = ch
			pos = nextPos
			continue
		}

		match = Match{Start: pos, HeadSpace: prevCh == 0 || s.isWhitespace(prevCh), Settings: s.settings}
		// Scans for a false positive first, if not found, scans for profanity
		if s.scanFalsePositive(pos, s.falsePositiveTree.root, &match); match.WordType == 0 {
			s.scanProfanity(pos, 0, s.profanityTree.root, &match)
		}

		if match.WordType != 0 {
			if !s.settings.ConfidenceCalculator(&match) {
				goto ScanNextPos
			}
			matchCopy := match
			matches = append(matches, &matchCopy)
			if match.WordType == WordTypeProfanity && !s.settings.findAllProfanityMatches {
				return matches
			}
			prevCh = ch
			pos = match.End
			continue
		}

	ScanNextPos:
		prevCh = ch
		if hasHeadingWildcard {
			pos = nextPos
			continue
		}
		if next := s.skipUntilWhitespace(pos); next != pos {
			pos = next
			continue
		}
		pos = nextPos
	}

	return matches
}

func (s *scanner) shouldStartScanning(ch rune) bool {
	if s.settings.SanitizeLeetSpeak && s.leetSpeakCharacters[ch] != 0 {
		return true
	}
	return !s.isWhitespace(ch)
}

//nolint:gocognit,gocyclo
func (s *scanner) scanProfanity(pos int, prevCh rune, currentNode *node, match *Match) {
	wildcardPos := -1
	var wildcardNode *node

	for {
		ch, nextPos := s.nextCharAt(pos)
		if ch == 0 {
			break
		}

		ch = unicode.ToLower(ch)
		nextNode := currentNode.Next(ch)
		if nextNode == nil { //nolint:nestif
			if s.settings.SanitizeLeetSpeak {
				if lsCh, exists := s.leetSpeakCharacters[ch]; exists {
					if lsNode := currentNode.Next(lsCh); lsNode != nil {
						match.foundRealCharMatch = true
						if lsNode.word != nil { // match found at the current node
							s.updateMatchWithFoundNode(match, nextPos, lsNode)
						}
						s.scanProfanity(nextPos, lsCh, lsNode, match) // scan deeper
						if match.WordType == WordTypeProfanity {      // found a profanity, return
							break
						}
					}
				}
			}

			if !match.foundRealCharMatch {
				break
			}

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

			if s.settings.SanitizeWildcardCharacters {
				nextNode = currentNode.Next('*')
				if nextNode != nil {
					goto HandleNodeFound
				}
			}

			break
		}

	HandleNodeFound:
		pos = nextPos
		prevCh = ch
		currentNode = nextNode
		if !match.foundRealCharMatch {
			match.foundRealCharMatch = !s.isWhitespace(ch)
		}

		// If there is a matching detected
		if currentNode.word != nil {
			s.updateMatchWithFoundNode(match, nextPos, currentNode)
		}
	}

	// After all scans and no matching found, we may start a new scan for wildcard matching
	if match.WordType == 0 && wildcardPos >= 0 {
		for currCh, currNode := range wildcardNode.children {
			if currNode.word != nil { // match found at the current node
				s.updateMatchWithFoundNode(match, wildcardPos+1, currNode)
			}
			s.scanProfanity(wildcardPos+1, currCh, currNode, match) // scan deeper
			if match.WordType == WordTypeProfanity {                // found a profanity, return
				break
			}
		}
	}

	// When found a profanity, we do extra scans to make sure it's not a false positive
	if match.WordType > 0 && match.WordType < WordTypeFalsePositive {
		for pos = match.Start; pos < match.End; pos++ {
			if s.scanExactFalsePositive(pos, match); match.WordType == WordTypeFalsePositive {
				break
			}
		}
	}
}

//nolint:gocognit,gocyclo
func (s *scanner) scanFalsePositive(pos int, currentNode *node, match *Match) {
	for {
		ch, nextPos := s.nextCharAt(pos)
		if ch == 0 {
			break
		}

		ch = unicode.ToLower(ch)
		nextNode := currentNode.Next(ch)
		if nextNode == nil { //nolint:nestif
			if s.settings.SanitizeLeetSpeak {
				if lsCh, exists := s.leetSpeakCharacters[ch]; exists {
					if lsNode := currentNode.Next(lsCh); lsNode != nil {
						if lsNode.word != nil { // match found at the current node
							s.updateMatchWithFoundNode(match, nextPos, lsNode)
						}
						s.scanFalsePositive(nextPos, lsNode, match)  // scan deeper
						if match.WordType == WordTypeFalsePositive { // found a false positive, return
							break
						}
					}
				}
			}

			if s.settings.SanitizeWildcardCharacters {
				nextNode = currentNode.Next('*')
				if nextNode != nil {
					goto HandleNodeFound
				}
			}

			break
		}

	HandleNodeFound:
		pos = nextPos
		currentNode = nextNode

		// If there is a matching detected
		if currentNode.word != nil {
			s.updateMatchWithFoundNode(match, nextPos, currentNode)
		}
	}
}

// scanExactFalsePositive cans for exact match of false positive without applying
// any transformation of casing, leet speak, or special characters
func (s *scanner) scanExactFalsePositive(pos int, match *Match) {
	currentNode := s.falsePositiveTree.root
	start := pos
	for {
		ch, nextPos := s.nextOrigCharAt(pos)
		if ch == 0 {
			break
		}

		nextNode := currentNode.Next(ch)
		if nextNode == nil {
			break
		}

		pos = nextPos
		currentNode = nextNode

		// If there is a matching detected
		if currentNode.word != nil {
			match.Start = start
			s.updateMatchWithFoundNode(match, pos, currentNode)
		}
	}
}

func (s *scanner) isWhitespace(ch rune) bool {
	if ch == ' ' {
		return true
	}
	if s.settings.SanitizeSpecialCharacters {
		if ch = s.specialCharacters[ch]; ch == ' ' {
			return true
		}
	}
	return false
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
	return s.nextCharOf(s.input, i)
}

func (s *scanner) nextOrigCharAt(i int) (rune, int) {
	return s.nextCharOf(s.inputOrig, i)
}

func (s *scanner) nextCharOf(input []rune, i int) (rune, int) {
	if i >= len(input) {
		return 0, i
	}
	ch := input[i]
	if s.settings.ProcessInputAsHTML { //nolint:nestif
		if ch == '<' { // HTML tag opening
			next := skipHTMLTag(input, i)
			if next != i {
				return s.nextCharOf(input, next)
			}
			return ch, i + 1
		}
		if ch == '&' { // HTML entity beginning
			ch2, next := decodeHTMLEntityAt(input, i)
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
	if match.WordType > node.word.wordType {
		return
	}

	tailSpace := s.isWhitespaceAt(end)
	if node.word.wordType < WordTypeFalsePositive {
		if !match.HeadSpace && node.word.wordFlag.RequireHeadSpace() {
			return
		}
		if !tailSpace && node.word.wordFlag.RequireTailSpace() {
			return
		}
	}

	match.End = end
	match.WordType = node.word.wordType
	match.Word = node.word.word
	match.TailSpace = tailSpace
	match.Text = s.inputOrig[match.Start:match.End]
}
