package main

import (
	"naderi-ruppert-queue/algorithm"
	"naderi-ruppert-queue/io"
)

func main() {

	io.Init(8)

	q := algorithm.NewBlockTree(4, 100)

	DoJob(q, job0)
	DoJob(q, job1)
	DoJob(q, job2)
	DoJob(q, job3)
	DoJob(q, job4)
	DoJob(q, job5)
	DoJob(q, job6)
	DoJob(q, job7)
	Wait()

	q.Print()
	io.DrawExecutionTimeline()
}

func job0(interfaceQ interface{}) {
	q := interfaceQ.(algorithm.BlockTree)
	Enqueue(q, "01", 0)
	Enqueue(q, "02", 0)
	Dequeue(q, 0)
	Enqueue(q, "03", 0)
	Enqueue(q, "04", 0)
	Enqueue(q, "05", 0)
}

func job1(interfaceQ interface{}) {
	q := interfaceQ.(algorithm.BlockTree)
	Enqueue(q, "11", 1)
	Enqueue(q, "12", 1)
	Dequeue(q, 1)
	Enqueue(q, "13", 1)
	Enqueue(q, "14", 1)
}

func job2(interfaceQ interface{}) {
	q := interfaceQ.(algorithm.BlockTree)
	Dequeue(q, 2)
	Dequeue(q, 2)
	Dequeue(q, 2)
	Enqueue(q, "21", 2)
}

func job3(interfaceQ interface{}) {
	q := interfaceQ.(algorithm.BlockTree)
	Dequeue(q, 3)
	Enqueue(q, "31", 3)
	Dequeue(q, 3)
}

func job4(interfaceQ interface{}) {
	q := interfaceQ.(algorithm.BlockTree)
	Dequeue(q, 4)
	Dequeue(q, 4)
	Dequeue(q, 4)
	Enqueue(q, "41", 4)
	Enqueue(q, "42", 4)
	Enqueue(q, "43", 4)
	Enqueue(q, "44", 4)
	Enqueue(q, "45", 4)
}

func job5(interfaceQ interface{}) {
	q := interfaceQ.(algorithm.BlockTree)
	Dequeue(q, 5)
	Enqueue(q, "51", 5)
	Dequeue(q, 5)
	Enqueue(q, "52", 5)
	Dequeue(q, 5)
	Enqueue(q, "53", 5)
	Dequeue(q, 5)
	Enqueue(q, "54", 5)
}

func job6(interfaceQ interface{}) {
	q := interfaceQ.(algorithm.BlockTree)
	Dequeue(q, 6)
	Enqueue(q, "61", 6)
}

func job7(interfaceQ interface{}) {
	q := interfaceQ.(algorithm.BlockTree)
	Dequeue(q, 7)
	Enqueue(q, "71", 7)
}
