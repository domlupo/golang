// https://leetcode.com/problems/flip-game/

package main

func generatePossibleNextMoves(currentState string) []string {
	var o []string

	for i := 0; i+1 < len(currentState); i++ {
		if currentState[i] == '+' && currentState[i+1] == '+' {
			o = append(o, currentState[:i]+"--"+currentState[i+2:])
		}
	}

	return o
}
