// https://leetcode.com/problems/count-primes/
// N is the passed n
// time complexity is O(log(log(N))
// space complexity is O(n)

package main

func countPrimes(n int) int {
	var ans int
	prime := make([]bool, n)

	for i := range n {
		prime[i] = true
	}

	for i := 2; i < n; i++ {
		if prime[i] {
			ans++
			for j := 2 * i; j < n; j += i {
				prime[j] = false
			}
		}
	}
	return ans
}
