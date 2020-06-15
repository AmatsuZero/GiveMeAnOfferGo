package 剑指Offer

import (
	"fmt"
)

func Atoi(input string) (err error, output int) {
	if len(input) == 0 {
		return fmt.Errorf("字符串长度为0"), output
	}
	chars := []byte(input)
	minus := false
	index := 0
	if chars[index] == '+' {
		index++
	} else if chars[index] == '-' {
		minus = true
		index++
	}
	if index != len(chars)-1 {
		e, num := strToIntCore(chars[index:], minus)
		err, output = e, int(num)
	}
	return err, output
}

func strToIntCore(digits []byte, minus bool) (err error, output int64) {
	for _, digit := range digits {
		if digit >= '0' && digit <= '9' {
			flag := int64(1)
			if minus {
				flag = -1
			}
			output = output*10 + flag*int64(digit-'0')
			if (!minus && output > int64(MaxInt)) || (minus && output < int64(MinInt)) {
				output = 0
				err = fmt.Errorf("数据溢出")
				break
			}
		} else {
			output = 0
			err = fmt.Errorf("包含非数字字符")
			break
		}
	}
	return
}
