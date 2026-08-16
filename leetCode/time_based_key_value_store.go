// https://leetcode.com/problems/time-based-key-value-store/
// M is the total number of keys, N is the number of data for a key
// time complexity is O(1) for set and O(log(N)) for get
// space complexity is O(M*N)

package main

type data struct {
	value     string
	timestamp int
}

type TimeMap struct {
	mp map[string][]data
}

func Constructor981() TimeMap {
	return TimeMap{
		mp: make(map[string][]data),
	}
}

func (this *TimeMap) Set(key string, value string, timestamp int) {
	d := data{value: value, timestamp: timestamp}
	this.mp[key] = append(this.mp[key], d)
}

func (this *TimeMap) Get(key string, timestamp int) string {
	vals, ok := this.mp[key]
	if !ok || vals[0].timestamp > timestamp {
		return ""
	}

	var l int
	r := len(vals) - 1
	for l <= r {
		m := (l + r) / 2
		ts := vals[m].timestamp
		if ts > timestamp {
			r = m - 1
		} else {
			l = m + 1
		}
	}

	return vals[r].value
}
