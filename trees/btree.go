package trees

import (
	"math"
	"slices"
	"sort"
)

type Btree struct {
	root    *BtreeNode
	order   int
	minKeys int
	maxKeys int
	height  int
}

type BtreeNode struct {
	keys     []int
	values   []any
	children []*BtreeNode
	isLeaf   bool
}

func NewBtree(order int) *Btree {
	if order < 3 {
		order = 3
	}
	return &Btree{
		order:   order,
		minKeys: int(math.Ceil(float64(order)/2)) - 1,
		maxKeys: order - 1,
		height:  0,
	}
}

func (b *Btree) Insert(key int, value int) {
	if b.root == nil {
		b.root = &BtreeNode{
			keys:   []int{key},
			values: []any{value},
			isLeaf: true,
		}
		b.height++
		return
	}

	// if the root is full, we need to split it
	if len(b.root.keys) == b.maxKeys {
		newRoot := &BtreeNode{
			keys:     []int{},
			values:   []any{},
			children: []*BtreeNode{b.root},
			isLeaf:   false,
		}
		b.root = newRoot
		b.splitChild(b.root, 0)
		b.insertNonFull(b.root, key, value)
		b.height++
	} else {
		b.insertNonFull(b.root, key, value)
	}
}

func (b *Btree) insertNonFull(node *BtreeNode, key int, value int) {
	if node.isLeaf {
		// find the insertion point using binary search
		ip := sort.Search(len(node.keys), func(i int) bool {
			return node.keys[i] >= key
		})

		// insert the key and value at ip
		node.keys = append(node.keys, 0)
		node.values = append(node.values, 0)
		copy(node.keys[ip+1:], node.keys[ip:])
		copy(node.values[ip+1:], node.values[ip:])
		node.keys[ip] = key
		node.values[ip] = value
	} else {
		// find the child to insert the key and value
		ip_child_idx := sort.Search(len(node.keys), func(i int) bool {
			return node.keys[i] >= key
		})

		// if the child is full, split it before going down
		if len(node.children[ip_child_idx].keys) == b.maxKeys {
			b.splitChild(node, ip_child_idx) // node is parent, ip_child_idx is index of child in parent.children
			// after splitting, the key might go into the new right sibling
			if key > node.keys[ip_child_idx] { // Compare with the key that was just promoted to parent
				ip_child_idx++ // If key is greater, target the new right sibling
			}
		}
		b.insertNonFull(node.children[ip_child_idx], key, value) // Descend into the correct child
	}
}

func (b *Btree) splitChild(parent *BtreeNode, index int) {
	child := parent.children[index]
	newSibling := &BtreeNode{
		isLeaf: child.isLeaf,
	}

	medianKey := child.keys[b.minKeys]
	medianVal := child.values[b.minKeys]

	// move the keys and values after the median to the new sibling
	newSibling.keys = append(newSibling.keys, child.keys[b.minKeys+1:]...)
	newSibling.values = append(newSibling.values, child.values[b.minKeys+1:]...)

	// if this isnt a leaf, move the children after the median
	if !child.isLeaf {
		newSibling.children = append(newSibling.children, child.children[b.minKeys+1:]...)
		child.children = child.children[:b.minKeys+1]
	}

	// remove the original childs keys and values and children
	child.keys = child.keys[:b.minKeys]
	child.values = child.values[:b.minKeys]

	// insert the median key and value into the parent
	ip := sort.Search(len(parent.keys), func(i int) bool {
		return parent.keys[i] >= medianKey
	})

	// shift keys and values to make space
	parent.keys = append(parent.keys, 0)
	copy(parent.keys[ip+1:], parent.keys[ip:])
	parent.keys[ip] = medianKey

	parent.values = append(parent.values, 0)
	copy(parent.values[ip+1:], parent.values[ip:])
	parent.values[ip] = medianVal

	// shift the children to make space
	parent.children = append(parent.children, nil)
	copy(parent.children[ip+2:], parent.children[ip+1:])
	parent.children[ip+1] = newSibling
}

