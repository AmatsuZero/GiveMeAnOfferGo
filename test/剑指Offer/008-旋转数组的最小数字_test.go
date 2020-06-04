package 剑指Offer

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinNumberInRotatedArray(t *testing.T) {
	// 有序无重复
	source := []int{3, 7, 9, 1, 2, 3}
	assert.Equal(t, 剑指Offer.MinNumberInRotatedArray(source), 1)
	// 有序有重复
	source = []int{1, 0, 1, 1, 1}
	assert.Equal(t, 剑指Offer.MinNumberInRotatedArray(source), 0)
	// 升序排序数组
	source = []int{1, 2, 3, 4, 5}
	assert.Equal(t, 剑指Offer.MinNumberInRotatedArray(source), 1)
	// 只有一个元素
	source = []int{1}
	assert.Equal(t, 剑指Offer.MinNumberInRotatedArray(source), 1)
}
