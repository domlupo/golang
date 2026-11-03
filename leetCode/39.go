// https://leetcode.com/problems/combination-sum/
// T is the target value and M is the min value in candidates
// time complexity is O(2^(T/M))
// space complexity is O(T/M), output space is larger but not included.

package main

func combinationSum(candidates []int, target int) [][]int {
	var ans [][]int
	backtrack39(&ans, []int{}, candidates, target)
	return ans
}

func backtrack39(ans *[][]int, comb, nums []int, remainder int) {
	if remainder == 0 {
		deepCopy := make([]int, len(comb))
		copy(deepCopy, comb)
		*ans = append(*ans, deepCopy)
		return
	} else if remainder < 0 {
		return
	}

	for ni, n := range nums {
		comb = append(comb, n)
		backtrack39(ans, comb, nums[ni:], remainder-n)
		comb = comb[:len(comb)-1]
	}
}
