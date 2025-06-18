package profanityout

import "strings"

type tree struct {
	root               *node
	hasHeadingWildcard bool
}

type node struct {
	children map[rune]*node
	word     *wordData
}

type wordData struct {
	word     string
	wordType WordType
	wordFlag WordFlag
}

type WordFlag uint8

const (
	wordFlagRequireHeadSpace WordFlag = 1
	wordFlagRequireTailSpace WordFlag = 2

	wordFlagDefault WordFlag = wordFlagRequireHeadSpace | wordFlagRequireTailSpace
)

func (flag WordFlag) RequireHeadSpace() bool {
	return flag&wordFlagRequireHeadSpace != 0
}

func (flag *WordFlag) SetRequireHeadSpace(val bool) {
	if val {
		*flag |= wordFlagRequireHeadSpace
	} else {
		*flag &= ^wordFlagRequireHeadSpace
	}
}

func (flag WordFlag) RequireTailSpace() bool {
	return flag&wordFlagRequireTailSpace != 0
}

func (flag *WordFlag) SetRequireTailSpace(val bool) {
	if val {
		*flag |= wordFlagRequireTailSpace
	} else {
		*flag &= ^wordFlagRequireTailSpace
	}
}

func (node *node) Next(next rune) *node {
	if node.children == nil {
		return nil
	}
	return node.children[next]
}

func newTree() *tree {
	return &tree{root: &node{children: make(map[rune]*node)}}
}

func (tree *tree) Add(word string, wordType WordType) {
	word = normalizeAsNFC(word)
	wordFlag := wordFlagDefault
	for strings.HasPrefix(word, "*") {
		word = strings.TrimPrefix(word, "*")
		wordFlag.SetRequireHeadSpace(false)
		tree.hasHeadingWildcard = true
	}
	for strings.HasSuffix(word, "*") {
		word = strings.TrimSuffix(word, "*")
		wordFlag.SetRequireTailSpace(false)
	}

	for _, w := range buildWordListHandleWildcard(word) {
		tree.add(w, wordType, wordFlag)
	}
}

func (tree *tree) add(word string, wordType WordType, flag WordFlag) {
	if len(word) == 0 {
		return
	}
	current := tree.root
	for _, ch := range word {
		next := current.Next(ch)
		if next == nil {
			next = &node{}
		}
		if current.children == nil {
			current.children = make(map[rune]*node)
		}
		current.children[ch] = next
		current = next
	}
	if current.word == nil {
		current.word = &wordData{wordFlag: wordFlagDefault}
	}
	current.word.word = word
	if current.word.wordType < wordType {
		current.word.wordType = wordType
	}
	current.word.wordFlag = flag
}

// buildWordListHandleWildcard if input contains wildcard "*", this returns all combinations
// which can be possible for matching.
//
// For example: xx*yy*zz -> []string{xx*yy*zz, xxyy*zz, xx*yyzz, xxyyzz}
func buildWordListHandleWildcard(word string) (res []string) {
	p := strings.SplitN(word, "*", 2) //nolint:mnd
	if len(p) == 1 {
		res = append(res, word)
		return
	}
	for _, sub := range buildWordListHandleWildcard(p[1]) {
		res = append(res, p[0]+sub)
		res = append(res, p[0]+"*"+sub)
	}
	return
}
