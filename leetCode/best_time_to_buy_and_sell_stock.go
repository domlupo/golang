// https://leetcode.com/problems/best-time-to-buy-and-sell-stock/

package main

func maxProfit(prices []int) int {
	m := prices[0]
	var profit int

	for _, p := range prices {
		profit = max(profit, p-m)
		m = min(m, p)
	}

	return profit
}
