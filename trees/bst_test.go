package trees

import (
	"reflect"
	"testing"
)

func newBSTWithValues(values ...int) *BST {
	bst := &BST{}
	for _, v := range values {
		bst.Insert(v)
	}
	return bst
}

func TestBST_Insert(t *testing.T) {
	tests := []struct {
		name          string
		initialValues []int
		insertValue   int
		expectedOrder []int
		expectedRoot  int
	}{
		{
			name:          "insert into empty tree",
			initialValues: []int{},
			insertValue:   10,
			expectedOrder: []int{10},
			expectedRoot:  10,
		},
		{
			name:          "insert smaller value",
			initialValues: []int{10},
			insertValue:   5,
			expectedOrder: []int{5, 10},
			expectedRoot:  10,
		},
		{
			name:          "insert larger value",
			initialValues: []int{10},
			insertValue:   15,
			expectedOrder: []int{10, 15},
			expectedRoot:  10,
		},
		{
			name:          "insert multiple values",
			initialValues: []int{10, 5, 15},
			insertValue:   3,
			expectedOrder: []int{3, 5, 10, 15},
			expectedRoot:  10,
		},
		{
			name:          "insert another larger value",
			initialValues: []int{10, 5, 15},
			insertValue:   20,
			expectedOrder: []int{5, 10, 15, 20},
			expectedRoot:  10,
		},
		{
			name:          "insert value causing deeper left branch",
			initialValues: []int{10, 5, 15, 3},
			insertValue:   1,
			expectedOrder: []int{1, 3, 5, 10, 15},
			expectedRoot:  10,
		},
		{
			name:          "insert value causing deeper right branch",
			initialValues: []int{10, 5, 15, 20},
			insertValue:   25,
			expectedOrder: []int{5, 10, 15, 20, 25},
			expectedRoot:  10,
		},
		{
			name:          "insert duplicate value (should go to right)",
			initialValues: []int{10, 5},
			insertValue:   5,
			expectedOrder: []int{5, 5, 10},
			expectedRoot:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := newBSTWithValues(tt.initialValues...)
			bst.Insert(tt.insertValue)

			if len(tt.initialValues) == 0 {
				if bst.root == nil && len(tt.expectedOrder) > 0 {
					t.Fatalf("root is nil, expected root value %d", tt.expectedRoot)
				}
				if bst.root != nil && bst.root.value != tt.expectedRoot {
					t.Errorf("root.value = %d; want %d", bst.root.value, tt.expectedRoot)
				}
			}

			actualOrder := bst.InOrderTraversal()
			if !reflect.DeepEqual(actualOrder, tt.expectedOrder) {
				t.Errorf("InOrderTraversal() = %v; want %v", actualOrder, tt.expectedOrder)
			}
		})
	}
}

func TestBST_Get(t *testing.T) {
	bst := newBSTWithValues(10, 5, 15, 3, 7, 12, 17)

	tests := []struct {
		name        string
		valueToGet  int
		expectFound bool
		expectedVal int
	}{
		{"get from empty tree", 10, false, 0},
		{"get existing root", 10, true, 10},
		{"get existing left child", 5, true, 5},
		{"get existing right child", 15, true, 15},
		{"get existing leaf (left)", 3, true, 3},
		{"get existing leaf (right)", 17, true, 17},
		{"get existing internal node (left)", 7, true, 7},
		{"get existing internal node (right)", 12, true, 12},
		{"get non-existing (smaller than min)", 1, false, 0},
		{"get non-existing (larger than max)", 20, false, 0},
		{"get non-existing (in between)", 11, false, 0},
		{"get non-existing (in between 2)", 6, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentBST := bst
			if tt.name == "get from empty tree" {
				currentBST = &BST{}
			}

			node := currentBST.Get(tt.valueToGet)
			found := node != nil

			if found != tt.expectFound {
				t.Errorf("Get(%d) found = %v; want %v", tt.valueToGet, found, tt.expectFound)
			}

			if tt.expectFound && found {
				if node.value != tt.expectedVal {
					t.Errorf("Get(%d) node.value = %d; want %d", tt.valueToGet, node.value, tt.expectedVal)
				}
			}
		})
	}
}

