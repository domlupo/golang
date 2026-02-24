// https://leetcode.com/problems/fair-candy-swap/

package main

func fairCandySwap(aliceSizes []int, bobSizes []int) []int {
	var aliceSum, bobSum int
	aliceSwaps := make(map[int]bool)

	for _, size := range aliceSizes {
		aliceSum += size
		aliceSwaps[size] = true
	}

	for _, size := range bobSizes {
		bobSum += size
	}

	for _, bobSize := range bobSizes {
		swapNeeded := (aliceSum - bobSum + (2 * bobSize)) / 2
		if aliceSwaps[swapNeeded] {
			return []int{swapNeeded, bobSize}
		}
	}

	panic("no solution found")
}
