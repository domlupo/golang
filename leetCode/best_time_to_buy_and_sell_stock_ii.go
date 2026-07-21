// https://leetcode.com/problems/best-time-to-buy-and-sell-stock-ii/
// N is the length of prices
// time complexity is O(N)
// space complexity is O(1)

package main

func maxProfit122(prices []int) int {
	var total int
	for i := 0; i < len(prices)-1; i++ {
		total += max(0, prices[i+1]-prices[i])
	}
	return total
}
