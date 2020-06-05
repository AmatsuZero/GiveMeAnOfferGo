package 剑指Offer

import "math"

/*
题目一：写一个函数，输入n，求斐波那契（Fibonacci）数列的第n项
*/

func Fibonacci(n uint) int64 {
	result := []int64{0, 1}
	if n < 2 {
		return result[n]
	}
	/*
		首先根据f（0）和f（1）算出f（2），再根据f（1）和f（2）算出f（3）……依此类推就可以算出第n项了。很容易理解，这种思路的时间复杂度是O（n）
	*/
	fibNMinusOne := int64(1)
	fibNMinusTwo := int64(0)
	fibN := int64(0)
	for i := uint(2); i <= n; i++ {
		fibN = fibNMinusOne + fibNMinusTwo
		fibNMinusTwo, fibNMinusOne = fibNMinusOne, fibN
	}
	return fibN
}

/*
题目二：一只青蛙一次可以跳上1级台阶，也可以跳上2级。求该青蛙跳上一个n级的台阶总共有多少种跳法
*/
func FrogSteps(n uint) int64 {
	result := []int64{0, 1}
	if n < 2 {
		return result[n]
	}
	fibNMinusOne := int64(1) // 按照第一次一个台阶
	fibNMinusTwo := int64(1) // 按照第一次两个台阶
	fibN := int64(0)
	for i := uint(2); i <= n; i++ {
		fibN = fibNMinusOne + fibNMinusTwo
		fibNMinusTwo, fibNMinusOne = fibNMinusOne, fibN
	}
	return fibN
}

/*
题目二：一只青蛙一次可以跳上1级台阶，也可以跳上2级, 还可以跳上N级。求该青蛙跳上一个n级的台阶总共有多少种跳法
*/
func FrogStepsN(n uint) int64 {
	return int64(math.Pow(2, float64(n-1)))
}
