// https://leetcode.com/problems/permutations/
// time complexity is O(n * n!)
// space complexity is O(n)

package main

func permute(nums []int) [][]int {
	var ans [][]int
	var curr []int
	backtrack3(nums, curr, &ans)
	return ans
}

func backtrack3(nums, curr []int, ans *[][]int) {
	if len(curr) == len(nums) {
		deepCopy := make([]int, len(curr))
		copy(deepCopy, curr)
		*ans = append(*ans, deepCopy)
		return
	}

	for _, n := range nums {
		var exists bool
		for _, v := range curr {
			if v == n {
				exists = true
				break
			}
		}

		if !exists {
			curr = append(curr, n)
			backtrack3(nums, curr, ans)
			curr = curr[:len(curr)-1]
		}
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
