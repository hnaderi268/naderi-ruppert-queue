package main

import (
	"naderi-ruppert-queue/algorithm"
	"naderi-ruppert-queue/io"
)

// Assumption: There will be no call of Enqueue(*, i), until Enqueue(v, i) terminates. (0<=i<#processes)
// Enqueue puts the given value v at the tail of the queue.
func Enqueue(q algorithm.BlockTree, v interface{}, pid int64) {
	io.Log(pid, io.START_ENQ, v)
	algorithm.RandomSleep(3)
	q.Enqueue(v, pid)
	io.Log(pid, io.END_ENQ)
}

// Assumption: There will be no call of Dequeue(i), until Dequeue(i) terminates. (O<=i<#processes)
// Dequeue removes the value v at the head of the algorithm and returns v.
func Dequeue(q algorithm.BlockTree, pid int64) interface{} {
	io.Log(pid, io.START_DEQ)
	algorithm.RandomSleep(3)
	v := q.Dequeue(pid)
	io.Log(pid, io.END_DEQ, v)
	return v
}
