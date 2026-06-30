// https://leetcode.com/problems/container-with-most-water/

package main

func maxArea(height []int) int {
	var ans, l int
	r := len(height) - 1

	for l < r {
		w := r - l
		h := min(height[l], height[r])
		ans = max(ans, w*h)
		if height[l] < height[r] {
			l++
		} else {
			r--
		}
	}

	return ans
}
