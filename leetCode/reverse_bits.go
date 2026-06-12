// https://leetcode.com/problems/reverse-bits

package main

func reverseBits1(n int) int {
	var a int
	for range 32 {
		a *= 2     // shift left
		a += n % 2 // add the least significant bit of n to a
		n /= 2     // shift n to the right
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

// input  = 7  0111
// output = 14 1110
//
// 0 = 0 * 2   0000
// 1 = 0 + 1   0001
// 3 = 7 / 2   0011
//
// 2 = 2 * 1   0010
// 3 = 2 + 1   0011
// 1 = 3 / 2   0001
//
// 6 = 3 * 2   0110
// 7 = 6 + 1   0111
// 0 = 1 / 2   0000
//
// 14 = 7 * 2  1110
// 14 = 14 + 0 1110
// 0 = 0 / 2   0000
