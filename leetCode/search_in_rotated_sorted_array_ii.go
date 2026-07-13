// https://leetcode.com/problems/search-in-rotated-sorted-array-ii
// N is the length of nums
// time complexity is O(N) in worst case, O(log(N)) in best case.
// space complexity is O(1)

package main

// This solution is not 100% optimal but very easy to understand.
// Once I get better at LeetCode I can try more complex solution.
// I believe three sum is related.
func search81(nums []int, target int) bool {
	var l, m int
	r := len(nums) - 1
	for l <= r {
		m = (l + r) / 2
		if nums[m] == target || nums[l] == target || nums[r] == target {
			return true
		} else if within81(target, nums[l], nums[m]) {
			r = m - 1
		} else if within81(target, nums[m], nums[r]) {
			l = m + 1
		} else if nums[l] == nums[m] {
			l++
		} else {
			r--
		}
	}

	return false
}

func within81(m, l, r int) bool {
	return m > l && m < r
}
