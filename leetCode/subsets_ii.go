// https://leetcode.com/problems/subsets-ii/
// N is the length of nums
// time complexity is O(N * 2^N)
// N is from creating deep copies, and 2^N is due to subsets.
// space complexity is O(N)
// Other approaches have a space complexity of O(log(N)) but I skipped.
// I should look into again once better at LeetCode.

package main

func subsetsWithDup(nums []int) [][]int {
	for i := 0; i < len(nums)-1; i++ {
		if i >= 0 && nums[i] > nums[i+1] {
			nums[i], nums[i+1] = nums[i+1], nums[i]
			i -= 2
		}
	}

	var ans [][]int
	var comb []int
	backtrack90(&ans, nums, comb, 0)
	return ans
}

func backtrack90(ans *[][]int, nums, comb []int, i int) {
	if i == len(nums) {
		deepCopy := make([]int, len(comb))
		copy(deepCopy, comb)
		*ans = append(*ans, deepCopy)
		return
	}

	comb = append(comb, nums[i])
	backtrack90(ans, nums, comb, i+1)
	comb = comb[:len(comb)-1]
	for i < len(nums)-1 && nums[i] == nums[i+1] {
		i++
	}
	backtrack90(ans, nums, comb, i+1)
}

// Below is the subsets created without skipping duplicates
// []
// [1], []
// [1, 2], [1], [2], []
// [1, 2, 2], [1, 2], [1, 2], [1], [2, 2], [2], [2], []
// The extras are a [1, 2] and a [2], which need skipped once
// They are skipped by executing i++ when a duplicate occurs
