package 剑指Offer

import "sort"

/*
题目：从扑克牌中随机抽5张牌，判断是不是一个顺子，即这5张牌是不是连续的。2～10为数字本身，A为1，J为11，Q为12，K为13，而大、小王可以看成任意数字
*/
func IsStraightCards(cards []int) bool {
	if cards == nil || len(cards) < 1{
		return false
	}
	// 现排序，便于计算间隔
	sort.Slice(cards, func(i, j int) bool {
		return cards[i] < cards[j]
	})
	numberOfZero, numberOfGap := 0, 0
	// 统计数组中0的个数 （Joker）
	for i := 0; i < len(cards) && cards[i] == 0; i++ {
		numberOfZero++
	}
	// 统计数组中间隔的个数
	small := numberOfZero // 从不是Joker的卡开始找
	big := small + 1
	for big < len(cards) {
		// 两个数相等，有顺子，不可能是顺子
		if cards[small] == cards[big] {
			return false
		}
		numberOfGap += cards[big] - cards[small] - 1 // 如果连着的话，两张卡相差应该为1
		small = big
		big++
	}
	return numberOfGap <= numberOfZero
}