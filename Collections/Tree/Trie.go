package Tree

import "github.com/AmatsuZero/GiveMeAnOfferGo/Objects"

type TrieNode struct {
	Key           Objects.Hashable
	Parent        *TrieNode
	Children      map[uint32]*TrieNode
	IsTerminating bool
}

func NewTrieNode(key Objects.Hashable, parent *TrieNode) *TrieNode {
	return &TrieNode{
		Key:           key,
		Parent:        parent,
		Children:      map[uint32]*TrieNode{},
		IsTerminating: false,
	}
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{root: NewTrieNode(nil, nil)}
}

func (trie *Trie) Insert(collection ...Objects.Hashable) {
	current := trie.root
	for _, element := range collection {
		if _, ok := current.Children[element.HashValue()]; !ok {
			current.Children[element.HashValue()] = NewTrieNode(element, current)
		}
		current = current.Children[element.HashValue()]
	}
	current.IsTerminating = true
}

func (trie *Trie) Contains(collection ...Objects.Hashable) bool {
	current := trie.root
	for _, element := range collection {
		child, ok := current.Children[element.HashValue()]
		if !ok {
			return false
		}
		current = child
	}
	return current.IsTerminating
}

func (trie *Trie) Remove(collection ...Objects.Hashable) {
	current := trie.root
	for _, element := range collection {
		child, ok := current.Children[element.HashValue()]
		if !ok {
			return
		}
		current = child
	}
	if !current.IsTerminating {
		return
	}
	current.IsTerminating = false

	for parent := current.Parent; parent != nil &&
		len(current.Children) == 0 &&
		!current.IsTerminating; {
		delete(parent.Children, current.Key.HashValue())
		current = parent
	}
}
