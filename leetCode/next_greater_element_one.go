// https://leetcode.com/problems/next-greater-element-i/

package main

func nextGreaterElement(nums1 []int, nums2 []int) []int {
	var s []int
	m := make(map[int]int)

	for _, n := range nums2 {
		if m[n] == 0 {
			m[n] = -1
		}

		for len(s) > 0 && s[len(s)-1] < n {
			m[s[len(s)-1]] = n
			s = s[:len(s)-1]
		}
		s = append(s, n)
	}

	var o []int
	for _, n := range nums1 {
		o = append(o, m[n])
	}
	return o
}
