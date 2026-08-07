// https://leetcode.com/problems/encode-and-decode-strings/
// N is the number of characters in passed strs, K is the number of strings
// time complexity is O(N)
// space complexity is O(K) due to adding delimiter

package main

type Codec struct{}

// Encodes a list of strings to a single string.
func (codec *Codec) Encode(strs []string) string {
	var bytes []byte
	for _, str := range strs {
		for _, c := range str {
			if c == '/' || c == '.' {
				bytes = append(bytes, '/')
			}
			bytes = append(bytes, byte(c))
		}
		bytes = append(bytes, '.')
	}
	return string(bytes)
}

// Decodes a single string to a list of strings.
func (codec *Codec) Decode(strs string) []string {
	var ans []string
	var bytes []byte
	for i := 0; i < len(strs); i++ {
		if i < len(strs)-1 && strs[i:i+2] == "//" {
			bytes = append(bytes, '/')
			i++
		} else if i < len(strs)-1 && strs[i:i+2] == "/." {
			bytes = append(bytes, '.')
			i++
		} else if strs[i] == '.' {
			ans = append(ans, string(bytes))
			bytes = []byte{}
		} else {
			bytes = append(bytes, strs[i])
		}
	}
	return ans
}
