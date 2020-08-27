package datastructure

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/array"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

type stack []float64
type queue []float64

func (s *stack) push(n float64) {
	*s = append(*s, n)
}

func (s *stack) pop(n float64) float64 {
	l := len(*s)
	if l == 0 {
		return math.NaN()
	}
	last := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return last
}

func TestMergeIntervals(t *testing.T) {
	assert.Equal(t, [][]int{{1, 6}, {8, 10}, {15, 18}}, array.Merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}))
	source := [][]int{
		{5, 100}, {3, 9}, {12, 78}, {45, 63}, {23, 27}, {1, 120}, {30, 53}, {26, 86}, {92, 97}, {27, 89},
	}

	t.Log(array.Merge(source))
}

func op(n, target int) int {
	step := 0
	for n != target {
		if n%3 == 0 {
			n /= 3
		} else if n%2 == 0 {
			n /= 2
		} else {
			n -= 1
		}
		step += 1
	}
	return step
}

func dictStr(s string) string {
	a, b := 0, len(s)-1
	newStr := make([]uint8, 0)
	for a <= b {
		left := false
		for i := 0; a-1 <= b; i++ {
			if s[a+i] < s[b-i] {
				left = true
				break
			} else if s[a+i] > s[b-i] {
				left = false
				break
			}
		}

		if left {
			newStr = append(newStr, s[a])
			a++
		} else {
			newStr = append(newStr, s[b])
			b--
		}
	}
	return string(newStr)
}

func TestFuck(t *testing.T) {

}
