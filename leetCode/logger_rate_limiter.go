// https://leetcode.com/problems/logger-rate-limiter/

package main

type Logger struct {
	last map[string]int
}

func Constructor() Logger {
	return Logger{
		last: make(map[string]int),
	}
}

func (this *Logger) ShouldPrintMessage(timestamp int, message string) bool {
	l, ok := this.last[message]
	if !ok || (ok && timestamp >= l+10) {
		this.last[message] = timestamp
		return true
	}
	return false
}