func TestBST_Remove(t *testing.T) {
	tests := []struct {
		name            string
		initialValues   []int
		valueToRemove   int
		expectRemoved   bool
		expectedOrder   []int
		expectedRootVal int
	}{
		{"remove from empty tree", []int{}, 10, false, []int{}, 0},
		{"remove non-existent from single node tree", []int{10}, 5, false, []int{10}, 10},
		{"remove non-existent from multi-node tree", []int{10, 5, 15}, 7, false, []int{5, 10, 15}, 10},

		{"remove leaf (root itself)", []int{10}, 10, true, []int{}, 0},
		{"remove leaf (left child)", []int{10, 5}, 5, true, []int{10}, 10},
		{"remove leaf (right child)", []int{10, 15}, 15, true, []int{10}, 10},
		{"remove leaf (deeper left)", []int{10, 5, 15, 3}, 3, true, []int{5, 10, 15}, 10},
		{"remove leaf (deeper right)", []int{10, 5, 15, 17}, 17, true, []int{5, 10, 15}, 10},

		{"remove root with only left child", []int{10, 5}, 10, true, []int{5}, 5},
		{"remove root with only right child", []int{10, 15}, 10, true, []int{15}, 15},
		{"remove internal with only left (node is left child)", []int{20, 10, 5}, 10, true, []int{5, 20}, 20},
		{"remove internal with only left (node is right child)", []int{10, 20, 15}, 20, true, []int{10, 15}, 10},
		{"remove internal with only right (node is left child)", []int{20, 10, 15}, 10, true, []int{15, 20}, 20},
		{"remove internal with only right (node is right child)", []int{10, 20, 25}, 20, true, []int{10, 25}, 10},

		{
			name:            "remove root with two children (successor is right child)",
			initialValues:   []int{10, 5, 15},
			valueToRemove:   10,
			expectRemoved:   true,
			expectedOrder:   []int{5, 15},
			expectedRootVal: 15,
		},
		{
			name:            "remove root with two children (successor deeper)",
			initialValues:   []int{10, 5, 20, 15, 25, 12},
			valueToRemove:   10,
			expectRemoved:   true,
			expectedOrder:   []int{5, 12, 15, 20, 25},
			expectedRootVal: 12,
		},
		{
			name:            "remove internal (left child) with two children (successor is its right child)",
			initialValues:   []int{50, 20, 70, 10, 30, 25, 35},
			valueToRemove:   20,
			expectRemoved:   true,
			expectedOrder:   []int{10, 25, 30, 35, 50, 70},
			expectedRootVal: 50,
		},
		{
			name:            "remove internal (right child) with two children (successor deeper)",
			initialValues:   []int{10, 5, 50, 30, 70, 20, 40, 60, 80, 35},
			valueToRemove:   50,
			expectRemoved:   true,
			expectedOrder:   []int{5, 10, 20, 30, 35, 40, 60, 70, 80},
			expectedRootVal: 10,
		},
		{
			name:            "remove node where successor's parent is the node being removed",
			initialValues:   []int{10, 5, 15, 12, 17},
			valueToRemove:   10,
			expectRemoved:   true,
			expectedOrder:   []int{5, 12, 15, 17},
			expectedRootVal: 12,
		},
		{
			name:            "remove node where successor has a right child",
			initialValues:   []int{10, 5, 20, 15, 25, 12, 17},
			valueToRemove:   10,
			expectRemoved:   true,
			expectedOrder:   []int{5, 12, 15, 17, 20, 25},
			expectedRootVal: 12,
		},
		{
			name:            "remove all nodes one by one (complex case)",
			initialValues:   []int{10, 5, 15, 3, 7, 12, 17},
			valueToRemove:   10,
			expectRemoved:   true,
			expectedOrder:   []int{3, 5, 7, 12, 15, 17},
			expectedRootVal: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := newBSTWithValues(tt.initialValues...)
			removed := bst.Remove(tt.valueToRemove)

			if removed != tt.expectRemoved {
				t.Errorf("Remove(%d) returned %v; want %v", tt.valueToRemove, removed, tt.expectRemoved)
			}

			actualOrder := bst.InOrderTraversal()

			if len(tt.expectedOrder) == 0 {
				if bst.root != nil {
					t.Errorf(
						"Root is not nil after Remove(%d) made tree empty. Root value: %d. InOrder: %v",
						tt.valueToRemove,
						bst.root.value,
						actualOrder,
					)
				}
			} else {
				if bst.root == nil {
					t.Fatalf(
						"Root is nil after Remove(%d), but expected order %v is not empty.",
						tt.valueToRemove,
						tt.expectedOrder,
					)
				}

				if bst.root.value != tt.expectedRootVal {
					t.Errorf(
						"Root value after Remove(%d) is %d; want %d. InOrder: %v",
						tt.valueToRemove,
						bst.root.value,
						tt.expectedRootVal,
						actualOrder,
					)
				}
			}

			if !reflect.DeepEqual(actualOrder, tt.expectedOrder) {
				t.Errorf(
					"InOrderTraversal() after Remove(%d) = %v; want %v",
					tt.valueToRemove,
					actualOrder,
					tt.expectedOrder,
				)
			}
		})
	}

	t.Run("sequential removals to empty", func(t *testing.T) {
		bst := newBSTWithValues(10, 5, 15, 3, 7, 12, 17)
		valuesToRemove := []int{3, 7, 5, 12, 17, 15, 10}
		expectedOrdersAfterRemove := [][]int{
			{5, 7, 10, 12, 15, 17},
			{5, 10, 12, 15, 17},
			{10, 12, 15, 17},
			{10, 15, 17},
			{10, 15},
			{10},
			{},
		}

		for i, val := range valuesToRemove {
			removed := bst.Remove(val)
			if !removed {
				t.Errorf("Sequential Remove(%d) failed, expected true", val)
			}
			actualOrder := bst.InOrderTraversal()
			if !reflect.DeepEqual(actualOrder, expectedOrdersAfterRemove[i]) {
				t.Errorf(
					"Sequential Remove(%d): InOrder = %v; want %v",
					val,
					actualOrder,
					expectedOrdersAfterRemove[i],
				)
			}
		}
		if bst.root != nil {
			t.Error("Tree should be empty after all sequential removals, but root is not nil")
		}
	})
}

