package 剑指Offer

import "math"

const gMaxValue = 6
/*
题目：把n个骰子扔在地上，所有骰子朝上一面的点数之和为s。输入n，打印出s的所有可能的值出现的概率。
*/
func ProbabilityOfDiceNumbers(num int) (result map[int]float64) {
	if num < 1 {
		return
	}
	result = map[int]float64{}
	probabilities := [2][]int{}
	probabilities[0] = make([]int, gMaxValue * num + 1)
	probabilities[1] = make([]int, gMaxValue * num + 1)
	flag := 0
	for i := 1; i <= gMaxValue; i++ {
		probabilities[flag][i] = 1
	}
	for k := 2; k <= num; k++ {
		for i := 0; i < k; i++ {
			probabilities[1-flag][i] = 0
		}
		for i := k; i <= gMaxValue * k; i++ {
			probabilities[1-flag][i] = 0
			for j := 1; j <= i && j <= gMaxValue; j++ {
				probabilities[1-flag][i] += probabilities[flag][i-j]
			}
		}
		flag = 1 - flag
	}
	total := math.Pow(gMaxValue, float64(num))
	for i := num; i <= gMaxValue * num; i++ {
		result[i] = float64(probabilities[flag][i]) / total
	}
	return
}