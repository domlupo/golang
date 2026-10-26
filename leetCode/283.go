// https://leetcode.com/problems/move-zeroes/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func moveZeroes(nums []int) {
	var z, n int
	for n < len(nums) {
		if n <= z || nums[n] == 0 {
			n++
		} else if nums[z] != 0 {
			z++
		} else {
			nums[n], nums[z] = nums[z], nums[n]
		}
	}
}
