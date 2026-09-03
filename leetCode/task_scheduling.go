// https://neetcode.io/problems/task-scheduling/
// N is the number of tasks
// time complexity is O(N)
// space complexity is O(1)

package main

func leastInterval(tasks []byte, n int) int {
	freqs := make([]int, 26)
	var maxFreq int

	for _, t := range tasks {
		freqs[t-byte('A')]++
		maxFreq = max(maxFreq, freqs[t-byte('A')])
	}

	var maxFreqCount int
	for _, f := range freqs {
		if f == maxFreq {
			maxFreqCount++
		}
	}

	time := (maxFreq-1)*(n+1) + maxFreqCount
	return max(len(tasks), time)
}
