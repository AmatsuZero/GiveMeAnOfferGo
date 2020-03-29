package tree

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/Collections/Tree"
	"github.com/AmatsuZero/GiveMeAnOfferGo/test/Utils"
	"testing"
)

func TestContains(t *testing.T) {
	getString := Utils.GetString
	trie := Tree.NewTrie()
	trie.Insert(getString("cute"))
	if trie.Contains(getString("cute")) {
		t.Log("cute is in the trie")
	}
}

func TestRemove(t *testing.T) {
	getString := Utils.GetString
	trie := Tree.NewTrie()
	trie.Insert(getString("cut"))
	trie.Insert(getString("cute"))

	t.Log("\n*** Before removing ***")
	if !trie.Contains(getString("cut")) {
		t.Fail()
	}
	t.Log("\n\"cut\" is in the trie")
	if !trie.Contains(getString("cute")) {
		t.Fail()
	}
	t.Log("\n\"cute\" is in the trie")

	t.Log("\n*** After removing cut***")
	trie.Remove(getString("cut"))
	if !trie.Contains(getString("cut")) &&
		trie.Contains(getString("cute")) {
		t.Log("\n\"cute\" is still in the trie")
	}

}
