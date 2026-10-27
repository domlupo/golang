// https://leetcode.com/problems/decrypt-string-from-alphabet-to-integer-mapping/
// N is the length of s
// time complexity is O(N)
// space complexity is O(1)

package main

func freqAlphabets(s string) string {
	numToByte := map[string]byte{
		"1":  'a',
		"2":  'b',
		"3":  'c',
		"4":  'd',
		"5":  'e',
		"6":  'f',
		"7":  'g',
		"8":  'h',
		"9":  'i',
		"10": 'j',
		"11": 'k',
		"12": 'l',
		"13": 'm',
		"14": 'n',
		"15": 'o',
		"16": 'p',
		"17": 'q',
		"18": 'r',
		"19": 's',
		"20": 't',
		"21": 'u',
		"22": 'v',
		"23": 'w',
		"24": 'x',
		"25": 'y',
		"26": 'z',
	}
	var bytes []byte
	var i int

	for i < len(s) {
		if i+2 < len(s) && s[i+2] == '#' {
			b := numToByte[s[i:i+2]]
			bytes = append(bytes, b)
			i += 3
		} else {
			b := numToByte[string(s[i])]
			bytes = append(bytes, b)
			i++
		}
	}

	return string(bytes)
}
