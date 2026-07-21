// https://leetcode.com/problems/longest-consecutive-sequence/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(N)

package main

func longestConsecutive(nums []int) int {
	numbers := make(map[int]bool)

	for _, n := range nums {
		numbers[n] = true
	}

	var longest int
	for n := range numbers {
		if numbers[n-1] {
			continue
		}
		var count int
		for numbers[n] {
			n++
			count++
		}
		longest = max(longest, count)
	}

	return longest
}
