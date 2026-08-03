// https://leetcode.com/problems/course-schedule-ii/
// M is the number of courses, N is the number of prereqs
// time complexity is O(M + N)
// space complexity is O(M + N)

package main

func findOrder(numCourses int, prerequisites [][]int) []int {
	degrees := make([]int, numCourses)
	prereqs := make([][]int, numCourses)

	for _, pr := range prerequisites {
		degrees[pr[0]]++
		prereqs[pr[1]] = append(prereqs[pr[1]], pr[0])
	}

	var ans []int
	var q []int
	for course, degree := range degrees {
		if degree == 0 {
			q = append(q, course)
		}
	}

	for len(q) != 0 {
		course := q[0]
		q = q[1:]
		ans = append(ans, course)
		for _, pr := range prereqs[course] {
			degrees[pr]--
			if degrees[pr] == 0 {
				q = append(q, pr)
			}
		}
	}

	if len(ans) == numCourses {
		return ans
	} else {
		return []int{}
	}
}
