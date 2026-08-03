// https://leetcode.com/problems/design-add-and-search-words-data-structure/
// N is the length of word
// time complexity is O(N)
// space complexity is O(N)
// O(N) because wildcards are limited to 2.
// If wildcards were not limited, time and space complexity would be 26^N.

package main

type WordDictionary struct {
	Next []*WordDictionary
	End  bool
}

func Constructor211() WordDictionary {
	return WordDictionary{
		Next: make([]*WordDictionary, 26),
	}
}

func (this *WordDictionary) AddWord(word string) {
	curr := this
	for _, c := range word {
		if curr.Next[c-'a'] == nil {
			curr.Next[c-'a'] = &WordDictionary{
				Next: make([]*WordDictionary, 26),
			}
		}
		curr = curr.Next[c-'a']
	}
	curr.End = true
}

func (this *WordDictionary) Search(word string) bool {
	curr := this
	for i, c := range word {
		if c == '.' {
			for _, n := range curr.Next {
				if n != nil && n.Search(word[i+1:]) {
					return true
				}
			}
			// no letter available to continue search
			return false
		} else if curr.Next[c-'a'] != nil {
			curr = curr.Next[c-'a']
		} else {
			return false
		}
	}

	// edge case: '.' last char returns false unless the specific letter is an end
	return curr.End
}
