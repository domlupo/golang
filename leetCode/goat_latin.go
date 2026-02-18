// https://leetcode.com/problems/goat-latin/

package main

func toGoatLatin(sentence string) string {
	isVowel := map[byte]bool{
		'a': true,
		'e': true,
		'i': true,
		'o': true,
		'u': true,
		'A': true,
		'E': true,
		'I': true,
		'O': true,
		'U': true,
	}

	var o []byte
	wordIndex := 1
	for i := 0; i < len(sentence); i++ {
		var b []byte
		for i < len(sentence) && sentence[i] != ' ' {
			b = append(b, sentence[i])
			i++
		}

		if !isVowel[b[0]] {
			b = append(b, b[0])
			b = b[1:]
		}

		b = append(b, 'm', 'a')

		for range wordIndex {
			b = append(b, 'a')
		}
		o = append(o, b...)
		o = append(o, ' ')
		wordIndex++
	}

	return string(o[:len(o)-1])
}
