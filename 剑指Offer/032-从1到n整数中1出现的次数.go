package 剑指Offer

/*
题目：输入一个整数n，求从1到n这n个整数的十进制表示中1出现的次数。例如输入12，从1到12这些整数中包含1 的数字有1，10，11和12，1一共出现了5次。
*/
func Ones(n int) int {
	/*
		f(n) = n1*f(10bit-1) + f(n – n1*10bit) + LEFT;
		其中
		if(n1 == 1)
			LEFT = n - 10bit+ 1;
		else
			LEFT = 10bit;
	*/
	if n == 0 {
		return 0
	}
	if n > 1 && n < 10 {
		return 1
	}
	count := 0
	highest := n
	bit := 0
	for highest >= 10 {
		highest /= 10
		bit++
	}
	weight := pow(10, bit)
	if highest == 1 {
		count = Ones(weight-1) + Ones(n-weight) + (n - weight + 1)
	} else {
		count = Ones(weight-1) + Ones(n-highest*weight) + weight
	}
	return count
}

func pow(a, b int) int {
	res := 1
	for i := b; i > 0; i-- {
		res *= a
	}
	return res
}
