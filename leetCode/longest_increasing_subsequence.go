// https://leetcode.com/problems/longest-increasing-subsequence/
// N is the length of nums
// time complexity is O(log(N)*N)
// space complexity is O(N)

package main

func lengthOfLIS(nums []int) int {
	sub := []int{nums[0]}

	for i := 1; i < len(nums); i++ {
		if nums[i] > sub[len(sub)-1] {
			sub = append(sub, nums[i])
		} else {
			pos := search300(sub, nums[i])
			sub[pos] = nums[i]
		}
	}

	return len(sub)
}

func search300(nums []int, n int) int {
	l, r := 0, len(nums)-1
	for l < r {
		m := (l + r) / 2
		if nums[m] < n {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}
