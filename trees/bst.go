package trees

import (
	"math"
)

type BSTNode struct {
	value int
	left  *BSTNode
	right *BSTNode
}

type BST struct {
	root *BSTNode
}

func (b *BST) Insert(value int) {
	newNode := &BSTNode{
		value: value,
	}

	if b.root == nil {
		b.root = newNode
	} else {
		c := b.root
		for c != nil {
			if value < c.value {
				if c.left == nil {
					c.left = newNode
					break
				}
				c = c.left
			} else {
				if c.right == nil {
					c.right = newNode
					break
				}
				c = c.right
			}
		}
	}
}

func (b *BST) Remove(value int) bool {
	if b.root == nil {
		return false
	}

	c := b.root
	var parent *BSTNode

	for c != nil {
		// traverse the tree first, assume we are not at the node to remove
		if value < c.value {
			parent = c
			c = c.left
		} else if value > c.value {
			parent = c
			c = c.right
		} else {
			// we are at the node to remove
			// case 1: no children or 1 child
			if c.left == nil {
				// replace the node with its right child (could be nil)
				if parent == nil { // removing the root
					b.root = c.right
				} else if c == parent.left {
					parent.left = c.right
				} else {
					parent.right = c.right
				}
				return true // root has been removed
			} else if c.right == nil {
				// replace the node with its left child (could be nil)
				if parent == nil {
					b.root = c.left
				} else if c == parent.left {
					parent.left = c.left
				} else {
					parent.right = c.left
				}
				return true
			} else { // case 2: 2 children
				// find the in-order successor (smallest node in right subtree)
				sp := c      // successor parent
				s := c.right // successor
				for s.left != nil {
					sp = s
					s = s.left
				}
				// copy values
				c.value = s.value

				// remove the successor
				if sp == c {
					sp.right = s.right
				} else {
					sp.left = s.right
				}
				return true
			}
		}
	}

	return false
}

func (b *BST) Get(value int) *BSTNode {
	if b.root == nil {
		return nil
	}

	c := b.root
	for c != nil {
		if value == c.value {
			return c
		} else if value < c.value {
			c = c.left
		} else {
			c = c.right
		}
	}

	return nil
}

// bfs
func (b *BST) GetMinDepth() int {
	if b.root == nil {
		return 0
	}

	queue := []*BSTNode{b.root}
	depth := 0

	for len(queue) > 0 {
		depth++
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			c := queue[0]
			queue = queue[1:]

			if c.left == nil && c.right == nil {
				return depth
			}

			if c.left != nil {
				queue = append(queue, c.left)
			}

			if c.right != nil {
				queue = append(queue, c.right)
			}
		}
	}

	return 0
}

// dfs
func (b *BST) GetMaxDepth() int {
	if b.root == nil {
		return 0
	}

	type NodeDepth struct {
		BSTNode *BSTNode
		Depth   int
	}

	stack := []NodeDepth{{b.root, 1}}
	maxDepth := 0

	for len(stack) > 0 {
		// pop the last element
		li := len(stack) - 1
		c := stack[li]
		stack = stack[:li]

		maxDepth = int(math.Max(float64(maxDepth), float64(c.Depth)))

		// push the right child first
		if c.BSTNode.right != nil {
			stack = append(stack, NodeDepth{c.BSTNode.right, c.Depth + 1})
		}

		// push the left child
		if c.BSTNode.left != nil {
			stack = append(stack, NodeDepth{c.BSTNode.left, c.Depth + 1})
		}
	}

	return maxDepth
}

// -- Helpers for Testing and Stuff --
func (b *BST) InOrderTraversal() []int {
	result := []int{}
	var traverse func(node *BSTNode)
	traverse = func(node *BSTNode) {
		if node == nil {
			return
		}
		traverse(node.left)
		result = append(result, node.value)
		traverse(node.right)
	}
	traverse(b.root)
	return result
}
