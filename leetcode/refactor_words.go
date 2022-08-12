package leetcode

import (
	"GiveMeAnOffer/defines"
	"strings"
)

func ReplaceWords(dict []string, sentence string) string {
	root := buildTree(dict)
	words := strings.Split(sentence, " ")
	for i := 0; i < len(words); i++ {
		prefix := FindPrefix(root, words[i])
		if len(prefix) > 0 {
			words[i] = prefix
		}
	}
	return strings.Join(words, " ")
}

func buildTree(dict []string) *defines.Trie {
	return defines.NewTrie(dict)
}

func FindPrefix(root *defines.Trie, word string) string {
	prefix := strings.Builder{}
	for _, ch := range word {
		prefix.WriteRune(ch)
		str := prefix.String()
		if root.Find(str) {
			return str
		}
	}
	return ""
}
