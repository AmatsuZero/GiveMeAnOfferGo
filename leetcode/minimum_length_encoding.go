package leetcode

import "GiveMeAnOffer/defines"

func MiniMumLengthEncoding(words []string) int {
	tree := &defines.TrieNode{}
	build := func() {
		node := tree
		for _, word := range words {
			arr := []rune(word)
			for i := len(arr) - 1; i >= 0; i-- {
				ch := arr[i]
				if node.Children[ch-'a'] == nil {
					node.Children[ch-'a'] = &defines.TrieNode{}
				}
				node = node.Children[ch-'a']
			}
		}
	}
	build()

	total := 0
	var dfs func(root *defines.TrieNode, len int)
	dfs = func(root *defines.TrieNode, len int) {
		isLeaf := true
		for _, child := range root.Children {
			if child != nil {
				isLeaf = false
				dfs(child, len+1)
			}
		}

		if isLeaf {
			total += len
		}
	}
	dfs(tree, 0)
	return total
}
