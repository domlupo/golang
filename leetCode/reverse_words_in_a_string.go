// https://leetcode.com/problems/reverse-words-in-a-string/
// N is the length of s
// time complexity is O(N)
// space complexity is O(N) because Golang does not allow string mutation.

package main

func reverseWords151(s string) string {
	var letters []byte

	// trim whitespace
	s = s + " "
	var addWhitespace bool
	for _, c := range s {
		if c != ' ' {
			letters = append(letters, byte(c))
			addWhitespace = true
		} else if addWhitespace {
			letters = append(letters, ' ')
			addWhitespace = false
		}
	}
	letters = letters[:len(letters)-1]

	// reverse  entire string
	l := 0
	r := len(letters) - 1
	for l < r {
		letters[l], letters[r] = letters[r], letters[l]
		l++
		r--
	}

	// reverse words
	l = 0
	for l < len(letters) {
		r = l
		for r < len(letters) && letters[r] != ' ' {
			r++
		}
		next := r + 1
		r--
		for l < r {
			letters[l], letters[r] = letters[r], letters[l]
			l++
			r--
		}
		l = next
	}

	return string(letters)
}