func (b *Btree) Remove(key int) bool {
	if b.root == nil || len(b.root.keys) == 0 {
		// the tree is empty or root is empty
		return false
	}

	removed := b.remove(b.root, key)

	// if the root node becomes empty after deletion
	if b.root != nil && len(b.root.keys) == 0 {
		if !b.root.isLeaf {
			// if root was an internal node and is now empty
			// its first child becomes the new root
			b.root = b.root.children[0]
			b.height--
		} else {
			b.root = nil
			b.height = 0
		}
	}

	return removed
}

func (b *Btree) remove(node *BtreeNode, key int) bool {
	// 1. find the index of the key or the child to decend into
	idx := sort.Search(len(node.keys), func(i int) bool {
		return node.keys[i] >= key
	})

	// 2. key found in current node
	if idx < len(node.keys) && node.keys[idx] == key {
		if node.isLeaf {
			// case 1: key is in a leaf node
			b.removeFromLeaf(node, idx)
			return true
		}
		// case 2: key is in an internal node
		return b.removeFromInternalNode(node, idx, key)
	}

	// 3. key not found in the current node, decend to appropriate child
	if node.isLeaf {
		return false // key not found
	}

	// the child to descend into
	childIdx := idx
	child := node.children[childIdx]

	// ensure the child has at least b.minKeys +1 keys before descending
	if len(child.keys) == b.minKeys {
		b.fillChild(node, childIdx)
		return b.remove(node, key)
	}

	return b.remove(child, key)
}

func (b *Btree) removeFromLeaf(node *BtreeNode, keyIdx int) {
	node.keys = slices.Delete(node.keys, keyIdx, keyIdx+1)
	node.values = slices.Delete(node.values, keyIdx, keyIdx+1)
}

func (b *Btree) removeFromInternalNode(node *BtreeNode, keyIdx int, key int) bool {
	lChild := node.children[keyIdx]
	rChild := node.children[keyIdx+1]

	// case 2a: left child has at least t keys (b.minKeys + 1)
	if len(lChild.keys) > b.minKeys {
		predKey, predVal := b.getPredecessor(lChild)
		node.keys[keyIdx] = predKey
		node.values[keyIdx] = predVal
		return b.remove(lChild, predKey) // remove the predecessor from the left child
	}

	// case 2b: right child has at least t keys
	if len(rChild.keys) > b.minKeys {
		succKey, succVal := b.getSuccessor(rChild)
		node.keys[keyIdx] = succKey
		node.values[keyIdx] = succVal
		return b.remove(rChild, succKey) // remove the successor from the right child
	}

	// case 2c: both left and right child have exactly t-1 keys
	// merge right and left
	mergedNode := b.mergeChildren(node, keyIdx)
	return b.remove(mergedNode, key)
}

func (b *Btree) getPredecessor(node *BtreeNode) (int, interface{}) {
	c := node
	for !c.isLeaf {
		c = c.children[len(c.children)-1]
	}
	lastIdx := len(c.keys) - 1
	return c.keys[lastIdx], c.values[lastIdx]
}

func (b *Btree) getSuccessor(node *BtreeNode) (int, interface{}) {
	c := node
	for !c.isLeaf {
		c = c.children[0]
	}
	return c.keys[0], c.values[0]
}

func (b *Btree) fillChild(parent *BtreeNode, childIdx int) {
	// try borrowing from the left
	if childIdx > 0 && len(parent.children[childIdx-1].keys) > b.minKeys {
		b.borrowFromLeft(parent, childIdx)
	} else if childIdx < len(parent.keys) && len(parent.children[childIdx+1].keys) > b.minKeys {
		// try borrowing from the right
		b.borrowFromRight(parent, childIdx)
	} else {
		// both siblings have b.minKeys, merge
		if childIdx < len(parent.keys) {
			b.mergeChildren(parent, childIdx)
		} else {
			// deficient child is the rightmost, merge with its left
			b.mergeChildren(parent, childIdx-1)
		}
	}
}

