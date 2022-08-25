package max_alive_year

// MaxAliveYear https://leetcode.cn/problems/living-people-lcci/
func MaxAliveYear(birth []int, death []int) (ans int) {
	data := make([]int, 102)
	for _, x := range birth {
		data[x-1900] += 1
	}
	for _, x := range death {
		data[x-1899] -= 1
	}

	for i := 1; i < len(data); i++ {
		data[i] += data[i-1]
		if data[i] > data[ans] {
			ans = i
		}
	}
	return 1900 + ans
}
