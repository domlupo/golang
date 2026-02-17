// https://leetcode.com/problems/find-anagram-mappings/

package main

func anagramMappings(nums1 []int, nums2 []int) []int {
	numIndexes := make(map[int]int)
	o := make([]int, len(nums1))

	for i, n := range nums2 {
		numIndexes[n] = i
	}

	for i, n := range nums1 {
		o[i] = numIndexes[n]
	}

	return o
}
