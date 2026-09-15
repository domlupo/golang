// https://leetcode.com/problems/permutations/
// time complexity is O(n * n!)
// space complexity is O(n)

package main

import "slices"

func permute(nums []int) [][]int {
	var ans [][]int
	backtrack46(&ans, []int{}, nums)
	return ans
}

func backtrack46(ans *[][]int, comb, nums []int) {
	if len(comb) == len(nums) {
		deepCopy := make([]int, len(comb))
		copy(deepCopy, comb)
		*ans = append(*ans, deepCopy)
		return
	}

	for _, n := range nums {
		if slices.Contains(comb, n) {
			continue
		}

		comb = append(comb, n)
		backtrack46(ans, comb, nums)
		comb = comb[:len(comb)-1]
	}
}

// 3 numbers is 6 permutations
// 4 numbers is 24 permutations

// 1 2 3
// 1 3 2
// 2 1 3
// 2 3 1
// 3 1 2
// 3 2 1

// 1 2 3 4
// 1 2 4 3
// 1 3 2 4
// 1 3 4 2
// 1 4 2 3
// 1 4 3
