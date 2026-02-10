// https://leetcode.com/problems/set-mismatch/

package main

// o(1) space is complex bit manipulation. Will revisit.

// o(n) space
func findErrorNums(nums []int) []int {
	counts := make([]int, len(nums)+1)
	var duplicate, lowSum, fullSum int

	for ni, n := range nums {
		counts[n] = counts[n] + 1
		if counts[n] == 2 {
			duplicate = n
			lowSum -= n
		}
		lowSum += n
		fullSum += (ni + 1)
	}

	return []int{duplicate, fullSum - lowSum}
}
