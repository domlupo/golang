// https://leetcode.com/problems/house-robber/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func rob(nums []int) int {
	var prev, curr int
	for _, n := range nums {
		prev, curr = curr, max(prev+n, curr)
	}
	return curr
}

// [1, 2, 3, 1]

// i = 0 -> 1
// i = 1 -> 1 or 2
// i = 2 -> (1 + 3) or max at prior index
// i = 3 -> (1 + 3) or (1 + 2) or (1 + 1)
// i = 3 -> max at prior index or curr index + max at two index back

// i = 0 -> prev, curr = 0, 1
// i = 1 -> prev, curr = 1, max(1, 0+2)
// i = 2 -> prev, curr = 2, max(2, 1+3)
// i = 3 -> prev, curr = 4, max(4, 2+1)
