// https://leetcode.com/problems/find-pivot-index/

package main

func pivotIndex(nums []int) int {
	var rSum int
	for _, n := range nums {
		rSum += n
	}

	var lSum int
	for i, n := range nums {
		rSum -= n
		if rSum == lSum {
			return i
		}
		lSum += n
	}

	return -1
}
