package algorithm

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// There exists a wait-free counter with logarithmic step complexity.
var counter int64

// TODO: move to UTIl and remove from block-tree.
// Returns an integer greater than all values previously returned by PID.
func PID() int64 {
	return atomic.AddInt64(&counter, 1)
}

func isPowerOfTwo(n int64) bool {
	return 0 == (n & (n - 1))
}

//func max(a, b int64) int64 {
//	return int64(math.Max(float64(a), float64(b)))
//}

func load(p *unsafe.Pointer) *block {
	return (*block)(atomic.LoadPointer(p))
}

func CASBlock(addr *unsafe.Pointer, old, new *block) bool {
	unsafeAddr := (*unsafe.Pointer)(unsafe.Pointer(addr))
	unsafeNew := unsafe.Pointer(new)
	return atomic.CompareAndSwapPointer(unsafeAddr, nil, unsafeNew)
}

func CASInt(oldAddr *int64, old int64, new int64) bool {
	return atomic.CompareAndSwapInt64(oldAddr, int64(old), int64(new))
}

// For left<=i<=right binarySearch returns min({i: blocks[i].sumEnq>=key})
func binarySearch(blocks []unsafe.Pointer, left int64, right int64, key int64) int64 {
	if right <= left {
		return left
	}
	mid := (left + right) / 2
	midBlock := load(&blocks[mid])
	if midBlock.sumEnq < key {
		return binarySearch(blocks, mid+1, right, key)
	}
	return binarySearch(blocks, left, mid, key)
}

func (q *BlockTree) Print() {
	queue := []*Node{}
	queue = append(queue, q.root)
	list := bfsUtil(queue, 0)
	fmt.Println("|sumEng sumDeq endLeft endRight (size, super, value] |")
	printUtil(list)
}

func bfsUtil(nodes []*Node, i int64) []*Node {
	if int64(len(nodes)) == i {
		return nodes
	}
	if nodes[i].left != nil {
		nodes = append(nodes, nodes[i].left)
	}
	if nodes[i].right != nil {
		nodes = append(nodes, nodes[i].right)
	}
	return bfsUtil(nodes, 1+1)
}

func printUtil(list []*Node) {
	for i, n := range list {
		fmt.Println(n, " ")
		if isPowerOfTwo(int64(i + 2)) { // Break line as height of n increases.
			fmt.Println()
		}
	}
}

// Returns the blocks stored in a node.
func (n Node) String() string {
	res := "["
	for _, block := range n.blocks {
		if !isEmpty(load(&block)) {
			isRoot := n.parent == nil
			isLeaf := n.left == nil
			res = res + load(&block).String(isRoot, isLeaf)
		}
	}
	res = res + "]"
	return fmt.Sprint(n.id, res)
}

// Return a block's information with regard to the position of the node containing it in the tree.
func (b *block) String(isRoot bool, isLeaf bool) string {
	res := fmt.Sprint(b.sumEnq, b.sumDeq)
	if isLeaf {
		res = fmt.Sprint(res, " ", b.super, " ", b.value)
	} else if isRoot {
		res = fmt.Sprint(res, " ", b.endLeft, b.endRight, b.size)
	} else {
		res = fmt.Sprint(res, " ", b.endLeft, b.endRight, b.super)
	}
	return fmt.Sprint(res, "|")
}
