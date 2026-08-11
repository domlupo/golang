// https://leetcode.com/problems/trapping-rain-water/
// N is the length of height
// time complexity is O(N)
// space complexity is O(1)

package main

func trap(height []int) int {
	var ans, l, lMax, rMax int
	r := len(height) - 1

	for l < r {
		lMax = max(lMax, height[l])
		rMax = max(rMax, height[r])
		hTop := min(lMax, rMax)
		var hBottom int
		if height[l] < height[r] {
			hBottom = height[l]
			l++
		} else {
			hBottom = height[r]
			r--
		}
		ans += max(0, hTop-hBottom)
	}

	return ans
}
