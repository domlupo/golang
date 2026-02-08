// https://leetcode.com/problems/range-addition-ii/

package main

import "math"

func maxCount(m int, n int, ops [][]int) int {
	if len(ops) == 0 {
		return m * n
	}

	minM := math.MaxInt
	minN := math.MaxInt
	for _, o := range ops {
		minM = min(minM, o[0])
		minN = min(minN, o[1])
	}
	return minM * minN
}
