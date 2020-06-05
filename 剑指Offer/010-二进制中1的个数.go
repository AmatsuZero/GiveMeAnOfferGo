package 剑指Offer

func NumberOfOne(n int) int {
	count := 0
	/*
		把一个整数减去1之后再和原来的整数做位与运算，得到的结果相当于是把整数的二进制表示中的最右边一个1变成0。
	*/
	for n != 0 {
		count++
		n = (n - 1) & n
	}
	return count
}
