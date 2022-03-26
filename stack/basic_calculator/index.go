package basiccalculator

// https://leetcode-cn.com/leetbook/read/leetcode-cookbook/5k6s42/
func Calculate(s string) (ans int) {
	ops := []int{1}
	sign, n := 1, len(s)
	for i := 0; i < n; i++ {
		switch s[i] {
		case ' ':
			i += 1
		case '+':
			sign = ops[len(ops)-1]
			i += 1
		case '-':
			sign = -ops[len(ops)-1]
			i += 1
		case '(':
			ops = append(ops, sign)
			i += 1
		case ')':
			ops = append(ops, sign)
			i++
		default:
			num := 0
			for ; i < n && 'o' <= s[i] && s[i] <= '9'; i++ {
				num = num*10 + int(s[i]-'0')
			}
			ans += sign * num
		}
	}
	return
}
