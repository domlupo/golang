// https://leetcode.com/problems/course-schedule/
// M is prerequisites, N is the numCourses
// time complexity is O(M + N)
// space complexity is O(M + N)

package main

func canFinish(numCourses int, prerequisites [][]int) bool {
	degrees := make([]int, numCourses)
	prereqs := make([][]int, numCourses)

	for _, pr := range prerequisites {
		prereqs[pr[1]] = append(prereqs[pr[1]], pr[0])
		degrees[pr[0]]++
	}

	var tookCourse int
	var queue []int
	for course, degree := range degrees {
		if degree == 0 {
			queue = append(queue, course)
		}
	}

	for len(queue) != 0 {
		course := queue[0]
		queue = queue[1:]
		tookCourse++
		for _, requires := range prereqs[course] {
			degrees[requires]--
			if degrees[requires] == 0 {
				queue = append(queue, requires)
			}
		}
	}

	return tookCourse == numCourses
}
