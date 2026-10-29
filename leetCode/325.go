// https://leetcode.com/problems/maximum-size-subarray-sum-equals-k/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(N)

package main

func maxSubArrayLen(nums []int, k int) int {
	sumLeftIndices := make(map[int]int)
	sumLeftIndices[0] = -1

	var prefixSum, ans int
	for i, n := range nums {
		prefixSum += n
		if left, ok := sumLeftIndices[prefixSum-k]; ok {
			ans = max(ans, i-left)
		}
		if _, ok := sumLeftIndices[prefixSum]; !ok {
			sumLeftIndices[prefixSum] = i
		}
	}

	return ans
}
