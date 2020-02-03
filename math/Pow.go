package math

func MyPow(x float64, n int) float64 {
	N := int64(n)
	if N < 0 {
		x = 1 / x
		N = -N
	}

	ans := float64(1)
	currentProduct := x
	for i := N; i != 0; i /= 2 {
		if (i % 2) == 1 {
			ans = ans * currentProduct
		}
		currentProduct *= currentProduct
	}
	return ans
}
