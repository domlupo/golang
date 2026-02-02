// https://leetcode.com/problems/reverse-bits

package main

func reverseBits1(n int) int {
	var a int
	for range 32 {
		a = a * 2  // shift left
		a += n % 2 // add the least significant bit of n to a
		n = n / 2  // shift n to the right
	}

	return a
}

func reverseBits2(n int) int {
	var a uint32
	for range 32 {
		a = a << 1           // shift left
		a += (uint32(n) & 1) // add the least significant bit of n to a
		n = n >> 1           // shift n to the right
	}

	return int(a)
}
