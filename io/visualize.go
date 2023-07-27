package io

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type processState int64

const (
	IDLE processState = 0 + iota
	STARTED
	WORKING
	FINISHED
)

func DrawExecutionTimeline() {
	printHeader()
	CloseWriters()
	mergedHistory := ReadAll()
	sort.Strings(mergedHistory)

	printlinearization(mergedHistory)
}
func printlinearization(mergedOutputs []string) {
	isworking := make([]processState, processesCount)
	prevProcess := -1

	for _, text := range mergedOutputs {
		var process = -1
		if strings.Contains(text, "[") {
			process, _ = strconv.Atoi(string(text[strings.Index(text, "[")+1]))
			isworking[process] = STARTED
		} else {
			process, _ = strconv.Atoi(string(text[strings.Index(text, "]")-1]))
			isworking[process] = FINISHED
		}
		if prevProcess != -1 && prevProcess != process {
			isworking[prevProcess] = (isworking[prevProcess] + 1) % 4 // IDLE -> STARTED -> WORKING ->FINISHED
		}
		prevProcess = process
		fmt.Println(drawProcessesStates(isworking), "|", text[16:])
	}
}

func printHeader() {
	fmt.Println()
	for i := 0; i < int(processesCount); i++ {
		fmt.Print(i)
	}
	fmt.Print(" | PID |         OP")
	fmt.Println()
	fmt.Println("---------+-----+-------------------")
}

func drawProcessesStates(isWorking []processState) any {
	res := ""
	for _, v := range isWorking {
		switch v {
		case IDLE:
			res = res + " "
		case STARTED:
			res = res + "▖"
		case WORKING:
			res = res + "▌"
		case FINISHED:
			res = res + "▘"
		}
	}
	return res
}
