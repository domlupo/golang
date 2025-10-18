// hard coded test setups which help manual testing
package main

type obstacleSetup int

const (
	noObstacle obstacleSetup = iota
	singleObstacle
	doubleObstacle
	multiObstacle
	horizontalObstacle
)

func setupObstacles(setup obstacleSetup, hc *hexConfig) (obstacles, pixelPathing) {
	obstacles := make(obstacles)
	obstacles.setupBorder()
	obstacles.setupInterior(setup, hc)
	pixelPathing := obstacles.setup(hc)
	return obstacles, pixelPathing
}

func (os obstacles) setupBorder() {
	gridPositions := []gridPosition{
		// north west corner
		{0, 0},
		{1, 1},
		// north
		{1, 0},
		{2, 0},
		{3, 0},
		{4, 0},
		{5, 0},
		{6, 0},
		{7, 0},
		{8, 0},
		{9, 0},
		{10, 0},
		{11, 0},
		{12, 0},
		{13, 0},
		{14, 0},
		{15, 0},
		{16, 0},
		{17, 0},
		{18, 0},
		{19, 0},
		{20, 0},
		{21, 0},
		{22, 0},
		{23, 0},
		{24, 0},
		{25, 0},
		{26, 0},
		// north east corner
		{27, 0},
		{27, 1},
		// east
		{28, 1},
		{27, 2},
		{28, 3},
		{27, 4},
		{28, 5},
		{27, 6},
		{28, 7},
		{27, 8},
		{28, 9},
		{27, 10},
		{28, 11},
		{27, 12},
		{28, 13},
		{27, 14},
		{28, 15},
		{27, 16},
		{28, 17},
		{27, 18},
		{28, 19},
		{27, 20},
		{28, 21},
		{27, 22},
		{28, 23},
		// south east corner
		{27, 24},
		{27, 23},
		// south west corner
		{0, 24},
		{1, 23},
		// south
		{1, 24},
		{2, 24},
		{3, 24},
		{4, 24},
		{5, 24},
		{6, 24},
		{7, 24},
		{8, 24},
		{9, 24},
		{10, 24},
		{11, 24},
		{12, 24},
		{13, 24},
		{14, 24},
		{15, 24},
		{16, 24},
		{17, 24},
		{18, 24},
		{19, 24},
		{20, 24},
		{21, 24},
		{22, 24},
		{23, 24},
		{24, 24},
		{25, 24},
		{26, 24},
		// west
		{0, 1},
		{0, 2},
		{0, 3},
		{0, 4},
		{0, 5},
		{0, 6},
		{0, 7},
		{0, 8},
		{0, 9},
		{0, 10},
		{0, 11},
		{0, 12},
		{0, 13},
		{0, 14},
		{0, 15},
		{0, 16},
		{0, 17},
		{0, 18},
		{0, 19},
		{0, 20},
		{0, 21},
		{0, 22},
		{0, 23},
	}

	for _, pos := range gridPositions {
		os.insert(pos, &obstacle{gridPos: pos, kind: highMountain})
	}
}

