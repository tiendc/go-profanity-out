package profanityout

import "strings"

type tree struct {
	root               *node
	hasHeadingWildcard bool
}

type node struct {
	children         map[rune]*node
	word             string
	wordType         WordType
	requireHeadSpace bool
	requireTailSpace bool
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
	requireHeadSpace := true
	requireTailSpace := true
	for strings.HasPrefix(word, "*") {
		word = strings.TrimPrefix(word, "*")
		requireHeadSpace = false
		tree.hasHeadingWildcard = true
	}
	for strings.HasSuffix(word, "*") {
		word = strings.TrimSuffix(word, "*")
		requireTailSpace = false
	}
	tree.add(word, wordType, requireHeadSpace, requireTailSpace)

	for _, w := range buildWordListHandleWildcard(word) {
		if w != word {
			tree.add(w, wordType, true, true)
		}
	}
}

func (tree *tree) add(word string, wordType WordType, requireHeadSpace, requireTailSpace bool) {
	if len(word) == 0 {
		return
	}
	current := tree.root
	for _, ch := range word {
		next := current.Next(ch)
		if next == nil {
			next = &node{requireHeadSpace: true, requireTailSpace: true}
		}
		if current.children == nil {
			current.children = make(map[rune]*node)
		}
		current.children[ch] = next
		current = next
	}
	current.word = word
	if current.requireHeadSpace {
		current.requireHeadSpace = requireHeadSpace
	}
	if current.requireTailSpace {
		current.requireTailSpace = requireTailSpace
	}
	if current.wordType < wordType {
		current.wordType = wordType
	}
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
