package 剑指Offer

import "math"

func Power(base float64, exponent int) float64 {
	/*
		当底数是零且指数是负数的时候，如果不做特殊处理，就会出现对0求倒数从而导致程序运行出错
	*/
	if equal(base, 0) && exponent < 0 {
		return math.NaN()
	}
	absExponent := uint(exponent)
	if exponent < 0 {
		absExponent = uint(-exponent)
	}
	result := powerWithUnsignedExponent(base, absExponent)
	if exponent < 0 {
		result = 1.0 / result
	}
	return result
}

func powerWithUnsignedExponent(base float64, exponent uint) float64 {
	if exponent == 0 {
		return 1
	}
	if exponent == 1 {
		return base
	}
	result := powerWithUnsignedExponent(base, exponent>>1) // 用右移代替除2
	result *= result
	if exponent&0x1 == 1 { // 判断是否是奇数，根据公式，奇数还需要乘以base
		result *= base
	}
	return result
}

func equal(num1, num2 float64) bool {
	return (num1-num2 > -0.0000001) && (num1-num2 < 0.0000001)
}
