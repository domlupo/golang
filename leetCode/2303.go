// https://leetcode.com/problems/calculate-amount-paid-in-taxes/
// N is the length of brackets
// time complexity is O(N)
// space complexity is O(1)

package main

func calculateTax(brackets [][]int, income int) float64 {
	var ans, lower float64
	for _, b := range brackets {
		upper := min(float64(b[0]), float64(income))
		rate := float64(b[1])
		ans += (rate * (upper - lower)) / 100.0
		lower = upper
	}
	return ans
}
