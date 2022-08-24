package copy_random_list

type Node struct {
	Val    int
	Next   *Node
	Random *Node
}

func CopyRandomList(head *Node) (ans *Node) {
	dict := map[*Node]*Node{} // 新旧节点
	replaceLater := map[*Node][]*Node{}
	cur := ans
	for n := head; n != nil; n = n.Next {
		newNode := &Node{Val: n.Val}
		if ans == nil {
			ans = newNode
			cur = ans
		} else {
			cur.Next = newNode
			cur = newNode
		}
		dict[n] = newNode
		if x, ok := replaceLater[n]; ok {
			for _, y := range x {
				y.Random = newNode
			}
		}
		if n.Random != nil {
			if x, ok := dict[n.Random]; ok {
				newNode.Random = x
			} else {
				replaceLater[n.Random] = append(replaceLater[n.Random], newNode)
			}
		}
	}
	return
}
