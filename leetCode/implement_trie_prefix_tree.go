// https://leetcode.com/problems/implement-trie-prefix-tree/
// time complexity is O(N)
// space complexity is O(N) for Insert and O(1) for Search and StartsWith

package main

type Trie struct {
	links []*Trie
	end   bool
}

func Constructor208() Trie {
	return Trie{
		links: make([]*Trie, 26),
	}
}

func (this *Trie) Insert(word string) {
	curr := this
	for _, c := range word {
		if curr.links[c-'a'] == nil {
			curr.links[c-'a'] = &Trie{
				links: make([]*Trie, 26),
			}
		}
		curr = curr.links[c-'a']
	}
	curr.end = true
}

func (this *Trie) Search(word string) bool {
	curr := this
	for _, c := range word {
		if curr.links[c-'a'] == nil {
			return false
		}
		curr = curr.links[c-'a']
	}
	return curr.end == true
}

func (this *Trie) StartsWith(prefix string) bool {
	curr := this
	for _, c := range prefix {
		if curr.links[c-'a'] == nil {
			return false
		}
		curr = curr.links[c-'a']
	}
	return true
}
