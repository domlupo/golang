// hard coded test setups which help manual testing
package main

type setup int

const (
	horizontal = iota
	horizontalMissile
	vertical
	diagonal
	complex
	complexSniper
	complexMissile
)

func setupState(gameState state, s setup) {
	switch s {
	case horizontal:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{5, 1}, infantryConfig))
	case horizontalMissile:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, unguidedGroundMissileConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{5, 1}, guidedAirMissileConfig))
	case vertical:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{1, 5}, infantryConfig))
	case diagonal:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{5, 5}, infantryConfig))
	case complex:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 1}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 1}, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 1}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 1}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 1}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 1}, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 5}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 5}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 5}, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 5}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 5}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 5}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 5}, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 10}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 10}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 10}, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 10}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 10}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 10}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 10}, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 15}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 15}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 15}, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 15}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 15}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 15}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 15}, infantryConfig))
	case complexSniper:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, sniperConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 1}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 1}, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 1}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 1}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 1}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 1}, sniperConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 5}, sniperConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 5}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 5}, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 5}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 5}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 5}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 5}, sniperConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 10}, sniperConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 10}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 10}, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 10}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 10}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 10}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 10}, sniperConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 15}, sniperConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 15}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 15}, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 15}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 15}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 15}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 15}, sniperConfig))
	case complexMissile:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, unguidedGroundMissileConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 1}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 1}, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 1}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 1}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 1}, guidedAirMissileConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 1}, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 5}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 5}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 5}, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 5}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 5}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 5}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 5}, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 10}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 10}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 10}, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 10}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 10}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 10}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 10}, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 15}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 15}, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 15}, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 15}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 15}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 15}, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 15}, infantryConfig))
	}
}
