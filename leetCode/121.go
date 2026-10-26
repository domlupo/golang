// https://leetcode.com/problems/best-time-to-buy-and-sell-stock/
// N is the length of prices
// time complexity is O(N)
// space complexity is O(1)

package main

func maxProfit(prices []int) int {
	cost := prices[0]
	var profit int

	for _, price := range prices {
		profit = max(profit, price-cost)
		cost = min(cost, price)
	}

	return profit
}
