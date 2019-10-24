package dynamicprogramming

func ClimbStairs(n int) int {
	if n < 1 {
		return 0
	} else if n == 1 {
		return 1
	} else if n == 2 {
		return 2
	}
	a := 1
	b := 2
	temp := a + b
	for i := 3; i <= n; i++ {
		temp = a + b
		a = b
		b = temp
	}
	return temp
}
