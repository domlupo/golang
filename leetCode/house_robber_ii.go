// https://leetcode.com/problems/house-robber-ii/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func rob213(nums []int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	return max(robHouses(nums[1:]), robHouses(nums[:len(nums)-1]))
}

func robHouses(nums []int) int {
	var prev, curr int
	for _, n := range nums {
		prev, curr = curr, max(curr, prev+n)
	}
	return curr
}

// 1 2 3 1
// prev, curr = 0, 0
// prev, curr = 0, max(0, 0+1)
// prev, curr = 1, max(1, 0+2)
// prev, curr = 2, max(2, 1+3)
// prev, curr = 4, max(4, 2+1
