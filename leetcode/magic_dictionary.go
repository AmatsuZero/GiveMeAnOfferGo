package leetcode

import "GiveMeAnOffer/defines"

type MagicDictionary struct {
	root *defines.Trie
}

func (d *MagicDictionary) BuildDict(dict []string) {
	if d.root == nil {
		d.root = defines.NewTrie(dict)
	}
}

func (d *MagicDictionary) Search(word string) bool {
	chArr := []rune(word)
	var dfs func(root *defines.TrieNode, i, edit int) bool
	dfs = func(root *defines.TrieNode, i, edit int) bool {
		if root == nil {
			return false
		}
		if root.IsWordEnd && i == len(word) && edit == 1 {
			return true
		}
		if i < len(word) && edit <= 1 {
			found := false
			for j := 0; j < defines.AlphabetSize && !found; j++ {
				next := edit
				if int(chArr[i]-'a') == j {
					next += 1
				}
				found = dfs(root.Children[j], i+1, next)
			}
		}
		return false
	}
	return dfs(d.root.Root, 0, 0)
}
