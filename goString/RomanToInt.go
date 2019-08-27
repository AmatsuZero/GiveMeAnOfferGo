package goString

import "strings"

var romanMap map[string]int

func init() {
	romanMap = map[string]int{
		"I":  1,
		"V":  5,
		"X":  10,
		"L":  50,
		"C":  100,
		"D":  500,
		"M":  1000,
		"IV": 4,
		"IX": 9,
		"XL": 40,
		"XC": 90,
		"CD": 400,
		"CM": 900,
	}
}

func RomanToInt(s string) (ret int) {
	symbols := make([]string, 0)
	for i, str := range strings.Split(s, "") {
		if i == 0 {
			symbols = append(symbols, str)
			continue
		}
		last := symbols[len(symbols)-1]
		switch {
		case last == "I" && (str == "X" || str == "V"),
			last == "X" && (str == "L" || str == "C"),
			last == "C" && (str == "D" || str == "M"):
			// 移除上一个
			symbols = symbols[:len(symbols)-1]
			symbols = append(symbols, last+str)
		default:
			symbols = append(symbols, str)
		}
	}
	// 计算
	for _, symbol := range symbols {
		ret += romanMap[symbol]
	}
	return
}
