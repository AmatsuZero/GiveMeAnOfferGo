package validate_stack_sequences

// ValidateStackSequences https: //leetcode.cn/problems/validate-stack-sequences/
func ValidateStackSequences(pushed []int, popped []int) bool {
	var stack []int
	j := 0
	for i := 0; i < len(pushed); i++ {
		stack = append(stack, pushed[i])
		for ; len(stack) > 0 && stack[len(stack)-1] == popped[j]; j++ { // 遇到需要 popped，出栈
			stack = stack[:len(stack)-1]
		}
	}

	if len(popped)-j != len(stack) {
		return false
	}

	for ; j < len(popped); j++ {
		if stack[len(stack)-1] != popped[j] {
			return false
		}
		stack = stack[:len(stack)-1]
	}

	return true
}
