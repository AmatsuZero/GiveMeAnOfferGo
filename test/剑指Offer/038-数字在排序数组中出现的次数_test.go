package 剑指Offer

import (
	"testing"

	"github.com/AmatsuZero/GiveMeAnOfferGo/剑指Offer"
	"github.com/stretchr/testify/assert"
)

func TestGetNumberOfK(t *testing.T) {
	arr := []int{1, 2, 2, 2, 3, 4, 5}
	assert.Equal(t, 3, 剑指Offer.GetNumberOfK(arr, 2))
	assert.Equal(t, 0, 剑指Offer.GetNumberOfK(arr, 6))      // 不存在
	assert.Equal(t, 1, 剑指Offer.GetNumberOfK(arr, 3))      // 只出现一次
	assert.Equal(t, 1, 剑指Offer.GetNumberOfK(arr, 5))      // 最大值
	assert.Equal(t, 1, 剑指Offer.GetNumberOfK(arr, 1))      // 最小值
	assert.Equal(t, 1, 剑指Offer.GetNumberOfK([]int{1}, 1)) // 只有一个元素
	assert.Equal(t, 0, 剑指Offer.GetNumberOfK([]int{}, 1))  // 空数组
}
