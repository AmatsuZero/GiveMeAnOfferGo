package leetcode

import "GiveMeAnOffer/defines"

type mapSumNode struct {
	children [defines.AlphabetSize]*mapSumNode
	value    int
}

type MapSum struct {
	root *mapSumNode
}

func (s *MapSum) Insert(key string, val int) {
	if s.root == nil {
		s.root = &mapSumNode{}
	}
	node := s.root
	for _, ch := range key {
		if node.children[ch-'a'] == nil {
			node.children[ch-'a'] = &mapSumNode{}
		}
		node = node.children[ch-'a']
	}
	node.value = val
}

func (s *MapSum) Sum(prefix string) int {
	var getSum func(node *mapSumNode) int
	getSum = func(node *mapSumNode) int {
		if node == nil {
			return 0
		}
		result := node.value
		for _, child := range node.children {
			result += getSum(child)
		}
		return result
	}

	node := s.root
	for _, ch := range prefix {
		if node.children[ch-'a'] == nil {
			return 0
		}
		node = node.children[ch-'a']
	}
	return getSum(node)
}
