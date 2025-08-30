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
	switch s {
	case horizontal:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}))
		gameState.addNode(newNode(redPlayer, gridPosition{5, 1}))
	case vertical:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}))
		gameState.addNode(newNode(redPlayer, gridPosition{1, 5}))
	case diagonal:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}))
		gameState.addNode(newNode(redPlayer, gridPosition{5, 5}))
	case complex:
		gameState.addNode(newNode(bluePlayer, gridPosition{1, 1}))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 1}))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 1}))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 1}))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 1}))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 1}))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 1}))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 5}))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 5}))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 5}))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 5}))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 5}))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 5}))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 5}))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 10}))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 10}))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 10}))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 10}))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 10}))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 10}))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 10}))

		gameState.addNode(newNode(bluePlayer, gridPosition{1, 15}))
		gameState.addNode(newNode(bluePlayer, gridPosition{5, 15}))
		gameState.addNode(newNode(bluePlayer, gridPosition{10, 15}))
		gameState.addNode(newNode(neutralPlayer, gridPosition{15, 15}))
		gameState.addNode(newNode(redPlayer, gridPosition{20, 15}))
		gameState.addNode(newNode(redPlayer, gridPosition{25, 15}))
		gameState.addNode(newNode(redPlayer, gridPosition{30, 15}))
	}
}
