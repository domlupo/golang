// https://leetcode.com/problems/kth-missing-positive-number/
// N is the length of arr
// time complexity is O(log(N))
// space complexity is O(1)

package main

func findKthPositive(arr []int, k int) int {
	l, r := 0, len(arr)-1
	for l <= r {
		m := (l + r) / 2
		if arr[m]-m-1 < k {
			l = m + 1
		} else {
			r = m - 1
		}
	}
	return l + k
}

// arr:     2 3 4 7 11
// m:       0 1 2 3  4
// arr[m]-m 1 1 1 3  6
