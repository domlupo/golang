// https://leetcode.com/problems/sentence-similarity/

package main

func areSentencesSimilar(sentence1 []string, sentence2 []string, similarPairs [][]string) bool {
	if len(sentence1) != len(sentence2) {
		return false
	}

	pairs := make(map[[2]string]bool)

	for _, p := range similarPairs {
		pairs[[2]string{p[0], p[1]}] = true
		pairs[[2]string{p[1], p[0]}] = true
	}

	for i := 0; i < len(sentence1); i++ {
		if sentence1[i] != sentence2[i] && !pairs[[2]string{sentence1[i], sentence2[i]}] {
			return false
		}
	}

	return true
}
