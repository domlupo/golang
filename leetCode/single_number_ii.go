// https://leetcode.com/problems/single-number-ii/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func singleNumber137(nums []int) int {
	var loner int
	for shift := range 32 {
		var bitSum int
		for _, n := range nums {
			bitSum += (n >> shift) & 1
		}
		// compute the bit of loner and place it
		lonerBit := bitSum % 3
		loner = loner | (lonerBit << shift)
	}
	// do not mistaken sign bit for MSB.
	negative := (loner & (1 << 31)) == 1<<31
	if negative {
		loner = loner - (1 << 32)
	}
	return loner
}
