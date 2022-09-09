package crawler_log_folder

// MinOperations https://leetcode.cn/problems/crawler-log-folder/
func MinOperations(logs []string) int {
	ans := 0
	for _, log := range logs {
		switch log {
		case "./":
			continue
		case "../":
			if ans == 0 {
				continue
			}
			ans -= 1
		default:
			ans += 1
		}
	}
	return ans
}
