package main

// file contains different setups

type setup int

const (
	horizontal = iota
	vertical
	diagonal
	complex
)

func setupState(gameState state, s setup) {
	var bluePlayerNode = newNode(bluePlayer)
	var redPlayerNode = newNode(redPlayer)
	var neutralPlayerNode = newNode(neutralPlayer)

	switch s {
	case horizontal:
		gameState.addNode(gridPosition{1, 1}, bluePlayerNode)
		gameState.addNode(gridPosition{5, 1}, redPlayerNode)
	case vertical:
		gameState.addNode(gridPosition{1, 1}, bluePlayerNode)
		gameState.addNode(gridPosition{1, 5}, redPlayerNode)
	case diagonal:
		gameState.addNode(gridPosition{1, 1}, bluePlayerNode)
		gameState.addNode(gridPosition{5, 5}, redPlayerNode)
	case complex:
		gameState.addNode(gridPosition{1, 1}, bluePlayerNode)
		gameState.addNode(gridPosition{5, 1}, bluePlayerNode)
		gameState.addNode(gridPosition{10, 1}, bluePlayerNode)
		gameState.addNode(gridPosition{15, 1}, neutralPlayerNode)
		gameState.addNode(gridPosition{20, 1}, redPlayerNode)
		gameState.addNode(gridPosition{25, 1}, redPlayerNode)
		gameState.addNode(gridPosition{30, 1}, redPlayerNode)

		gameState.addNode(gridPosition{1, 5}, bluePlayerNode)
		gameState.addNode(gridPosition{5, 5}, bluePlayerNode)
		gameState.addNode(gridPosition{10, 5}, bluePlayerNode)
		gameState.addNode(gridPosition{15, 5}, neutralPlayerNode)
		gameState.addNode(gridPosition{20, 5}, redPlayerNode)
		gameState.addNode(gridPosition{25, 5}, redPlayerNode)
		gameState.addNode(gridPosition{30, 5}, redPlayerNode)

		gameState.addNode(gridPosition{1, 10}, bluePlayerNode)
		gameState.addNode(gridPosition{5, 10}, bluePlayerNode)
		gameState.addNode(gridPosition{10, 10}, bluePlayerNode)
		gameState.addNode(gridPosition{15, 10}, neutralPlayerNode)
		gameState.addNode(gridPosition{20, 10}, redPlayerNode)
		gameState.addNode(gridPosition{25, 10}, redPlayerNode)
		gameState.addNode(gridPosition{30, 10}, redPlayerNode)

		gameState.addNode(gridPosition{1, 15}, bluePlayerNode)
		gameState.addNode(gridPosition{5, 15}, bluePlayerNode)
		gameState.addNode(gridPosition{10, 15}, bluePlayerNode)
		gameState.addNode(gridPosition{15, 15}, neutralPlayerNode)
		gameState.addNode(gridPosition{20, 15}, redPlayerNode)
		gameState.addNode(gridPosition{25, 15}, redPlayerNode)
		gameState.addNode(gridPosition{30, 15}, redPlayerNode)
	}
}
