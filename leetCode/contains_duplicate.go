// https://leetcode.com/problems/contains-duplicate/

package main

func containsDuplicate(nums []int) bool {
	seen := make(map[int]bool)

	for _, n := range nums {
		_, ok := seen[n]
		if ok {
			return true
		}
		seen[n] = true
	}

	return false
}
