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

func setupState(gameState state, s setup, initialNodeHealth int, nodeRadius int) {
	switch s {
	case horizontal:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{5, 1}, initialNodeHealth, nodeRadius, infantryConfig))
	case horizontalMissile:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, initialNodeHealth, nodeRadius, unguidedGroundMissileConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{5, 1}, initialNodeHealth, nodeRadius, guidedAirMissileConfig))
	case vertical:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{1, 5}, initialNodeHealth, nodeRadius, infantryConfig))
	case diagonal:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{5, 5}, initialNodeHealth, nodeRadius, infantryConfig))
	case complex:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 1}, initialNodeHealth, nodeRadius, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 5}, initialNodeHealth, nodeRadius, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 10}, initialNodeHealth, nodeRadius, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 15}, initialNodeHealth, nodeRadius, infantryConfig))
	case complexSniper:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, initialNodeHealth, nodeRadius, sniperConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 1}, initialNodeHealth, nodeRadius, sniperConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 5}, initialNodeHealth, nodeRadius, sniperConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 5}, initialNodeHealth, nodeRadius, sniperConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 10}, initialNodeHealth, nodeRadius, sniperConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 10}, initialNodeHealth, nodeRadius, sniperConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 15}, initialNodeHealth, nodeRadius, sniperConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 15}, initialNodeHealth, nodeRadius, sniperConfig))
	case complexMissile:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, initialNodeHealth, nodeRadius, unguidedGroundMissileConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{15, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 1}, initialNodeHealth, nodeRadius, guidedAirMissileConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 1}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 1}, initialNodeHealth, nodeRadius, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 5}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 5}, initialNodeHealth, nodeRadius, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 10}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 10}, initialNodeHealth, nodeRadius, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 15}, initialNodeHealth, nodeRadius, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 15}, initialNodeHealth, nodeRadius, infantryConfig))
	default:
		panic("setup does not exist")
	}
}
