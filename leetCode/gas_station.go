// https://leetcode.com/problems/gas-station/
// N is the length of gas
// time complexity is O(N)
// space complexity is O(1)

package main

import "math"

func canCompleteCircuit(gas []int, cost []int) int {
	var sum, start int
	minimum := math.MaxInt

	for i := 0; i < len(gas); i++ {
		sum += (gas[i] - cost[i])
		minimum = min(sum, minimum)
		if sum == minimum {
			start = i + 1
		}
	}

	if sum < 0 {
		return -1
	} else if start == len(gas) {
		return 0
	} else {
		return start
	}
}
