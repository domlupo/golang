// hard coded test setups which help manual testing
package main

type setup int

const (
	horizontal = iota
	horizontalMissile
	vertical
	diagonal
	complex
	horizontalOffsetFull
	verticalOffsetFull
)

func setupState(gameState state, s setup, hexConfig hexConfig, initialNodeHealth int) {
	switch s {
	case horizontal:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 23}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{3, 23}, initialNodeHealth, hexConfig, infantryConfig))
	case horizontalMissile:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 23}, initialNodeHealth, hexConfig, unguidedGroundMissileConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{3, 23}, initialNodeHealth, hexConfig, guidedAirMissileConfig))
	case vertical:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{1, 23}, initialNodeHealth, hexConfig, infantryConfig))
	case diagonal:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{5, 5}, initialNodeHealth, hexConfig, infantryConfig))
	case complex:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 1}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 5}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 5}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 5}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 5}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 5}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 5}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 9}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 9}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 9}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 9}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 9}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 9}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 13}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 13}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 13}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 13}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 13}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 13}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 17}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 17}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 17}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 17}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 17}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 17}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 21}, initialNodeHealth, hexConfig, infantryConfig))
	case horizontalOffsetFull:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 1}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 6}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 6}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 6}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 6}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 6}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 6}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 11}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 11}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 11}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 11}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 11}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 11}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 16}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 16}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 16}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 16}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 16}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 16}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 21}, initialNodeHealth, hexConfig, infantryConfig))
	case verticalOffsetFull:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 1}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 2}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 4}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 3}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 3}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 6}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 6}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 7}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 9}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 8}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 8}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 11}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 11}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 12}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 14}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 13}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 13}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 16}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 16}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 17}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 19}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 18}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 18}, initialNodeHealth, hexConfig, infantryConfig))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(bluePlayer, gridPosition{6, 21}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{11, 22}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(neutralPlayer, gridPosition{16, 22}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{21, 23}, initialNodeHealth, hexConfig, infantryConfig))
		gameState.addNode(newNode(redPlayer, gridPosition{26, 23}, initialNodeHealth, hexConfig, infantryConfig))
	default:
		panic("setup does not exist")
	}
}
