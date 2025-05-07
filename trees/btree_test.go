package trees

import (
	"reflect"
	"testing"
)

func TestBtree_InsertAndGet_Simple(t *testing.T) {
	b := NewBtree(3)

	testCases := []struct {
		key   int
		value int
	}{
		{10, 100},
		{20, 200},
		{5, 50},
		{15, 150},
		{25, 250},
	}

	for _, tc := range testCases {
		b.Insert(tc.key, tc.value)
	}

	for _, tc := range testCases {
		val, found := b.Get(tc.key)
		if !found {
			t.Errorf("Get(%d): expected found=true, got false", tc.key)
		}
		if val != tc.value {
			t.Errorf("Get(%d): expected value=%v, got %v", tc.key, tc.value, val)
		}
	}

	_, found := b.Get(1000)
	if found {
		t.Errorf("Get(1000): expected found=false, got true for non-existent key")
	}

	emptyTree := NewBtree(3)
	_, foundEmpty := emptyTree.Get(1)
	if foundEmpty {
		t.Errorf("Get(1) on empty tree: expected found=false, got true")
	}
}

func TestBtree_GetKeysInOrder_Simple(t *testing.T) {
	b := NewBtree(3)
	keysToInsert := []int{10, 20, 5, 15, 25, 3, 30}
	expectedOrder := []int{3, 5, 10, 15, 20, 25, 30}

	for _, key := range keysToInsert {
		b.Insert(key, key*10)
	}

	actualOrder := b.GetKeysInOrder()
	if !reflect.DeepEqual(actualOrder, expectedOrder) {
		t.Errorf("GetKeysInOrder(): expected %v, got %v", expectedOrder, actualOrder)
	}

	emptyTree := NewBtree(3)
	if len(emptyTree.GetKeysInOrder()) != 0 {
		t.Errorf("GetKeysInOrder() on empty tree: expected empty slice, got %v", emptyTree.GetKeysInOrder())
	}
}

func TestBtree_Insert_Splits(t *testing.T) {
	tests := []struct {
		name           string
		order          int
		inserts        []int
		expectedKeys   []int
		expectedHeight int
	}{
		{
			name:           "Order 3, no split initially",
			order:          3,
			inserts:        []int{10, 20},
			expectedKeys:   []int{10, 20},
			expectedHeight: 1,
		},
		{
			name:           "Order 3, leaf split, no root split",
			order:          3,
			inserts:        []int{10, 20, 5},
			expectedKeys:   []int{5, 10, 20},
			expectedHeight: 2,
		},
		{
			name:           "Order 3, root split",
			order:          3,
			inserts:        []int{10, 20, 5, 15},
			expectedKeys:   []int{5, 10, 15, 20},
			expectedHeight: 2,
		},
		{
			name:           "Order 3, further splits",
			order:          3,
			inserts:        []int{10, 20, 5, 15, 25, 3, 7},
			expectedKeys:   []int{3, 5, 7, 10, 15, 20, 25},
			expectedHeight: 3,
		},
		{
			name:           "Order 5, root split",
			order:          5,
			inserts:        []int{10, 20, 30, 40, 50},
			expectedKeys:   []int{10, 20, 30, 40, 50},
			expectedHeight: 2,
		},
		{
			name:           "Order 4, several splits leading to height 3",
			order:          4,
			inserts:        []int{1, 2, 3, 4, 5, 6, 0, -1, -2, 7, 8},
			expectedKeys:   []int{-2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8},
			expectedHeight: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBtree(tc.order)
			for i, key := range tc.inserts {
				b.Insert(key, key*10)
				t.Logf(
					"After inserting %d (item %d/%d), tree height: %d, keys: %v",
					key,
					i+1,
					len(tc.inserts),
					b.height,
					b.GetKeysInOrder(),
				)
			}

			actualKeys := b.GetKeysInOrder()
			if !reflect.DeepEqual(actualKeys, tc.expectedKeys) {
				t.Errorf("GetKeysInOrder() got %v, want %v", actualKeys, tc.expectedKeys)
			}

			if b.height != tc.expectedHeight {
				t.Errorf("height got %d, want %d", b.height, tc.expectedHeight)
			}

			if b.FindMinDepth() != tc.expectedHeight {
				t.Errorf("FindMinDepth() got %d, want %d", b.FindMinDepth(), tc.expectedHeight)
			}
			if b.FindMaxDepth() != tc.expectedHeight {
				t.Errorf("FindMaxDepth() got %d, want %d", b.FindMaxDepth(), tc.expectedHeight)
			}
		})
	}
}

func TestBtree_DepthMethods(t *testing.T) {
	t.Run("Empty tree", func(t *testing.T) {
		b := NewBtree(3)
		if b.FindMinDepth() != 0 {
			t.Errorf("FindMinDepth() on empty tree: expected 0, got %d", b.FindMinDepth())
		}
		if b.FindMaxDepth() != 0 {
			t.Errorf("FindMaxDepth() on empty tree: expected 0, got %d", b.FindMaxDepth())
		}
		if b.height != 0 {
			t.Errorf("height on empty tree: expected 0, got %d", b.height)
		}
	})

	t.Run("Tree with one node", func(t *testing.T) {
		b := NewBtree(3)
		b.Insert(10, 100)
		if b.FindMinDepth() != 1 {
			t.Errorf("FindMinDepth() on single node tree: expected 1, got %d", b.FindMinDepth())
		}
		if b.FindMaxDepth() != 1 {
			t.Errorf("FindMaxDepth() on single node tree: expected 1, got %d", b.FindMaxDepth())
		}
		if b.height != 1 {
			t.Errorf("height on single node tree: expected 1, got %d", b.height)
		}
	})
}
