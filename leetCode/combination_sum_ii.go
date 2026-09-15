// https://leetcode.com/problems/combination-sum-ii/
// N is the length of candidates
// time complexity is O(2^N)
// space complexity is O(N)

package main

import "slices"

func combinationSum2(candidates []int, target int) [][]int {
	slices.Sort(candidates)
	var ans [][]int
	backtrack40(&ans, []int{}, candidates, target)
	return ans
}

func backtrack40(ans *[][]int, comb, nums []int, remainder int) {
	if remainder == 0 {
		deepCopy := make([]int, len(comb))
		copy(deepCopy, comb)
		*ans = append(*ans, deepCopy)
		return
	} else if remainder < 0 {
		return
	}

	for i := 0; i < len(nums); i++ {
		comb = append(comb, nums[i])
		backtrack40(ans, comb, nums[i+1:], remainder-nums[i])
		comb = comb[:len(comb)-1]
		for i+1 < len(nums) && nums[i] == nums[i+1] {
			i++
		}
	}
}
