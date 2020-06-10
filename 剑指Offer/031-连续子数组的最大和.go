package 剑指Offer

/*
题目：输入一个整型数组，数组里有正数也有负数。数组中一个或连续的多个整数组成一个子数组。求所有子数组的和的最大值。要求时间复杂度为O（n）
*/
func FindGreatestSumOfSubArray(pData []int) int {
	if len(pData) == 0 {
		return 0
	}
	nCurSum, nGreatestSum := 0, MinInt
	for i := 0; i < len(pData); i++ {
		if nCurSum <= 0 {
			nCurSum = pData[i]
		} else {
			nCurSum += pData[i]
		}
		if nCurSum > nGreatestSum {
			nGreatestSum = nCurSum
		}
	}
	return nGreatestSum
}
