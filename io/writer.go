package io

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var (
	logFiles       []*os.File
	writers        []*bufio.Writer
	logPath              = "execution-"
	processesCount int64 = 8
)

// Init initializes n writers from 0 to n-1. writers[i] writes to logFile[i].
// A process(goroutine) should only call one writer to have concurrency.
func Init(n int64) {
	processesCount = n
	logFiles = make([]*os.File, processesCount)
	writers = make([]*bufio.Writer, processesCount)
	for pid := 0; pid < int(processesCount); pid++ {
		writers[pid] = writer(int64(pid))
	}
}

// CloseWriters closes all logFiles.
func CloseWriters() {
	var pid int64 = 0
	for pid < processesCount {
		CloseWriter(pid)
		pid++
	}
}

// CloseWriter closes writers[i]. It should be called after writing logs of process pid has finished.
// If not, then nothing will be written to opened logFiles.
func CloseWriter(pid int64) {
	writers[pid].Flush()
	logFiles[pid].Close()
}

func write(str string, pid int64) {
	writers[pid].WriteString(fmt.Sprint(time.Now().UnixMicro()) + str + "\n")
}

func writer(pid int64) *bufio.Writer {
	file := fmt.Sprintf("%s%d.txt", logPath, pid)
	logFiles[pid], _ = os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	return bufio.NewWriter(logFiles[pid])
}
