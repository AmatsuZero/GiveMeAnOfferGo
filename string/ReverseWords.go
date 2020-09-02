package string

/*
给定一个字符串，逐个翻转字符串中的每个单词。



示例 1：

输入: "the sky is blue"
输出: "blue is sky the"
示例 2：

输入: " hello world! "
输出:"world! hello"
解释: 输入字符串可以在前面或者后面包含多余的空格，但是反转后的字符不能包括。
示例 3：

输入: "a good  example"
输出: "example good a"
解释: 如果两个单词间有多余的空格，将反转后单词间的空格减少到只含一个。


说明：

无空格字符构成一个单词。
输入字符串可以在前面或者后面包含多余的空格，但是反转后的字符不能包括。
如果两个单词间有多余的空格，将反转后单词间的空格减少到只含一个。


进阶：

请选用 C 语言的用户尝试使用O(1) 额外空间复杂度的原地解法。

作者：力扣 (LeetCode)
链接：https://leetcode-cn.com/leetbook/read/array-and-string/crmp5/
来源：力扣（LeetCode）
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
*/
func ReverseWords(s string) string {
	strBuffer := []rune(s)
	// 全局反转
	for left, right := 0, len(strBuffer)-1; left < right; left, right = left+1, right-1 {
		strBuffer[left], strBuffer[right] = strBuffer[right], strBuffer[left]
	}
	// 反转单词
	for i := 0; i < len(strBuffer); i++ {
		end := i
		for ; strBuffer[end] != ' ' && end < len(strBuffer); end++ {

		}
		for start := i; start < end; start++ {

		}
	}
	return string(strBuffer)
}
