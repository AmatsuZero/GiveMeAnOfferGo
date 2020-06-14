package 剑指Offer

/*
输入一个递增排序的数组和一个数字s，在数组中查找两个数，使得它们的和正好是s。如果有多对数字的和等于s，输出任意一对即可
*/
func FindNumberWithSum(data []int, sum int) (found bool, num1, num2 int) {
	if len(data) == 1 {
		return
	}
	ahead, behind := len(data)-1, 0
	for ahead > behind {
		curSum := data[ahead] + data[behind]
		if curSum == sum {
			found = true
			num1 = data[behind]
			num2 = data[ahead]
			return
		} else if curSum > sum {
			ahead--
		} else {
			behind++
		}
	}
	return
}

/*
题目二：输入一个正数s，打印出所有和为s的连续正数序列（至少含有两个数）。例如输入15，由于1＋2＋3＋4＋5＝4＋5＋6＝7＋8＝15，所以结果打印出3个连续序列1～5、4～6和7～8
*/
func FindContinuousSequence(sum int) (output [][]int) {
	if sum < 3 {
		return
	}

	small, big := 1, 2
	middle := (1 + sum) / 2
	curSum := small + big

	for small < middle {
		if curSum == sum {
			output = append(output, generateSequence(small, big))
		}
		for curSum > sum && small < middle {
			curSum -= small
			small++
			if curSum == sum {
				output = append(output, generateSequence(small, big))
			}
		}
		big++
		curSum += big
	}
	return
}

func generateSequence(small, big int) []int {
	arr := make([]int, big-small+1)
	for i := range arr {
		arr[i] = small + i
	}
	return arr
}
