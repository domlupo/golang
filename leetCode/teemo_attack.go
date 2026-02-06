// https://leetcode.com/problems/teemo-attacking/

package main

func findPoisonedDuration(timeSeries []int, duration int) int {
	var d int
	for i := 0; i+1 < len(timeSeries); i++ {
		d += min(timeSeries[i+1]-timeSeries[i], duration)
	}
	return d + duration
}
