// https://leetcode.com/problems/product-of-array-except-self/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func productExceptSelf(nums []int) []int {
	var ans []int
	product := 1
	for _, n := range nums {
		ans = append(ans, product)
		product *= n
	}
	product = 1
	for i := len(nums) - 1; i >= 0; i-- {
		ans[i] = ans[i] * product
		product *= nums[i]
	}
	return ans
}
