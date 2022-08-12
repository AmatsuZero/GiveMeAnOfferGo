package defines

const AlphabetSize = 26

type TrieNode struct {
	Children  [AlphabetSize]*TrieNode
	IsWordEnd bool
}

type Trie struct {
	Root *TrieNode
}

func NewTrie(words []string) *Trie {
	root := &Trie{&TrieNode{}}
	if len(words) == 0 {
		return root
	}
	for _, word := range words {
		root.Insert(word)
	}
	return root
}

func (t *Trie) Insert(word string) {
	wordLen, cur := len(word), t.Root
	for i := 0; i < wordLen; i++ {
		index := word[i] - 'a'
		if cur.Children[index] == nil {
			cur.Children[index] = &TrieNode{}
		}
		cur = cur.Children[index]
	}
	cur.IsWordEnd = true
}

func (t *Trie) Find(word string) bool {
	wordLen, cur := len(word), t.Root
	for i := 0; i < wordLen; i++ {
		index := word[i] - 'a'
		if cur.Children[index] == nil {
			return false
		}
		cur = cur.Children[index]
	}
	return cur.IsWordEnd
}
