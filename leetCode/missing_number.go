// https://leetcode.com/problems/missing-number/

package main

func missingNumber(nums []int) int {
	var sum int
	var noMissing int
	for ni, n := range nums {
		sum += n
		noMissing += (ni + 1)
	}
	return noMissing - sum
}
