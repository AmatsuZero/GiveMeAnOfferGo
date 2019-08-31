package string

import "strconv"

func CountAndSay(n int) string {
	if n == 1 {
		return "1"
	}
	ans := "11"
	for n > 2 {
		n--
		tmp := ""
		now := '+'
		count := 0
		for _, v := range ans {
			if now == '+' {
				now = v
				count++
				continue
			}
			if now == v {
				count++
				continue
			} else {
				//int to string
				tmp = tmp + strconv.Itoa(count)
				tmp = tmp + string(now)
				count = 1
				now = v
				continue
			}
		}
		tmp = tmp + strconv.Itoa(count)
		tmp = tmp + string(now)
		ans = tmp
	}
	return ans
}
