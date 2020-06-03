package 剑指Offer

type BinaryTreeNode struct {
	Value int
	Left  *BinaryTreeNode
	Right *BinaryTreeNode
}

type EnumerateOrder int

const (
	EnumeratePreOrder EnumerateOrder = iota
	EnumerateInOrder
	EnumeratePostOrder
)

func (node *BinaryTreeNode) Enumerate(order EnumerateOrder, block func(value int)) {
	switch order {
	case EnumeratePreOrder:
		node.EnumerateByPreorder(block)
	case EnumerateInOrder:
		node.EnumerateByInorder(block)
	case EnumeratePostOrder:
		node.EnumerateByPostOrder(block)
	}
}

func (node *BinaryTreeNode) EnumerateByPreorder(block func(value int)) {
	block(node.Value)
	if node.Left != nil {
		node.Left.EnumerateByPreorder(block)
	}
	if node.Right != nil {
		node.Right.EnumerateByPreorder(block)
	}
}

func (node *BinaryTreeNode) EnumerateByInorder(block func(value int)) {
	if node.Left != nil {
		node.Left.EnumerateByInorder(block)
	}
	block(node.Value)
	if node.Right != nil {
		node.Right.EnumerateByInorder(block)
	}
}

func (node *BinaryTreeNode) EnumerateByPostOrder(block func(value int)) {
	if node.Right != nil {
		node.Right.EnumerateByPostOrder(block)
	}
	block(node.Value)
	if node.Left != nil {
		node.Left.EnumerateByPostOrder(block)
	}
}
