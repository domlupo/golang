// https://leetcode.com/problems/intersection-of-two-arrays-ii/

package main

func intersect(nums1 []int, nums2 []int) []int {
	counts := make(map[int]int)

	for _, n := range nums1 {
		counts[n] = counts[n] + 1
	}

	var i []int
	for _, n := range nums2 {
		if c, ok := counts[n]; ok && c > 0 {
			i = append(i, n)
			counts[n] = counts[n] - 1
		}
	}

	return i
}
