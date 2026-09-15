// https://leetcode.com/problems/subsets-ii/
// N is the length of nums
// time complexity is O(N * 2^N)
// N is from creating deep copies, and 2^N is due to subsets.
// space complexity is O(N)

package main

import "slices"

func subsetsWithDup(nums []int) [][]int {
	var ans [][]int
	slices.Sort(nums)
	backtrack90(&ans, []int{}, nums, 0)
	return ans
}

func backtrack90(ans *[][]int, comb, nums []int, i int) {
	if i == len(nums) {
		deepCopy := make([]int, len(comb))
		copy(deepCopy, comb)
		*ans = append(*ans, deepCopy)
		return
	}

	comb = append(comb, nums[i])
	backtrack90(ans, comb, nums, i+1)
	comb = comb[:len(comb)-1]
	for i+1 < len(nums) && nums[i] == nums[i+1] {
		i++
	}
	backtrack90(ans, comb, nums, i+1)
}

// Below is the subsets created without skipping duplicates
// []
// [1], []
// [1, 2], [1], [2], []
// [1, 2, 2], [1, 2], [1, 2], [1], [2, 2], [2], [2], []
// The extras are a [1, 2] and a [2], which need skipped once
// They are skipped by executing i++ when a duplicate occurs
