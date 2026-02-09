// https://leetcode.com/problems/minimum-index-sum-of-two-lists/

package main

import "math"

func findRestaurant(list1 []string, list2 []string) []string {
	l := make(map[string]int)

	for i, s := range list1 {
		l[s] = i
	}

	m := math.MaxInt
	var o []string
	for i2, s := range list2 {
		i1, ok := l[s]
		sum := i1 + i2
		if ok && sum < m {
			o = []string{s}
			m = sum
		} else if ok && sum == m {
			o = append(o, s)
		}
	}

	return o
}
