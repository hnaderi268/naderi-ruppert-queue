package io

import "fmt"

type opState int64

const (
	START_ENQ opState = iota
	START_DEQ
	END_ENQ
	END_DEQ
)

// Assumption: There will be no call of Log(i,...), until Log(i,...) terminates. This happens if
// Log(i,...) is only called by one process.
// Log writes a log to logFiles[pid] containing v with format opState.
func Log(pid int64, opState opState, v ...interface{}) {
	switch opState {
	case START_ENQ:
		write(fmt.Sprint(startf(pid), " ", v[0], " -> ▢▢▢"), pid)
	case START_DEQ:
		write(fmt.Sprint(startf(pid)), pid)
	case END_ENQ:
		write(fmt.Sprint(endf(pid)), pid)
	case END_DEQ:
		write(fmt.Sprint(endf(pid), "    ▢▢▢ -> ", v[0]), pid)
	}
}

func startf(pid int64) string {
	return fmt.Sprint("[", pid, "  | ")
}

func endf(pid int64) string {
	return fmt.Sprint(" ", pid, "] |", "\t")
}
