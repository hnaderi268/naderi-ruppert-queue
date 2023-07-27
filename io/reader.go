package io

import (
	"os"
)

import (
	"bufio"
	"fmt"
)

// ReadAll returns a slice containing all the lines written to logFiles.
func ReadAll() []string {
	output := make([]string, 0)
	var pid int64 = 0
	for pid < processesCount {
		readerPID := reader(pid)
		for readerPID.Scan() {
			output = append(output, readerPID.Text())
		}
		pid++
	}
	return output
}

func reader(pid int64) *bufio.Scanner {
	file := fmt.Sprintf("%s%d.txt", logPath, pid)
	logFile, _ := os.Open(file)
	return bufio.NewScanner(logFile)
}
