// https://leetcode.com/problems/rotate-array/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func rotate189(nums []int, k int) {
	reverse189(nums, 0, len(nums)-1)
	k = k % len(nums)
	reverse189(nums, 0, k-1)
	reverse189(nums, k, len(nums)-1)
}

func reverse189(nums []int, l, r int) {
	for l < r {
		nums[l], nums[r] = nums[r], nums[l]
		l++
		r--
	}
}
