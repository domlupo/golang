// https://leetcode.com/problems/contains-duplicate-ii/

package main

func containsNearbyDuplicate(nums []int, k int) bool {
	lastSeen := make(map[int]int)

	for i, n := range nums {
		ls, ok := lastSeen[n]
		if ok && i-ls <= k {
			return true
		}
		lastSeen[n] = i
	}

	return false
}
