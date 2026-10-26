// https://leetcode.com/problems/happy-number/
// N is the value of n
// time complexity is O(log(N))
// space complexity is O(1)

package main

func isHappy(n int) bool {
	slow, fast := n, sumOfDigitsSquared(n)
	for slow != fast {
		slow = sumOfDigitsSquared(slow)
		fast = sumOfDigitsSquared(fast)
		fast = sumOfDigitsSquared(fast)
	}
	return slow == 1
}

func sumOfDigitsSquared(n int) int {
	var sum int
	for n != 0 {
		digit := n % 10
		n /= 10
		sum += (digit * digit)
	}
	return sum
}
