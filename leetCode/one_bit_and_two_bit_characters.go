// https://leetcode.com/problems/1-bit-and-2-bit-characters/

package main

func isOneBitCharacter(bits []int) bool {
	var i int
	for i < len(bits)-1 {
		if bits[i] == 0 {
			i += 1
		} else {
			i += 2
		}
	}

	return i == len(bits)-1
}
