package depth_first_search

import (
	dfs "github.com/AmatsuZero/GiveMeAnOfferGo/depth-first-search"
	"testing"
)

func TestMinimumLengthEncoding(t *testing.T) {
	if dfs.MinimumLengthEncoding([]string{"time", "me", "bell"}) != 10 {
		t.Fail()
	}
}
