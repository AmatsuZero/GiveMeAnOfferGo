package perfect_squares

import "container/list"

// NumSquaries https://leetcode.cn/problems/perfect-squares/
func NumSquares(n int) (ans int) {
	queue := list.New()
	queue.PushBack(0)

	visited := map[int]bool{}
	visited[0] = true

	for queue.Len() > 0 {
		ans += 1
		for i := queue.Len(); i > 0; i-- {
			node := queue.Remove(queue.Front()).(int)
			for j := 0; j < n; j++ {
				num := node + j*j
				if num == n {
					return
				}
				if num > n {
					break
				}
				if !visited[num] {
					visited[num] = true
					queue.PushBack(num)
				}
			}
		}
	}

	return
}
