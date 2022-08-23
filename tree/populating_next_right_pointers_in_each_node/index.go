package populating_next_right_pointers_in_each_node

type Node struct {
	Val   int
	Left  *Node
	Right *Node
	Next  *Node
}

// Connect https://leetcode.cn/problems/populating-next-right-pointers-in-each-node/
func Connect(root *Node) *Node {
	var dfs func(lv int, node *Node)
	layers := map[int][]*Node{}

	dfs = func(lv int, node *Node) {
		if node == nil {
			return
		}
		layers[lv] = append(layers[lv], node)
		dfs(lv+1, node.Left)
		dfs(lv+1, node.Right)
	}
	dfs(0, root)
	for _, nodes := range layers {
		var tmp *Node
		for _, node := range nodes {
			if tmp != nil {
				tmp.Next = node
			}
			tmp = node
		}
	}
	return root
}
