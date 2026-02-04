// https://leetcode.com/problems/intersection-of-two-arrays/

package main

func intersection(nums1 []int, nums2 []int) []int {
	contains := make(map[int]bool)
	pushed := make(map[int]bool)

	for _, n := range nums1 {
		contains[n] = true
	}

	i := []int{}
	for _, n := range nums2 {
		_, cOk := contains[n]
		_, pOk := pushed[n]
		if cOk && !pOk {
			i = append(i, n)
			pushed[n] = true
		}
	}
	return i
}
