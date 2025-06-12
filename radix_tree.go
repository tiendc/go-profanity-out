package profanityout

type tree struct {
	root *node
}

type node struct {
	children map[rune]*node
	word     string
	wordType WordType
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
	current.word = word
	if current.wordType < wordType {
		current.wordType = wordType
	}
}
