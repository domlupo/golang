// https://leetcode.com/problems/maximum-product-subarray/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

import "math"

func maxProduct(nums []int) int {
	ans := math.MinInt

	var prefix, suffix int
	for i := 0; i < len(nums); i++ {
		if prefix == 0 {
			prefix = 1
		}
		if suffix == 0 {
			suffix = 1
		}
		prefix *= nums[i]
		suffix *= nums[len(nums)-1-i]
		ans = max(ans, prefix, suffix)
	}

	return ans
}
