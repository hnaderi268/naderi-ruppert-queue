package algorithm

import "unsafe"

// Assumption: There will be no call of Enqueue(*, i), until Enqueue(v, i) terminates. (0<=i<#processes)
// Enqueue puts the given value v at the tail of the queue.
func (q *BlockTree) Enqueue(v interface{}, pid int64) {
	leaf := q.leaves[pid]
	leafBlock := load(&leaf.blocks[leaf.head-1])
	b := block{
		value:  v,
		sumEnq: leafBlock.sumEnq + 1,
		sumDeq: leafBlock.sumDeq,
		super:  leaf.parent.head,
	}
	q.append(b, pid)
}

// Assumption: There will be no call of Dequeue(i), until Dequeue(i) terminates. (O<=i<#processes)
// Dequeue removes the value v at the head of the algorithm and returns v.
func (q *BlockTree) Dequeue(pid int64) interface{} {
	leaf := q.leaves[pid]
	leafBlock := load(&leaf.blocks[leaf.head-1])
	b := block{
		sumEnq: leafBlock.sumEnq,
		sumDeq: leafBlock.sumDeq + 1,
		super:  leaf.parent.head,
	}
	q.append(b, pid)
	c, d := leaf.indexDequeue(pid, leaf.head-1, 1)
	v := q.findResponse(pid, c, d)
	return v
}

func (q *BlockTree) append(b block, i int64) {
	leaf := q.leaves[i]
	leaf.blocks[leaf.head] = unsafe.Pointer(&b)
	leaf.head = leaf.head + 1
	q.propagate(leaf.parent)
}

func (q *BlockTree) propagate(n *Node) {
	RandomSleep(2)
	if !n.refresh() {
		n.refresh()
	}
	if n != q.root {
		q.propagate(n.parent)
	}
}

func (n *Node) refresh() bool {
	head := n.head
	n.left.tryAdvance()
	n.right.tryAdvance()
	newBlock := n.createBlock(head)
	if isEmpty(&newBlock) {
		return true
	}
	b := load(&n.blocks[head])
	if !isEmpty(b) {
		return false
	}
	result := CASBlock(&n.blocks[head], nil, &newBlock)
	n.advance(head)
	return result
}

// Precondition: node.blocks [b] is not null, was propagated to root, and contains at least i Dequeue operations.
// Returns b, 1 such that ith dequeue in node.blocks[b] is the ith dequeue of root.blocks [b]
func (node *Node) indexDequeue(id int64, b int64, i int64) (int64, int64) {
	if node.parent == nil { // node is the root
		return b, i
	}
	curr := load(&node.blocks[b])
	prev := load(&node.blocks[b-1])
	sup := curr.super

	supBlock := load(&node.parent.blocks[sup])
	// if node is a left child
	if node.parent.left == node {
		if supBlock.endLeft < b {
			sup++
		}
		supBlock = load(&node.parent.blocks[sup])
		prevSupBlock := load(&node.parent.blocks[sup-1])
		// compute index of dequeue in superblock
		prevSupEndLeftBlock := load(&node.blocks[prevSupBlock.endLeft])
		i += prev.sumDeq - prevSupEndLeftBlock.sumDeq
		return node.parent.indexDequeue(id, sup, i)
	} else {
		// if node is a right child
		if supBlock.endRight < b {
			sup++
		}
		supBlock = load(&node.parent.blocks[sup])
		prevSupBlock := load(&node.parent.blocks[sup-1])
		// compute index of dequeue in superblock
		supEndLeftBlock := load(&node.parent.left.blocks[supBlock.endLeft])
		prevSupEndLeftBlock := load(&node.parent.left.blocks[prevSupBlock.endLeft])
		prevSupEndRightBlock := load(&node.blocks[prevSupBlock.endRight])
		i += prev.sumDeq - prevSupEndRightBlock.sumDeq + supEndLeftBlock.sumDeq - prevSupEndLeftBlock.sumDeq
	}
	return node.parent.indexDequeue(id, sup, 1)
}

// findResponse(id,b,i) returns the response to the ith dequeue in bth block in q.root. (1<=i<=q.root
func (q *BlockTree) findResponse(id int64, b int64, i int64) interface{} {
	this := load(&q.root.blocks[b])
	prev := load(&q.root.blocks[b-1])
	numEnq := this.sumEnq - prev.sumEnq

	if prev.size+numEnq < i { // queue is empty when dequeue occurs
		return nil
	}
	// response is the eth enqueue in the root
	e := i + prev.sumEnq - prev.size

	// find min be <= b with root.blocks[be] sumEng >= e
	be := binarySearch(q.root.blocks, 1, b, e)
	containingBlock := load(&q.root.blocks[be-1])

	// rank of enqueue within its block
	ie := e - containingBlock.sumEnq
	return q.root.getEnqueue(id, be, ie)
}

// getEnqueue() returns argument of the ith enqueue in node.blocks[b]
// 1<=i, node.blocks[b] is non-null and contains at least i enqueues
func (node *Node) getEnqueue(id int64, b int64, i int64) interface{} {
	curr := load(&node.blocks[b])
	prev := load(&node.blocks[b-1])

	// Base case when node is a leaf.
	if node.left == nil {
		return curr.value
	}

	// sumLeft is the number of enqueues in node.blocks[1..b] from node's left child
	endLeftBlock := load(&node.left.blocks[curr.endLeft])
	sumLeft := endLeftBlock.sumEnq

	// sumPrevleft is the number of enqueues in node.blocks[1..b-1] from node's left child
	var sumPrevLeft int64 = 0
	if curr.endLeft > 0 {
		prevEndLeftBlock := load(&node.left.blocks[prev.endLeft])
		sumPrevLeft = prevEndLeftBlock.sumEnq
	}
	if i <= sumLeft-sumPrevLeft { // required enqueue is in node.left
		// find min=minimum index in [node.blocks[b-1].endLeft+1..node.blocks[b].endLeft] such that \\TODO
		// node.left.blocks[min].sumEnq >= 1 + sumPrevleft
		min := binarySearch(node.left.blocks, prev.endLeft+1, curr.endLeft, i+sumPrevLeft)
		minb := load(&node.left.blocks[min-1])
		mini := i - (minb.sumEnq - sumPrevLeft)
		return node.left.getEnqueue(id, min, mini)
	}
	// required enqueue is in node. right
	i = i - (sumLeft - sumPrevLeft)

	//sumPrevRight is thenumber of enqueues in node.blocks[1..b-1] from node's right child
	prevEndRightBlock := load(&node.right.blocks[prev.endRight])
	sumPrevRight := prevEndRightBlock.sumEnq

	min := binarySearch(node.right.blocks, prev.endRight+1, curr.endRight, i+sumPrevRight)
	minb := load(&node.right.blocks[min-1])
	mini := i - (minb.sumEnq - sumPrevRight)
	return node.right.getEnqueue(id, min, mini)
}

// tryAdvance advances n if needed.
func (n *Node) tryAdvance() {
	childHead := n.head
	childBlock := load(&n.blocks[childHead])
	if !isEmpty(childBlock) { // check 1f block is not empty
		n.advance(childHead)
	}
}

// advance sets n.blocks[h].super and increments v.head from h to h+1
func (n *Node) advance(h int64) {
	if n.parent != nil {
		h_p := n.parent.head
		b_h := load(&n.blocks[h])
		CASInt(&b_h.super, 0, h_p)
	}
	CASInt(&n.head, h, h+1)
}
