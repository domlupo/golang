// https://leetcode.com/problems/two-sum-ii-input-array-is-sorted/
// time complexity is O(N)
// space complexity is O(1)

package main

func twoSum167(numbers []int, target int) []int {
	var l int
	r := len(numbers) - 1
	for l < r {
		sum := numbers[l] + numbers[r]
		if sum == target {
			break
		} else if sum < target {
			l++
		} else {
			r--
		}
	}

	return []int{l + 1, r + 1}
}
