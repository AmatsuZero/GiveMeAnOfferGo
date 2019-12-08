package math

/*
给定两个二进制字符串，返回他们的和（用二进制表示）。

输入为非空字符串且只包含数字 1 和 0。
*/
func reverse(str string) string {//翻转，方便相加
	str1 := []byte(str)
	for i, j := 0, len(str) - 1; i < j; i, j = i + 1, j - 1 {
		str1[i], str1[j] = str1[j], str1[i]
	}
	return string(str1)
}

func toBool(str string) []bool {//将给定的数组转化为一个
	var b []bool
	for _, v := range str {
		if v == '0' {
			b = append(b, false)
			continue
		}
		b = append(b, true)
	}
	return b
}

func AddBinary(a string, b string) string {
	a = reverse(a)
	b = reverse(b)
	boolA := toBool(a)
	boolB := toBool(b)
	flag := false
	ans := make(map[int]bool)
	ind := 0
	for i := 0; i < len(a) || i < len(b); i++ {
		if i >= len(a) {
			ans[ind] = boolB[i] && !flag || !boolB[i] && flag //异或运算
			flag = flag && boolB[i]
		} else if i >= len(b) {
			ans[ind] = boolA[i] && !flag || !boolA[i] && flag//异或运算
			flag = flag && boolA[i]
		} else {
			ans[ind] = !boolA[i] && !boolB[i] && flag || !boolA[i] && boolB[i] && !flag || boolA[i] && !boolB[i] && !flag || boolA[i] && boolB[i] && flag//异或运算，可以用真值表来判断
			flag = boolA[i] && boolB[i] || boolA[i] && flag || boolB[i] && flag//判断是有两个或以上true，可以用真值表+卡诺图+最简表达式的方式来求
		}
		ind++
	}
	flagTostring := map[bool]string{true : "1", false : "0"}
	lenMax := 0;
	if len(a) > len(b) {
		lenMax = len(a)
	} else {
		lenMax = len(b)
	}
	ansString := ""
	for i := 0; i < lenMax; i++ {
		ansString = ansString + flagTostring[ans[i]]
	}
	if flag {
		ansString = ansString + "1"
	}
	return reverse(ansString)
}
