// https://leetcode.com/problems/first-bad-version/

package main

func firstBadVersion(n int) int {
	l := 1
	r := n

	for l < r {
		mid := l + (r-l)/2
		if isBadVersion(mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}

	return l
}

// provided
func isBadVersion(version int) bool {
	return true
}
