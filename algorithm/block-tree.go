package algorithm

import (
	"unsafe"
)

type BlockTree struct {
	root   *Node
	leaves []*Node
}

type Node struct {
	id     int64
	parent *Node
	left   *Node
	right  *Node
	blocks []unsafe.Pointer // []block
	head   int64
}

// A block implicitly shows the operations propagated to a node concurrently.
type block struct {
	value    interface{}
	sumEnq   int64
	sumDeq   int64
	super    int64
	size     int64
	endLeft  int64
	endRight int64
}

// Returns an BlockTree with tree length h and upper bound e engueues. (e, h>1)
func blockTree(h int64, e int64) BlockTree {
	if h == 1 {
		return baseCaseBlockTree(e, PID())
	}
	root := emptyNode(e, PID())
	leftSubTree := blockTree(h-1, e)
	rightSubTree := blockTree(h-1, e)
	root.left = leftSubTree.root
	root.right = rightSubTree.root
	leftSubTree.root.parent = root
	rightSubTree.root.parent = root

	return BlockTree{
		root:   root,
		leaves: append(leftSubTree.leaves, rightSubTree.leaves...),
	}
}

// Base case of block tree is a single node.
func baseCaseBlockTree(e int64, i int64) BlockTree {
	n := emptyNode(e, i)
	leaves := make([]*Node, 1)
	leaves[0] = n
	return BlockTree{root: n, leaves: leaves}
}

// Creates a node with id i, which can contain e ops.
func emptyNode(e int64, i int64) *Node {
	b := block{
		value:    nil,
		sumEnq:   0,
		sumDeq:   0,
		super:    1,
		size:     0,
		endLeft:  0,
		endRight: 0,
	}

	blocks := make([]unsafe.Pointer, e)
	blocks[0] = unsafe.Pointer(&b)
	n := &Node{
		parent: nil,
		left:   nil,
		right:  nil,
		blocks: blocks,
		head:   1,
		id:     i,
	}
	return n
}

func isEmpty(b *block) bool {
	return b == nil || *b == block{}
}

// Creates a block to be installed to n.blocks[i] from n's children's blocks.
func (n *Node) createBlock(i int64) block {
	b := block{
		endLeft:  n.left.head - 1,
		endRight: n.right.head - 1,
	}

	left := load(&n.left.blocks[b.endLeft])
	right := load(&n.right.blocks[b.endRight])

	b.sumEnq = left.sumEnq + right.sumEnq
	b.sumDeq = left.sumDeq + right.sumDeq

	prev := load(&n.blocks[i-1])
	numEnq := b.sumEnq - prev.sumEnq
	numDeq := b.sumDeq - prev.sumDeq

	if numEnq+numDeq == 0 { // Created block does not contain any operations.
		return block{}
	}

	if n.parent == nil { // Size of the queue after a block, is only computed in the root.
		b.size = max(0, prev.size+numEnq-numDeq)
	}

	return b
}
