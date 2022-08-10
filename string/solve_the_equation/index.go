package solve_the_equation

import (
	"strconv"
	"unicode"
)

// SolveEquation https://leetcode.cn/problems/solve-the-equation/
func SolveEquation(equation string) string {
	factor, val := 0, 0
	i, n, sign := 0, len(equation), 1 // 等式左边符号默认为正
	for i < n {
		if equation[i] == '=' {
			sign = -1 // 等式右边符号默认为负
			i += 1
			continue
		}

		s := sign
		if equation[i] == '+' { // 去掉前边的符号
			i += 1
		} else if equation[i] == '-' {
			s = -s
			i += 1
		}

		num, valid := 0, false
		if i < n && unicode.IsDigit(rune(equation[i])) {
			valid = true
			num = num*10 + int(equation[i]-'0')
			i += 1
		}

		if i < n && equation[i] == 'x' { // 变量
			if valid {
				s *= num
			}
			factor += s
			i += 1
		} else { // 数值
			val += s * num
		}
	}
	if factor == 0 {
		if val == 0 {
			return "Infinite solutions"
		}
		return "No solution"
	}
	return "x=" + strconv.Itoa(-val/factor)
}
