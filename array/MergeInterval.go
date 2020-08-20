package array

import "sort"

/*
给出一个区间的集合，请合并所有重叠的区间。



示例 1:

输入: intervals = [[1,3],[2,6],[8,10],[15,18]]
输出: [[1,6],[8,10],[15,18]]
解释: 区间 [1,3] 和 [2,6] 重叠, 将它们合并为 [1,6].
示例 2:

输入: intervals = [[1,4],[4,5]]
输出: [[1,5]]
解释: 区间 [1,4] 和 [4,5] 可被视为重叠区间。
注意：输入类型已于2019年4月15日更改。 请重置默认代码定义以获取新方法签名。



提示：

intervals[i][0] <= intervals[i][1]

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/merge-intervals
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/

func Merge(intervals [][]int) (merged [][]int) {
	if len(intervals) == 0 {
		return
	}
	// 先排序，确定可以合并的范围
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] <= intervals[j][0]
	})
	for _, interval := range intervals {
		// 如果结果为空或者当前数组的头小于结果最后元素的尾,说明两个数组没有相交, 则提交到结果
		if len(merged) == 0 || interval[0] > merged[len(merged)-1][1] {
			merged = append(merged, interval)
			// 	否则如果当前数组尾大于结果区最后一个元素的尾, 则说明两个数组相交, 则更新结果区最后一个元素的尾
		} else if interval[1] > merged[len(merged)-1][1] {
			merged[len(merged)-1][1] = interval[1]
		}
	}
	return
}
