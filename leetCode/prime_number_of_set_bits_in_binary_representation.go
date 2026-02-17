// https://leetcode.com/problems/prime-number-of-set-bits-in-binary-representation/

package main

func countPrimeSetBits(left int, right int) int {
	prime := map[int]bool{
		2:  true,
		3:  true,
		5:  true,
		7:  true,
		11: true,
		13: true,
		17: true,
		19: true,
	}

	var primeSetBits int
	for i := left; i <= right; i++ {
		if prime[countSetBits1(i)] {
			primeSetBits++
		}
	}
	return primeSetBits
}

func countSetBits1(n int) int {
	var count int
	for n > 0 {
		n &= n - 1
		count++
	}
	return count
}