func (os obstacles) setupInterior(setup obstacleSetup, hc *hexConfig) {
	var gridPositions []gridPosition

	switch setup {
	case noObstacle:
	case singleObstacle:
		gridPositions = []gridPosition{
			{3, 3},
			{3, 4},
			{4, 4},
			{5, 3},
			{5, 4},
			{6, 4},
			{7, 4},
			{8, 3},
			{5, 5},
			{5, 6},
			{5, 7},
			{5, 8},
			{5, 9},
			{5, 10},
			{5, 11},
			{5, 12},
			{5, 13},
			{5, 14},
			{5, 15},
			{5, 16},
			{5, 17},
			{5, 18},
			{5, 19},
			{5, 20},
			{5, 21},
			{5, 22},
		}
	case doubleObstacle:
		gridPositions = []gridPosition{
			{3, 3},
			{3, 4},
			{4, 4},
			{5, 3},
			{5, 4},
			{6, 4},
			{7, 4},
			{8, 3},
			{5, 5},
			{5, 6},
			{5, 7},
			{5, 8},
			{5, 9},
			{5, 10},
			{5, 11},
			{5, 12},
			{5, 13},
			{5, 14},
			{5, 15},
			{5, 16},
			{5, 17},
			{5, 18},
			{5, 19},
			{5, 20},
			{5, 21},
			{5, 22},

			{19, 3},
			{19, 4},
			{19, 5},
			{19, 6},
			{19, 7},
			{19, 8},
			{19, 9},
			{19, 10},
			{19, 11},
			{19, 12},
			{19, 13},
			{19, 14},
			{19, 15},
			{19, 16},
			{19, 17},
			{19, 18},
			{19, 19},
			{19, 20},
			{19, 21},
			{19, 22},
		}
	case multiObstacle:
		gridPositions = []gridPosition{
			{2, 4},
			{3, 3},
			{3, 4},
			{4, 4},
			{5, 4},
			{6, 3},
			{3, 5},
			{3, 6},
			{3, 7},
			{3, 8},
			{3, 9},
			{3, 10},
			{3, 11},
			{3, 12},
			{3, 13},
			{3, 14},
			{3, 15},
			{3, 16},
			{3, 17},
			{3, 18},
			{3, 19},

			{7, 2},

			{11, 3},
			{12, 3},

			{9, 3},
			{9, 4},
			{9, 5},
			{9, 6},
			{9, 7},
			{9, 8},
			{9, 9},
			{9, 10},
			{9, 11},
			{9, 12},
			{9, 13},
			{9, 14},
			{9, 15},
			{9, 16},
			{9, 17},
			{9, 18},
			{9, 19},
			{9, 20},
			{9, 21},
			{9, 22},
			{9, 23},

			{13, 4},
			{13, 5},
			{13, 6},
			{13, 7},
			{13, 8},
			{13, 9},
			{13, 10},
			{13, 11},
			{13, 12},
			{13, 13},
			{13, 14},
			{13, 15},
			{13, 16},
			{13, 17},
			{13, 18},
			{13, 19},
			{13, 20},
			{13, 21},
			{13, 22},
		}
	case horizontalObstacle:
		gridPositions = []gridPosition{
			{1, 3},
			{2, 3},
			{3, 3},
			{4, 3},
			{5, 3},
			{6, 3},
			{7, 3},
			{8, 3},
			{9, 3},
			{10, 3},
			{10, 4},
		}
	default:
		panic("obstacle setup does not exist")
	}

	for _, pos := range gridPositions {
		os.insert(pos, &obstacle{gridPos: pos, kind: lowMountain})
	}
}

type nodeSetup int

const (
	noNodes nodeSetup = iota
	nodeSurrounded
	northWest
	horizontal
	horizontalMissile
	vertical
	diagonal
	lowNodeCount
	noNeutral
	leftVsRight
	leftVsRightHorizontalOffset
	leftVsRightVerticalOffset
	topVsBottom
	topVsBottomHorizontalOffset
	zoomHorizontal
	zoomVertical
)

type nodeInfo struct {
	gridPos    gridPosition
	player     player
	unitConfig unitConfig
}

func setupNodes(s nodeSetup, hc *hexConfig, initialNodeHealth int, o obstacles, h hexagons) nodes {
	nodes := make(nodes)
	nodes.setupPositions(s, hc, initialNodeHealth)
	nodes.setupControl(hc, h, o)
	return nodes
}

