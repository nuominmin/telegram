package telegram

import (
	"fmt"
)

// TreeNode 代表树中的一个节点
type TreeNode struct {
	Name     string
	Children []*TreeNode
}

// Tree 代表一个多根节点的树
type Tree struct {
	RootNodes []*TreeNode
}

func NewTree() *Tree {
	return &Tree{}
}

// AddNode 向树中添加一个节点，并递归添加其子节点
func (t *Tree) AddNode(node *TreeNode) (*TreeNode, error) {
	if t.FindNode(node.Name) != nil {
		return nil, fmt.Errorf("Error: Node with name '%s' already exists\n", node.Name)
	}

	// 创建根节点并递归添加子节点
	rootNode := &TreeNode{Name: node.Name}
	t.addChildren(rootNode, node.Children)
	t.RootNodes = append(t.RootNodes, rootNode)
	return rootNode, nil
}

// addChildren 递归地为一个节点添加子节点
func (t *Tree) addChildren(parent *TreeNode, children []*TreeNode) {
	for _, n := range children {
		childNode := &TreeNode{Name: n.Name}
		t.addChildren(childNode, n.Children) // 递归添加子节点
		parent.Children = append(parent.Children, childNode)
	}
}

// FindNode 在整个树中递归查找名字为 name 的节点
func (t *Tree) FindNode(name string) *TreeNode {
	for _, root := range t.RootNodes {
		if node := root.FindNode(name); node != nil {
			return node
		}
	}
	return nil
}

// FindNode 在当前节点的子树中递归查找名字为 name 的节点
func (n *TreeNode) FindNode(name string) *TreeNode {
	if n.Name == name {
		return n
	}
	for _, child := range n.Children {
		if result := child.FindNode(name); result != nil {
			return result
		}
	}
	return nil
}

// PrintTree 打印节点及其子节点
func (n *TreeNode) PrintTree(indent string, isLast bool) {
	fmt.Print(indent)
	if isLast {
		fmt.Print("└── ")
		indent += "    "
	} else {
		fmt.Print("├── ")
		indent += "│   "
	}
	fmt.Println(n.Name)

	for i := 0; i < len(n.Children); i++ {
		n.Children[i].PrintTree(indent, i == len(n.Children)-1)
	}
}

// PrintTree 打印整个树
func (t *Tree) PrintTree() {
	for i := 0; len(t.RootNodes) > 0 && i < len(t.RootNodes); i++ {
		t.RootNodes[i].PrintTree("", i == len(t.RootNodes)-1)
	}
}
