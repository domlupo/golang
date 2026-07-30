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

// 0 = max(0, nums[n], 0)
// 1 = max(nums[n-1], nums[n])
// 2 = max(mem[n-1], mem[n-2] + nums[n])
// 3 = max(mem[n-1], mem[n-2] + nums[n]
