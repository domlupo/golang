// https://leetcode.com/problems/fixed-point/

package main

func fixedPoint(arr []int) int {
	l := 0
	r := len(arr) - 1
	ans := -1

	for l <= r {
		mid := l + (r-l)/2
		if arr[mid] == mid {
			ans = mid
			r = mid - 1
		} else if arr[mid] < mid {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}

	return ans
}
