// https://leetcode.com/problems/subsets/
// time complexity is O(n * 2^n)
// space complexity O(2^n)

package main

func subsets(nums []int) [][]int {
	var ans [][]int
	backtrack78(&ans, []int{}, nums, 0)
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

// 1
// _

// -----------------------

// 1 2
// 2

// 2
// _

// -----------------------

// 1 2 3
// 1 2

// 2 3
// 2

// 2 3
// 3

// 3
// _

// -----------------------
