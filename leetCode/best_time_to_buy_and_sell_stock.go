// https://leetcode.com/problems/best-time-to-buy-and-sell-stock/

package main

import "math"

func maxProfit(prices []int) int {
	m := math.MaxInt
	var profit int

	for _, p := range prices {
		profit = max(profit, p-m)
		m = min(m, p)
	}

	return profit
}
