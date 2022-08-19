package number_of_students_doing_homework_at_a_given_time

// BusyStudent https://leetcode.cn/problems/number-of-students-doing-homework-at-a-given-time/
func BusyStudent(startTime []int, endTime []int, queryTime int) int {
	clock := map[int]int{}
	for i := 0; i < len(startTime); i++ {
		for j := startTime[i]; j <= endTime[i]; j++ {
			clock[j] += 1
		}
	}
	return clock[queryTime]
}