func (b *Btree) borrowFromLeft(parent *BtreeNode, childIdx int) {
	child := parent.children[childIdx]
	lSibling := parent.children[childIdx-1]

	// key from parent moves down to child
	keyFromParent := parent.keys[childIdx-1]
	valFromParent := parent.values[childIdx-1]

	// prepend key and value
	child.keys = append([]int{keyFromParent}, child.keys...)
	child.values = append([]any{valFromParent}, child.values...)

	// last key from left sibling moves up to the parent
	parent.keys[childIdx-1] = lSibling.keys[len(lSibling.keys)-1]
	parent.values[childIdx-1] = lSibling.values[len(lSibling.values)-1]

	// if not a leaf, move the child pointer from left sibling to child
	if !lSibling.isLeaf {
		child.children = append([]*BtreeNode{lSibling.children[len(lSibling.children)-1]}, child.children...)
		lSibling.children = lSibling.children[:len(lSibling.children)-1]
	}
}

func (b *Btree) borrowFromRight(parent *BtreeNode, childIdx int) {
	child := parent.children[childIdx]
	rSibling := parent.children[childIdx+1]

	// key from parent moves down to child
	keyFromParent := parent.keys[childIdx]
	valFromParent := parent.values[childIdx]

	// append key and value to child
	child.keys = append(child.keys, keyFromParent)
	child.values = append(child.values, valFromParent)

	// first key from right sibling moves up to the parent
	parent.keys[childIdx] = rSibling.keys[0]
	parent.values[childIdx] = rSibling.values[0]

	// remove the first key from the right sibling
	rSibling.keys = rSibling.keys[1:]
	rSibling.values = rSibling.values[1:]

	// if not a leaf, move child from right sibling to child
	if !rSibling.isLeaf {
		child.children = append(child.children, rSibling.children[0])
		rSibling.children = rSibling.children[1:]
	}
}

func (b *Btree) mergeChildren(parent *BtreeNode, keyIdx int) *BtreeNode {
	lChild := parent.children[keyIdx]
	rChild := parent.children[keyIdx+1]

	// key from parent to be moved down
	keyFromParent := parent.keys[keyIdx]
	valFromParent := parent.values[keyIdx]

	// append keyFromparent and all k/v from rChild to lChild
	lChild.keys = append(lChild.keys, keyFromParent)
	lChild.values = append(lChild.values, valFromParent)
	lChild.keys = append(lChild.keys, rChild.keys...)
	lChild.values = append(lChild.values, rChild.values...)

	// if not leaves, append children from right to left
	if !lChild.isLeaf {
		lChild.children = append(lChild.children, rChild.children...)
	}

	// remove key and right child pointer from parent
	parent.keys = slices.Delete(parent.keys, keyIdx, keyIdx+1)
	parent.values = slices.Delete(parent.values, keyIdx, keyIdx+1)
	parent.children = slices.Delete(parent.children, keyIdx+1, keyIdx+2)

	return lChild // the new merged node
}

func (b *Btree) Get(key int) (any, bool) {
	if b.root == nil {
		return nil, false
	}

	return b.search(b.root, key)
}

func (b *Btree) search(node *BtreeNode, key int) (any, bool) {
	idx := sort.Search(len(node.keys), func(i int) bool {
		return node.keys[i] >= key
	})

	if idx < len(node.keys) && node.keys[idx] == key {
		return node.values[idx], true
	} else if node.isLeaf {
		return nil, false
	} else {
		return b.search(node.children[idx], key)
	}
}

func (b *Btree) FindMaxDepth() int {
	if b.root == nil {
		return 0
	}
	return b.height
}

func (b *Btree) FindMinDepth() int {
	if b.root == nil {
		return 0
	}
	return b.height
}

// -- Helpers for Testing and Stuff --
func (b *Btree) GetKeysInOrder() []int {
	var result []int
	var traverse func(node *BtreeNode)
	traverse = func(node *BtreeNode) {
		if node == nil {
			return
		}
		if node.isLeaf {
			result = append(result, node.keys...)
			return
		}
		// For internal nodes
		for i := 0; i < len(node.keys); i++ {
			traverse(node.children[i])
			result = append(result, node.keys[i])
		}
		traverse(node.children[len(node.keys)]) // Traverse the rightmost child
	}
	traverse(b.root)
	return result
}
