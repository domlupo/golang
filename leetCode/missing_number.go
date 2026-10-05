// https://leetcode.com/problems/missing-number/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func missingNumber(nums []int) int {
	ans := len(nums)
	for i := range nums {
		ans ^= i
		ans ^= nums[i]
	}
	return ans
}
