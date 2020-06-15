package 剑指Offer

import "unsafe"

/*
题目：求1＋2＋…＋n，要求不能使用乘除法、for、while、if、else、switch、case等关键字及条件判断语句（A?B:C）
*/
func SumSolutionRecursive(n int) (sum int) {
	shouldTerminate := func(num int) int {
		return 0
	}
	arr := []func(n int) int {shouldTerminate, SumSolutionRecursive}
	boolean := func(n int) int {
		bValue := !!(n != 0)
		return *(*int)(unsafe.Pointer(&bValue))&1
	}
	return n + arr[boolean(n)](n-1)
}