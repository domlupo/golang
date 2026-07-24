// https://leetcode.com/problems/maximum-product-subarray/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

import "math"

func maxProduct(nums []int) int {
	product := 1
	maxProduct := math.MinInt

	for i := range nums {
		product *= nums[i]
		maxProduct = max(maxProduct, product)
		if product == 0 {
			product = 1
		}
	}

	product = 1
	for i := len(nums) - 1; i >= 0; i-- {
		product *= nums[i]
		maxProduct = max(maxProduct, product)
		if product == 0 {
			product = 1
		}
	}

	return maxProduct
}
