// https://leetcode.com/problems/remove-duplicates-from-sorted-array
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func removeDuplicates(nums []int) int {
	var l int
	for r := 1; r < len(nums); r++ {
		if nums[l] != nums[r] {
			l++
			nums[l] = nums[r]
		}
	}
	return l + 1
}
