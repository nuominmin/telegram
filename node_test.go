package telegram

import (
	"testing"
)

func TestNode(t *testing.T) {
	// 创建树并添加节点
	tree := NewTree()

	tree.AddNode(&TreeNode{
		Name: "nodeA",
		Children: []*TreeNode{
			{Name: "nodeA1"},
			{Name: "nodeA2", Children: []*TreeNode{
				{Name: "nodeA21"},
			}},
		},
	})

	tree.AddNode(&TreeNode{
		Name: "nodeB",
		Children: []*TreeNode{
			{Name: "nodeB1"},
		},
	})

	tree.AddNode(&TreeNode{
		Name: "nodeC",
		Children: []*TreeNode{
			{Name: "nodeC1"},
		},
	})

	tree.PrintTree()
}
