// https://leetcode.com/problems/student-attendance-record-i/

package main

func checkRecord(s string) bool {
	var absentCount, lateCount int
	for _, c := range s {
		switch c {
		case 'A':
			absentCount++
			lateCount = 0
		case 'P':
			lateCount = 0
		case 'L':
			lateCount++
		}
		if absentCount == 2 || lateCount == 3 {
			return false
		}
	}

	return true
}
