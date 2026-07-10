// https://leetcode.com/problems/subsets/
// time complexity is complex
// space complexity is complex

package main

func subsets(nums []int) [][]int {
	var ans [][]int
	var comb []int
	backtrack78(&ans, comb, nums, 0)
	return ans
}

func backtrack78(ans *[][]int, comb, nums []int, index int) {
	if len(nums) == index {
		deepCopy := make([]int, len(comb))
		copy(deepCopy, comb)
		*ans = append(*ans, deepCopy)
		return
	}

	comb = append(comb, nums[index])
	backtrack78(ans, comb, nums, index+1)
	comb = comb[:len(comb)-1]
	backtrack78(ans, comb, nums, index+1)
}

// nums = [1, 2, 3]
// comb index
// [] 0
//
// [1] 1
// [] 1
//
// [1, 2] 2
// [1] 2
// [2] 2
// [] 2
//
// [1, 2, 3] 3
// [1, 2] 3
// [1, 3] 3
// [1] 3
// [2, 3] 3
// [2] 3
// [3] 3
// [] 3
