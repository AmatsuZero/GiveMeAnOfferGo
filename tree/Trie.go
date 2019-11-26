package tree

type trieNode struct {
	children map[rune]*trieNode
	isEnd    bool
}

func newTrieNode() *trieNode {
	return &trieNode{
		children: make(map[rune]*trieNode, 26),
		isEnd:    false,
	}
}

type Trie struct {
	root *trieNode
}

func NewTrie() *Trie {
	return &Trie{root: newTrieNode()}
}

func (trie *Trie) Insert(word string) {
	node := trie.root
	for _, w := range word {
		_, ok := node.children[w]
		if !ok {
			node.children[w] = newTrieNode()
		}
		node = node.children[w]
	}
	node.isEnd = true
}

func (trie *Trie) Search(word string) bool {
	node := trie.root
	for _, w := range word {
		_, ok := node.children[w]
		if !ok {
			return false
		}
		node = node.children[w]
	}
	return node.isEnd
}

func (trie *Trie) StartWith(prefix string) bool {
	node := trie.root
	for _, w := range prefix {
		_, ok := node.children[w]
		if !ok {
			return false
		}
		node = node.children[w]
	}
	return true
}
