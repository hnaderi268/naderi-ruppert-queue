package main

import (
	"sync"
)

var wg sync.WaitGroup

type job func(interface{})

func DoJob(input interface{}, task job) {
	wg.Add(1)
	go func() {
		task(input)
		wg.Done()
	}()
}

func Wait() {
	wg.Wait()
}