func (nodes nodes) setupPositions(s nodeSetup, hc *hexConfig, initialNodeHealth int) {
	var nodeInfos []nodeInfo

	switch s {
	case noNodes:
	case nodeSurrounded:
		player := bluePlayer
		unitConfig := infantryConfig
		nodeInfos = []nodeInfo{
			{gridPosition{3, 1}, player, unitConfig},
			{gridPosition{4, 1}, player, unitConfig},
			{gridPosition{2, 2}, player, unitConfig},
			{gridPosition{3, 2}, player, unitConfig},
			{gridPosition{4, 2}, player, unitConfig},
			{gridPosition{3, 3}, player, unitConfig},
			{gridPosition{4, 3}, player, unitConfig},
		}
	case northWest:
		nodeInfos = []nodeInfo{
			{gridPosition{2, 2}, bluePlayer, infantryConfig},
		}
	case horizontal:
		nodeInfos = []nodeInfo{
			{gridPosition{1, 23}, bluePlayer, infantryConfig},
			{gridPosition{3, 23}, redPlayer, infantryConfig},
		}
	case horizontalMissile:
		nodeInfos = []nodeInfo{
			{gridPosition{1, 23}, bluePlayer, unguidedGroundMissileConfig},
			{gridPosition{3, 23}, redPlayer, guidedAirMissileConfig},
		}
	case vertical:
		nodeInfos = []nodeInfo{
			{gridPosition{1, 21}, bluePlayer, infantryConfig},
			{gridPosition{1, 23}, redPlayer, infantryConfig},
		}
	case diagonal:
		nodeInfos = []nodeInfo{
			{gridPosition{2, 2}, bluePlayer, infantryConfig},
			{gridPosition{5, 5}, redPlayer, infantryConfig},
		}
	case lowNodeCount:
		nodeInfos = []nodeInfo{
			{gridPosition{2, 2}, bluePlayer, infantryConfig},
			{gridPosition{2, 5}, bluePlayer, infantryConfig},
			{gridPosition{20, 5}, redPlayer, infantryConfig},
		}
	case noNeutral:
		unitConfig := infantryConfig
		nodeInfos = []nodeInfo{
			{gridPosition{2, 2}, bluePlayer, unitConfig},
			{gridPosition{6, 1}, bluePlayer, unitConfig},
			{gridPosition{11, 1}, bluePlayer, unitConfig},
			{gridPosition{15, 1}, bluePlayer, unitConfig},
			{gridPosition{21, 1}, bluePlayer, unitConfig},
			{gridPosition{26, 1}, bluePlayer, unitConfig},

			{gridPosition{1, 6}, bluePlayer, unitConfig},
			{gridPosition{6, 6}, bluePlayer, unitConfig},
			{gridPosition{11, 6}, bluePlayer, unitConfig},
			{gridPosition{16, 6}, bluePlayer, unitConfig},
			{gridPosition{21, 6}, bluePlayer, unitConfig},
			{gridPosition{26, 6}, bluePlayer, unitConfig},

			{gridPosition{1, 11}, bluePlayer, unitConfig},
			{gridPosition{6, 11}, bluePlayer, unitConfig},
			{gridPosition{11, 11}, bluePlayer, unitConfig},

			{gridPosition{16, 11}, redPlayer, unitConfig},
			{gridPosition{21, 11}, redPlayer, unitConfig},
			{gridPosition{26, 11}, redPlayer, unitConfig},

			{gridPosition{1, 16}, redPlayer, unitConfig},
			{gridPosition{6, 16}, redPlayer, unitConfig},
			{gridPosition{11, 16}, redPlayer, unitConfig},
			{gridPosition{16, 16}, redPlayer, unitConfig},
			{gridPosition{21, 16}, redPlayer, unitConfig},
			{gridPosition{26, 16}, redPlayer, unitConfig},

			{gridPosition{1, 21}, redPlayer, unitConfig},
			{gridPosition{6, 21}, redPlayer, unitConfig},
			{gridPosition{11, 21}, redPlayer, unitConfig},
			{gridPosition{16, 21}, redPlayer, unitConfig},
			{gridPosition{21, 21}, redPlayer, unitConfig},
			{gridPosition{26, 21}, redPlayer, unitConfig},
		}
	case leftVsRight:
		unitConfig := infantryConfig
		nodeInfos = []nodeInfo{
			{gridPosition{2, 2}, bluePlayer, unitConfig},
			{gridPosition{6, 1}, bluePlayer, unitConfig},
			{gridPosition{11, 1}, neutralPlayer, unitConfig},
			{gridPosition{15, 1}, neutralPlayer, unitConfig},
			{gridPosition{21, 1}, redPlayer, unitConfig},
			{gridPosition{26, 1}, redPlayer, unitConfig},

			{gridPosition{1, 6}, bluePlayer, unitConfig},
			{gridPosition{6, 6}, bluePlayer, unitConfig},
			{gridPosition{11, 6}, neutralPlayer, unitConfig},
			{gridPosition{16, 6}, neutralPlayer, unitConfig},
			{gridPosition{21, 6}, redPlayer, unitConfig},
			{gridPosition{26, 6}, redPlayer, unitConfig},

			{gridPosition{1, 11}, bluePlayer, unitConfig},
			{gridPosition{6, 11}, bluePlayer, unitConfig},
			{gridPosition{11, 11}, neutralPlayer, unitConfig},
			{gridPosition{16, 11}, neutralPlayer, unitConfig},
			{gridPosition{21, 11}, redPlayer, unitConfig},
			{gridPosition{26, 11}, redPlayer, unitConfig},

			{gridPosition{1, 16}, bluePlayer, unitConfig},
			{gridPosition{6, 16}, bluePlayer, unitConfig},
			{gridPosition{11, 16}, neutralPlayer, unitConfig},
			{gridPosition{16, 16}, neutralPlayer, unitConfig},
			{gridPosition{21, 16}, redPlayer, unitConfig},
			{gridPosition{26, 16}, redPlayer, unitConfig},

			{gridPosition{1, 21}, bluePlayer, unitConfig},
			{gridPosition{6, 21}, bluePlayer, unitConfig},
			{gridPosition{11, 21}, neutralPlayer, unitConfig},
			{gridPosition{16, 21}, neutralPlayer, unitConfig},
			{gridPosition{21, 21}, redPlayer, unitConfig},
			{gridPosition{26, 21}, redPlayer, unitConfig},
		}
	case leftVsRightHorizontalOffset:
		unitConfig := infantryConfig
		nodeInfos = []nodeInfo{
			{gridPosition{2, 2}, bluePlayer, unitConfig},
			{gridPosition{6, 1}, bluePlayer, unitConfig},
			{gridPosition{11, 1}, neutralPlayer, unitConfig},
			{gridPosition{16, 1}, neutralPlayer, unitConfig},
			{gridPosition{21, 1}, redPlayer, unitConfig},
			{gridPosition{26, 1}, redPlayer, unitConfig},

			{gridPosition{1, 1}, bluePlayer, unitConfig},
			{gridPosition{6, 6}, bluePlayer, unitConfig},
			{gridPosition{11, 6}, neutralPlayer, unitConfig},
			{gridPosition{16, 6}, neutralPlayer, unitConfig},
			{gridPosition{21, 6}, redPlayer, unitConfig},
			{gridPosition{26, 6}, redPlayer, unitConfig},

			{gridPosition{1, 11}, bluePlayer, unitConfig},
			{gridPosition{6, 11}, bluePlayer, unitConfig},
			{gridPosition{11, 11}, neutralPlayer, unitConfig},
			{gridPosition{16, 11}, neutralPlayer, unitConfig},
			{gridPosition{21, 11}, redPlayer, unitConfig},
			{gridPosition{26, 11}, redPlayer, unitConfig},

			{gridPosition{1, 16}, bluePlayer, unitConfig},
			{gridPosition{6, 16}, bluePlayer, unitConfig},
			{gridPosition{11, 16}, neutralPlayer, unitConfig},
			{gridPosition{16, 16}, neutralPlayer, unitConfig},
			{gridPosition{21, 16}, redPlayer, unitConfig},
			{gridPosition{26, 16}, redPlayer, unitConfig},

			{gridPosition{1, 21}, bluePlayer, unitConfig},
			{gridPosition{6, 21}, bluePlayer, unitConfig},
			{gridPosition{11, 21}, neutralPlayer, unitConfig},
			{gridPosition{16, 21}, neutralPlayer, unitConfig},
			{gridPosition{21, 21}, redPlayer, unitConfig},
			{gridPosition{26, 21}, redPlayer, unitConfig},
		}
	case leftVsRightVerticalOffset:
		unitConfig := infantryConfig
		nodeInfos = []nodeInfo{
			{gridPosition{2, 2}, bluePlayer, unitConfig},
			{gridPosition{6, 1}, bluePlayer, unitConfig},
			{gridPosition{11, 2}, neutralPlayer, unitConfig},
			{gridPosition{16, 4}, neutralPlayer, unitConfig},
			{gridPosition{21, 3}, redPlayer, unitConfig},
			{gridPosition{26, 3}, redPlayer, unitConfig},

			{gridPosition{1, 6}, bluePlayer, unitConfig},
			{gridPosition{6, 6}, bluePlayer, unitConfig},
			{gridPosition{11, 7}, neutralPlayer, unitConfig},
			{gridPosition{16, 9}, neutralPlayer, unitConfig},
			{gridPosition{21, 8}, redPlayer, unitConfig},
			{gridPosition{26, 8}, redPlayer, unitConfig},

			{gridPosition{1, 11}, bluePlayer, unitConfig},
			{gridPosition{6, 11}, bluePlayer, unitConfig},
			{gridPosition{11, 12}, neutralPlayer, unitConfig},
			{gridPosition{16, 14}, neutralPlayer, unitConfig},
			{gridPosition{21, 13}, redPlayer, unitConfig},
			{gridPosition{26, 13}, redPlayer, unitConfig},

			{gridPosition{1, 16}, bluePlayer, unitConfig},
			{gridPosition{6, 16}, bluePlayer, unitConfig},
			{gridPosition{11, 17}, neutralPlayer, unitConfig},
			{gridPosition{16, 19}, neutralPlayer, unitConfig},
			{gridPosition{21, 18}, redPlayer, unitConfig},
			{gridPosition{26, 18}, redPlayer, unitConfig},

			{gridPosition{1, 21}, bluePlayer, unitConfig},
			{gridPosition{6, 21}, bluePlayer, unitConfig},
			{gridPosition{11, 22}, neutralPlayer, unitConfig},
			{gridPosition{16, 22}, neutralPlayer, unitConfig},
			{gridPosition{21, 23}, redPlayer, unitConfig},
			{gridPosition{26, 23}, redPlayer, unitConfig},
		}
	case topVsBottom:
		unitConfig := infantryConfig
		nodeInfos = []nodeInfo{
			{gridPosition{2, 2}, bluePlayer, unitConfig},
			{gridPosition{6, 1}, bluePlayer, unitConfig},
			{gridPosition{11, 1}, bluePlayer, unitConfig},
			{gridPosition{16, 1}, bluePlayer, unitConfig},
			{gridPosition{21, 1}, bluePlayer, unitConfig},
			{gridPosition{26, 1}, bluePlayer, unitConfig},

			{gridPosition{1, 5}, bluePlayer, unitConfig},
			{gridPosition{6, 5}, bluePlayer, unitConfig},
			{gridPosition{11, 5}, bluePlayer, unitConfig},
			{gridPosition{16, 5}, bluePlayer, unitConfig},
			{gridPosition{21, 5}, bluePlayer, unitConfig},
			{gridPosition{26, 5}, bluePlayer, unitConfig},

			{gridPosition{1, 9}, neutralPlayer, unitConfig},
			{gridPosition{6, 9}, neutralPlayer, unitConfig},
			{gridPosition{11, 9}, neutralPlayer, unitConfig},
			{gridPosition{16, 9}, neutralPlayer, unitConfig},
			{gridPosition{21, 9}, neutralPlayer, unitConfig},
			{gridPosition{26, 9}, neutralPlayer, unitConfig},

			{gridPosition{1, 13}, neutralPlayer, unitConfig},
			{gridPosition{6, 13}, neutralPlayer, unitConfig},
			{gridPosition{11, 13}, neutralPlayer, unitConfig},
			{gridPosition{16, 13}, neutralPlayer, unitConfig},
			{gridPosition{21, 13}, neutralPlayer, unitConfig},
			{gridPosition{26, 13}, neutralPlayer, unitConfig},

			{gridPosition{1, 17}, redPlayer, unitConfig},
			{gridPosition{6, 17}, redPlayer, unitConfig},
			{gridPosition{11, 17}, redPlayer, unitConfig},
			{gridPosition{16, 17}, redPlayer, unitConfig},
			{gridPosition{21, 17}, redPlayer, unitConfig},
			{gridPosition{26, 17}, redPlayer, unitConfig},

			{gridPosition{1, 21}, redPlayer, unitConfig},
			{gridPosition{6, 21}, redPlayer, unitConfig},
			{gridPosition{11, 21}, redPlayer, unitConfig},
			{gridPosition{16, 21}, redPlayer, unitConfig},
			{gridPosition{21, 21}, redPlayer, unitConfig},
			{gridPosition{26, 21}, redPlayer, unitConfig},
		}
	case topVsBottomHorizontalOffset:
		unitConfig := infantryConfig
		nodeInfos = []nodeInfo{
			{gridPosition{2, 2}, bluePlayer, unitConfig},
			{gridPosition{6, 1}, bluePlayer, unitConfig},
			{gridPosition{11, 1}, bluePlayer, unitConfig},
			{gridPosition{16, 1}, bluePlayer, unitConfig},
			{gridPosition{21, 1}, bluePlayer, unitConfig},
			{gridPosition{26, 1}, bluePlayer, unitConfig},

			{gridPosition{1, 6}, bluePlayer, unitConfig},
			{gridPosition{6, 6}, bluePlayer, unitConfig},
			{gridPosition{11, 6}, bluePlayer, unitConfig},
			{gridPosition{16, 6}, bluePlayer, unitConfig},
			{gridPosition{21, 6}, bluePlayer, unitConfig},
			{gridPosition{26, 6}, bluePlayer, unitConfig},

			{gridPosition{1, 11}, neutralPlayer, unitConfig},
			{gridPosition{6, 11}, neutralPlayer, unitConfig},
			{gridPosition{11, 11}, neutralPlayer, unitConfig},
			{gridPosition{16, 11}, neutralPlayer, unitConfig},
			{gridPosition{21, 11}, neutralPlayer, unitConfig},
			{gridPosition{26, 9}, neutralPlayer, unitConfig},

			{gridPosition{1, 16}, redPlayer, unitConfig},
			{gridPosition{6, 16}, redPlayer, unitConfig},
			{gridPosition{11, 16}, redPlayer, unitConfig},
			{gridPosition{16, 16}, redPlayer, unitConfig},
			{gridPosition{21, 16}, redPlayer, unitConfig},
			{gridPosition{26, 16}, redPlayer, unitConfig},

			{gridPosition{1, 21}, redPlayer, unitConfig},
			{gridPosition{6, 21}, redPlayer, unitConfig},
			{gridPosition{11, 21}, redPlayer, unitConfig},
			{gridPosition{16, 21}, redPlayer, unitConfig},
			{gridPosition{21, 21}, redPlayer, unitConfig},
			{gridPosition{26, 21}, redPlayer, unitConfig},
		}
	case zoomHorizontal:
		nodeInfos = []nodeInfo{
			{gridPosition{2, 2}, bluePlayer, infantryConfig},
			{gridPosition{6, 2}, redPlayer, infantryConfig},
		}
	default:
		panic("node setup does not exist")
	}

	for _, ni := range nodeInfos {
		node := newNode(ni.player, ni.gridPos, initialNodeHealth, hc, ni.unitConfig)
		nodes[ni.gridPos] = &node
	}
}