func TestBST_GetMinDepth(t *testing.T) {
	tests := []struct {
		name          string
		initialValues []int
		expectedDepth int
	}{
		{"empty tree", []int{}, 0},
		{"single node", []int{10}, 1},
		{"two nodes (left)", []int{10, 5}, 2},
		{"two nodes (right)", []int{10, 15}, 2},
		{"three nodes balanced", []int{10, 5, 15}, 2},
		{"three nodes skewed left", []int{10, 5, 3}, 3},
		{"skewed left (3 nodes)", []int{30, 20, 10}, 3},
		{"skewed right (3 nodes)", []int{10, 20, 30}, 3},
		{"complex tree 1", []int{10, 5, 15, 3, 7, 12}, 3},
		{"complex tree 2 (min depth 2)", []int{10, 5, 15, 3}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := newBSTWithValues(tt.initialValues...)
			depth := bst.GetMinDepth()
			if depth != tt.expectedDepth {
				t.Errorf("GetMinDepth() for %v = %d; want %d", tt.initialValues, depth, tt.expectedDepth)
			}
		})
	}
}

func TestBST_GetMaxDepth(t *testing.T) {
	tests := []struct {
		name          string
		initialValues []int
		expectedDepth int
	}{
		{"empty tree", []int{}, 0},
		{"single node", []int{10}, 1},
		{"two nodes (left)", []int{10, 5}, 2},
		{"two nodes (right)", []int{10, 15}, 2},
		{"three nodes balanced", []int{10, 5, 15}, 2},
		{"skewed left (3 nodes)", []int{30, 20, 10}, 3},
		{"skewed right (3 nodes)", []int{10, 20, 30}, 3},
		{"complex tree 1", []int{10, 5, 15, 3, 7, 12, 17, 1, 4, 6, 8, 11, 13, 16, 18}, 4},
		{"complex tree 2", []int{10, 5, 15, 3, 20, 2, 25, 1}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := newBSTWithValues(tt.initialValues...)
			depth := bst.GetMaxDepth()
			if depth != tt.expectedDepth {
				t.Errorf("GetMaxDepth() for %v = %d; want %d", tt.initialValues, depth, tt.expectedDepth)
			}
		})
	}
}
