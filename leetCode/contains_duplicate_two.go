// https://leetcode.com/problems/contains-duplicate-ii/

package main

func containsNearbyDuplicate(nums []int, k int) bool {
	seen := make(map[int]int)

	for ni, n := range nums {
		if i, ok := seen[n]; ok && ni-i <= k {
			return true
		}
		seen[n] = ni
	}

	return false
}
