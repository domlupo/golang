// https://leetcode.com/problems/find-the-difference/

package main

// bit manipulation
func findTheDifference(s string, t string) byte {
	var sum byte
	for _, l := range t {
		sum ^= byte(l)
	}
	for _, l := range s {
		sum ^= byte(l)
	}
	return sum
}

// simple math
func findTheDifference2(s string, t string) byte {
	var sum byte
	for _, l := range t {
		sum += byte(l)
	}
	for _, l := range s {
		sum -= byte(l)
	}
	return sum
}
