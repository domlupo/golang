// Hard coded test setups which help manual testing.

package main

type obstacleSetup int

const (
	noObstacle obstacleSetup = iota
	singleObstacle
	doubleObstacle
	tripleObstacle
	multiObstacle
	borderObstacle
	horizontalObstacle
)

func obstacleSetups() []obstacleSetup {
	return []obstacleSetup{
		noObstacle,
		singleObstacle,
		doubleObstacle,
		tripleObstacle,
		multiObstacle,
		borderObstacle,
		horizontalObstacle}
}

func (os obstacleSetup) string() string {
	switch os {
	case noObstacle:
		return "noObstacle"
	case singleObstacle:
		return "singleObstacle"
	case doubleObstacle:
		return "doubleObstacle"
	case tripleObstacle:
		return "tripleObstacle"
	case multiObstacle:
		return "multiObstacle"
	case borderObstacle:
		return "borderObstacle"
	case horizontalObstacle:
		return "horizontalObstacle"
	default:
		panic("obstacleSetup.string: unknown obstacle setup")
	}
}

func setupObstacles(setup obstacleSetup, g grid) (obstacles, obstacleLUT) {
	obstacles := make(obstacles)
	obstacles.setupBorder()
	obstacles.setupInterior(setup)
	oLUT := obstacles.setup(g)
	return obstacles, oLUT
}

func (os obstacles) setupBorder() {
	tiles := []Tile{
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

	for _, t := range tiles {
		os.insert(t, &obstacle{tile: t, kind: highMountain})
	}
}

func (os obstacles) setupInterior(setup obstacleSetup) {
	var tiles []Tile

	switch setup {
	case noObstacle:
	case singleObstacle:
		tiles = []Tile{
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
		tiles = []Tile{
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

			{10, 3},
			{10, 4},
			{10, 5},
			{10, 6},
			{10, 7},
			{10, 8},
			{10, 9},
			{10, 10},
			{10, 11},
			{10, 12},
			{10, 13},
			{10, 14},
			{10, 15},
			{10, 16},
			{10, 17},
			{10, 18},
			{10, 19},
			{10, 20},
			{10, 21},
			{10, 22},
		}
	case tripleObstacle:
		tiles = []Tile{
			{8, 6},
			{8, 7},

			{10, 1},
			{10, 2},
			{10, 3},
			{10, 4},
			{10, 5},
			{10, 6},
			{10, 7},

			{12, 6},
			{12, 7},
		}
	case multiObstacle:
		tiles = []Tile{
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
	case borderObstacle:
		tiles = []Tile{
			{22, 2},

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
		}
	case horizontalObstacle:
		tiles = []Tile{
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
		panic("obstacles.setupInterior: unknown obstacle setup")
	}

	for _, t := range tiles {
		os.insert(t, &obstacle{tile: t, kind: lowMountain})
	}
}

type nodeSetup int

const (
	noNodes nodeSetup = iota
	nodeSurrounded
	northWest
	horizontal
	horizontalNodeMissileCruise
	vertical
	diagonal
	lowNodeCount
	noNeutral
	leftVsRight
	leftVsRightHorizontalOffset
	leftVsRightVerticalOffset
	leftVsRightVerticalOffsetMissileCruise
	leftVsRightVerticalOffsetFighters
	leftVsRightVerticalOffsetAircraftAndOrdnance
	topVsBottom
	topVsBottomHorizontalOffset
	zoomHorizontal
)

func nodeSetups() []nodeSetup {
	return []nodeSetup{
		//noNodes, TODO: debug why this hangs
		nodeSurrounded,
		northWest,
		horizontal,
		horizontalNodeMissileCruise,
		//vertical, TODO debug why this hangs
		diagonal,
		lowNodeCount,
		noNeutral,
		leftVsRight,
		leftVsRightHorizontalOffset,
		leftVsRightVerticalOffset,
		leftVsRightVerticalOffsetMissileCruise,
		leftVsRightVerticalOffsetFighters,
		leftVsRightVerticalOffsetAircraftAndOrdnance,
		topVsBottom,
		topVsBottomHorizontalOffset,
		zoomHorizontal}
}

func (ns nodeSetup) string() string {
	switch ns {
	case noNodes:
		return "noNodes"
	case nodeSurrounded:
		return "nodeSurrounded"
	case northWest:
		return "northWest"
	case horizontal:
		return "horizontal"
	case horizontalNodeMissileCruise:
		return "horizontalNodeMissileCruise"
	case vertical:
		return "vertical"
	case diagonal:
		return "diagonal"
	case lowNodeCount:
		return "lowNodeCount"
	case noNeutral:
		return "noNeutral"
	case leftVsRight:
		return "leftVsRight"
	case leftVsRightHorizontalOffset:
		return "leftVsRightHorizontalOffset"
	case leftVsRightVerticalOffset:
		return "leftVsRightVerticalOffset"
	case leftVsRightVerticalOffsetMissileCruise:
		return "leftVsRightVerticalOffsetMissileCruise"
	case leftVsRightVerticalOffsetFighters:
		return "leftVsRightVerticalOffsetFighters"
	case leftVsRightVerticalOffsetAircraftAndOrdnance:
		return "leftVsRightVerticalOffsetAircraftAndOrdnance"
	case topVsBottom:
		return "topVsBottom"
	case topVsBottomHorizontalOffset:
		return "topVsBottomHorizontalOffset"
	case zoomHorizontal:
		return "zoomHorizontal"
	default:
		panic("nodeSetup.string: unknown node setup")
	}
}

type nodeSetupInput struct {
	tile       Tile
	player     player
	unitConfig unitConfig
}

func setupNodes(
	s nodeSetup,
	g grid,
	initialNodeHealth int,
	h hexagons,
	oLUT obstacleLUT,
	ucProvider unitConfigProvider) (nodes, nodeLUT) {
	nodes := make(nodes)
	nodes.setupPositions(s, initialNodeHealth, ucProvider)
	nodes.setupControl(g, h, oLUT)
	fileHandler := newFileHandler()
	nodeLookup := fileHandler.readNodeLUT(nodes, oLUT, g)
	return nodes, nodeLookup
}

func (nodes nodes) setupPositions(
	s nodeSetup, initialNodeHealth int, ucProvider unitConfigProvider) {
	var nodeSetupInputs []nodeSetupInput

	infantryUC := ucProvider.provide(infantry)
	fighterUC := ucProvider.provide(fighter)
	missileCruiseUC := ucProvider.provide(missileCruise)
	missileHeatSeekingUC := ucProvider.provide(missileHeatSeeking)
	missileBallisticUC := ucProvider.provide(missileBallistic)
	missileAntiBallisticUC := ucProvider.provide(missileAntiBallistic)
	missileRadarActiveUC := ucProvider.provide(missileRadarActive)
	switch s {
	case noNodes:
	case nodeSurrounded:
		player := pOne
		nodeSetupInputs = []nodeSetupInput{
			{Tile{3, 1}, player, infantryUC},
			{Tile{4, 1}, player, infantryUC},
			{Tile{2, 2}, player, infantryUC},
			{Tile{3, 2}, player, infantryUC},
			{Tile{4, 2}, player, infantryUC},
			{Tile{3, 3}, player, infantryUC},
			{Tile{4, 3}, player, infantryUC},
		}
	case northWest:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, infantryUC},
		}
	case horizontal:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{1, 23}, pOne, infantryUC},
			{Tile{3, 23}, pTwo, infantryUC},
		}
	case horizontalNodeMissileCruise:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{1, 23}, pOne, missileCruiseUC},
			{Tile{3, 23}, pTwo, missileHeatSeekingUC},
		}
	case vertical:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{1, 21}, pOne, infantryUC},
			{Tile{1, 23}, pTwo, infantryUC},
		}
	case diagonal:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, infantryUC},
			{Tile{5, 5}, pTwo, infantryUC},
		}
	case lowNodeCount:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, infantryUC},
			{Tile{2, 5}, pOne, infantryUC},
			{Tile{20, 5}, pTwo, infantryUC},
		}
	case noNeutral:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, infantryUC},
			{Tile{6, 1}, pOne, infantryUC},
			{Tile{11, 1}, pOne, infantryUC},
			{Tile{15, 1}, pOne, infantryUC},
			{Tile{21, 1}, pOne, infantryUC},
			{Tile{26, 1}, pOne, infantryUC},

			{Tile{1, 6}, pOne, infantryUC},
			{Tile{6, 6}, pOne, infantryUC},
			{Tile{11, 6}, pOne, infantryUC},
			{Tile{16, 6}, pOne, infantryUC},
			{Tile{21, 6}, pOne, infantryUC},
			{Tile{26, 6}, pOne, infantryUC},

			{Tile{1, 11}, pOne, infantryUC},
			{Tile{6, 11}, pOne, infantryUC},
			{Tile{11, 11}, pOne, infantryUC},

			{Tile{16, 11}, pTwo, infantryUC},
			{Tile{21, 11}, pTwo, infantryUC},
			{Tile{26, 11}, pTwo, infantryUC},

			{Tile{1, 16}, pTwo, infantryUC},
			{Tile{6, 16}, pTwo, infantryUC},
			{Tile{11, 16}, pTwo, infantryUC},
			{Tile{16, 16}, pTwo, infantryUC},
			{Tile{21, 16}, pTwo, infantryUC},
			{Tile{26, 16}, pTwo, infantryUC},

			{Tile{1, 21}, pTwo, infantryUC},
			{Tile{6, 21}, pTwo, infantryUC},
			{Tile{11, 21}, pTwo, infantryUC},
			{Tile{16, 21}, pTwo, infantryUC},
			{Tile{21, 21}, pTwo, infantryUC},
			{Tile{26, 21}, pTwo, infantryUC},
		}
	case leftVsRight:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, infantryUC},
			{Tile{6, 1}, pOne, infantryUC},
			{Tile{11, 1}, pNeutral, infantryUC},
			{Tile{15, 1}, pNeutral, infantryUC},
			{Tile{21, 1}, pTwo, infantryUC},
			{Tile{26, 1}, pTwo, infantryUC},

			{Tile{1, 6}, pOne, infantryUC},
			{Tile{6, 6}, pOne, infantryUC},
			{Tile{11, 6}, pNeutral, infantryUC},
			{Tile{16, 6}, pNeutral, infantryUC},
			{Tile{21, 6}, pTwo, infantryUC},
			{Tile{26, 6}, pTwo, infantryUC},

			{Tile{1, 11}, pOne, infantryUC},
			{Tile{6, 11}, pOne, infantryUC},
			{Tile{11, 11}, pNeutral, infantryUC},
			{Tile{16, 11}, pNeutral, infantryUC},
			{Tile{21, 11}, pTwo, infantryUC},
			{Tile{26, 11}, pTwo, infantryUC},

			{Tile{1, 16}, pOne, infantryUC},
			{Tile{6, 16}, pOne, infantryUC},
			{Tile{11, 16}, pNeutral, infantryUC},
			{Tile{16, 16}, pNeutral, infantryUC},
			{Tile{21, 16}, pTwo, infantryUC},
			{Tile{26, 16}, pTwo, infantryUC},

			{Tile{1, 21}, pOne, infantryUC},
			{Tile{6, 21}, pOne, infantryUC},
			{Tile{11, 21}, pNeutral, infantryUC},
			{Tile{16, 21}, pNeutral, infantryUC},
			{Tile{21, 21}, pTwo, infantryUC},
			{Tile{26, 21}, pTwo, infantryUC},
		}
	case leftVsRightHorizontalOffset:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, infantryUC},
			{Tile{6, 1}, pOne, infantryUC},
			{Tile{11, 1}, pNeutral, infantryUC},
			{Tile{16, 1}, pNeutral, infantryUC},
			{Tile{21, 1}, pTwo, infantryUC},
			{Tile{26, 1}, pTwo, infantryUC},

			{Tile{1, 1}, pOne, infantryUC},
			{Tile{6, 6}, pOne, infantryUC},
			{Tile{11, 6}, pNeutral, infantryUC},
			{Tile{16, 6}, pNeutral, infantryUC},
			{Tile{21, 6}, pTwo, infantryUC},
			{Tile{26, 6}, pTwo, infantryUC},

			{Tile{1, 11}, pOne, infantryUC},
			{Tile{6, 11}, pOne, infantryUC},
			{Tile{11, 11}, pNeutral, infantryUC},
			{Tile{16, 11}, pNeutral, infantryUC},
			{Tile{21, 11}, pTwo, infantryUC},
			{Tile{26, 11}, pTwo, infantryUC},

			{Tile{1, 16}, pOne, infantryUC},
			{Tile{6, 16}, pOne, infantryUC},
			{Tile{11, 16}, pNeutral, infantryUC},
			{Tile{16, 16}, pNeutral, infantryUC},
			{Tile{21, 16}, pTwo, infantryUC},
			{Tile{26, 16}, pTwo, infantryUC},

			{Tile{1, 21}, pOne, infantryUC},
			{Tile{6, 21}, pOne, infantryUC},
			{Tile{11, 21}, pNeutral, infantryUC},
			{Tile{16, 21}, pNeutral, infantryUC},
			{Tile{21, 21}, pTwo, infantryUC},
			{Tile{26, 21}, pTwo, infantryUC},
		}
	case leftVsRightVerticalOffset:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, infantryUC},
			{Tile{6, 1}, pOne, infantryUC},
			{Tile{11, 2}, pNeutral, infantryUC},
			{Tile{16, 4}, pNeutral, infantryUC},
			{Tile{21, 3}, pTwo, infantryUC},
			{Tile{26, 3}, pTwo, infantryUC},

			{Tile{1, 6}, pOne, infantryUC},
			{Tile{6, 6}, pOne, infantryUC},
			{Tile{11, 7}, pNeutral, infantryUC},
			{Tile{16, 9}, pNeutral, infantryUC},
			{Tile{21, 8}, pTwo, infantryUC},
			{Tile{26, 8}, pTwo, infantryUC},

			{Tile{1, 11}, pOne, infantryUC},
			{Tile{6, 11}, pOne, infantryUC},
			{Tile{11, 12}, pNeutral, infantryUC},
			{Tile{16, 14}, pNeutral, infantryUC},
			{Tile{21, 13}, pTwo, infantryUC},
			{Tile{26, 13}, pTwo, infantryUC},

			{Tile{1, 16}, pOne, infantryUC},
			{Tile{6, 16}, pOne, infantryUC},
			{Tile{11, 17}, pNeutral, infantryUC},
			{Tile{16, 19}, pNeutral, infantryUC},
			{Tile{21, 18}, pTwo, infantryUC},
			{Tile{26, 18}, pTwo, infantryUC},

			{Tile{1, 21}, pOne, infantryUC},
			{Tile{6, 21}, pOne, infantryUC},
			{Tile{11, 22}, pNeutral, infantryUC},
			{Tile{16, 22}, pNeutral, infantryUC},
			{Tile{21, 23}, pTwo, infantryUC},
			{Tile{26, 23}, pTwo, infantryUC},
		}
	case leftVsRightVerticalOffsetMissileCruise:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, infantryUC},
			{Tile{6, 1}, pOne, infantryUC},
			{Tile{11, 2}, pNeutral, infantryUC},
			{Tile{16, 4}, pNeutral, infantryUC},
			{Tile{21, 3}, pTwo, infantryUC},
			{Tile{26, 3}, pTwo, infantryUC},

			{Tile{1, 6}, pOne, infantryUC},
			{Tile{6, 6}, pOne, infantryUC},
			{Tile{11, 7}, pNeutral, infantryUC},
			{Tile{16, 9}, pNeutral, infantryUC},
			{Tile{21, 8}, pTwo, infantryUC},
			{Tile{26, 8}, pTwo, infantryUC},

			{Tile{1, 11}, pOne, missileCruiseUC},
			{Tile{6, 11}, pOne, infantryUC},
			{Tile{11, 12}, pNeutral, infantryUC},
			{Tile{16, 14}, pNeutral, infantryUC},
			{Tile{21, 13}, pTwo, infantryUC},
			{Tile{26, 13}, pTwo, missileHeatSeekingUC},

			{Tile{1, 16}, pOne, infantryUC},
			{Tile{6, 16}, pOne, infantryUC},
			{Tile{11, 17}, pNeutral, infantryUC},
			{Tile{16, 19}, pNeutral, infantryUC},
			{Tile{21, 18}, pTwo, infantryUC},
			{Tile{26, 18}, pTwo, infantryUC},

			{Tile{1, 21}, pOne, infantryUC},
			{Tile{6, 21}, pOne, infantryUC},
			{Tile{11, 22}, pNeutral, infantryUC},
			{Tile{16, 22}, pNeutral, infantryUC},
			{Tile{21, 23}, pTwo, infantryUC},
			{Tile{26, 23}, pTwo, infantryUC},
		}
	case leftVsRightVerticalOffsetFighters:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, fighterUC},
			{Tile{6, 1}, pOne, infantryUC},
			{Tile{11, 2}, pNeutral, infantryUC},
			{Tile{16, 4}, pNeutral, infantryUC},
			{Tile{21, 3}, pTwo, infantryUC},
			{Tile{26, 3}, pTwo, infantryUC},

			{Tile{1, 6}, pOne, infantryUC},
			{Tile{6, 6}, pOne, infantryUC},
			{Tile{11, 7}, pNeutral, infantryUC},
			{Tile{16, 9}, pNeutral, infantryUC},
			{Tile{21, 8}, pTwo, infantryUC},
			{Tile{26, 8}, pTwo, infantryUC},

			{Tile{1, 11}, pOne, infantryUC},
			{Tile{6, 11}, pOne, infantryUC},
			{Tile{11, 12}, pNeutral, infantryUC},
			{Tile{16, 14}, pNeutral, infantryUC},
			{Tile{21, 13}, pTwo, infantryUC},
			{Tile{26, 13}, pTwo, infantryUC},

			{Tile{1, 16}, pOne, infantryUC},
			{Tile{6, 16}, pOne, infantryUC},
			{Tile{11, 17}, pNeutral, infantryUC},
			{Tile{16, 19}, pNeutral, infantryUC},
			{Tile{21, 18}, pTwo, infantryUC},
			{Tile{26, 18}, pTwo, infantryUC},

			{Tile{1, 21}, pOne, infantryUC},
			{Tile{6, 21}, pOne, infantryUC},
			{Tile{11, 22}, pNeutral, infantryUC},
			{Tile{16, 22}, pNeutral, infantryUC},
			{Tile{21, 23}, pTwo, infantryUC},
			{Tile{26, 23}, pTwo, infantryUC},
		}

	case leftVsRightVerticalOffsetAircraftAndOrdnance:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, infantryUC},
			{Tile{6, 1}, pOne, infantryUC},
			{Tile{11, 2}, pNeutral, infantryUC},
			{Tile{16, 4}, pNeutral, infantryUC},
			{Tile{21, 3}, pTwo, infantryUC},
			{Tile{26, 3}, pTwo, infantryUC},

			{Tile{1, 6}, pOne, infantryUC},
			{Tile{6, 6}, pOne, infantryUC},
			{Tile{11, 7}, pNeutral, infantryUC},
			{Tile{16, 9}, pNeutral, infantryUC},
			{Tile{21, 8}, pTwo, infantryUC},
			{Tile{26, 8}, pTwo, infantryUC},

			{Tile{1, 11}, pOne, fighterUC},
			{Tile{6, 11}, pOne, infantryUC},
			{Tile{11, 12}, pNeutral, infantryUC},
			{Tile{16, 14}, pNeutral, infantryUC},
			{Tile{21, 13}, pTwo, infantryUC},
			{Tile{26, 13}, pTwo, fighterUC},

			{Tile{1, 16}, pOne, missileCruiseUC},
			{Tile{6, 16}, pOne, infantryUC},
			{Tile{11, 17}, pNeutral, infantryUC},
			{Tile{16, 19}, pNeutral, infantryUC},
			{Tile{21, 18}, pTwo, missileHeatSeekingUC},
			{Tile{26, 18}, pTwo, infantryUC},

			{Tile{1, 21}, pOne, missileBallisticUC},
			{Tile{6, 21}, pOne, missileRadarActiveUC},
			{Tile{11, 22}, pNeutral, infantryUC},
			{Tile{16, 22}, pNeutral, infantryUC},
			{Tile{21, 23}, pTwo, infantryUC},
			{Tile{26, 23}, pTwo, missileAntiBallisticUC},
		}
	case topVsBottom:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, infantryUC},
			{Tile{6, 1}, pOne, infantryUC},
			{Tile{11, 1}, pOne, infantryUC},
			{Tile{16, 1}, pOne, infantryUC},
			{Tile{21, 1}, pOne, infantryUC},
			{Tile{26, 1}, pOne, infantryUC},

			{Tile{1, 5}, pOne, infantryUC},
			{Tile{6, 5}, pOne, infantryUC},
			{Tile{11, 5}, pOne, infantryUC},
			{Tile{16, 5}, pOne, infantryUC},
			{Tile{21, 5}, pOne, infantryUC},
			{Tile{26, 5}, pOne, infantryUC},

			{Tile{1, 9}, pNeutral, infantryUC},
			{Tile{6, 9}, pNeutral, infantryUC},
			{Tile{11, 9}, pNeutral, infantryUC},
			{Tile{16, 9}, pNeutral, infantryUC},
			{Tile{21, 9}, pNeutral, infantryUC},
			{Tile{26, 9}, pNeutral, infantryUC},

			{Tile{1, 13}, pNeutral, infantryUC},
			{Tile{6, 13}, pNeutral, infantryUC},
			{Tile{11, 13}, pNeutral, infantryUC},
			{Tile{16, 13}, pNeutral, infantryUC},
			{Tile{21, 13}, pNeutral, infantryUC},
			{Tile{26, 13}, pNeutral, infantryUC},

			{Tile{1, 17}, pTwo, infantryUC},
			{Tile{6, 17}, pTwo, infantryUC},
			{Tile{11, 17}, pTwo, infantryUC},
			{Tile{16, 17}, pTwo, infantryUC},
			{Tile{21, 17}, pTwo, infantryUC},
			{Tile{26, 17}, pTwo, infantryUC},

			{Tile{1, 21}, pTwo, infantryUC},
			{Tile{6, 21}, pTwo, infantryUC},
			{Tile{11, 21}, pTwo, infantryUC},
			{Tile{16, 21}, pTwo, infantryUC},
			{Tile{21, 21}, pTwo, infantryUC},
			{Tile{26, 21}, pTwo, infantryUC},
		}
	case topVsBottomHorizontalOffset:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, infantryUC},
			{Tile{6, 1}, pOne, infantryUC},
			{Tile{11, 1}, pOne, infantryUC},
			{Tile{16, 1}, pOne, infantryUC},
			{Tile{21, 1}, pOne, infantryUC},
			{Tile{26, 1}, pOne, infantryUC},

			{Tile{1, 6}, pOne, infantryUC},
			{Tile{6, 6}, pOne, infantryUC},
			{Tile{11, 6}, pOne, infantryUC},
			{Tile{16, 6}, pOne, infantryUC},
			{Tile{21, 6}, pOne, infantryUC},
			{Tile{26, 6}, pOne, infantryUC},

			{Tile{1, 11}, pNeutral, infantryUC},
			{Tile{6, 11}, pNeutral, infantryUC},
			{Tile{11, 11}, pNeutral, infantryUC},
			{Tile{16, 11}, pNeutral, infantryUC},
			{Tile{21, 11}, pNeutral, infantryUC},
			{Tile{26, 9}, pNeutral, infantryUC},

			{Tile{1, 16}, pTwo, infantryUC},
			{Tile{6, 16}, pTwo, infantryUC},
			{Tile{11, 16}, pTwo, infantryUC},
			{Tile{16, 16}, pTwo, infantryUC},
			{Tile{21, 16}, pTwo, infantryUC},
			{Tile{26, 16}, pTwo, infantryUC},

			{Tile{1, 21}, pTwo, infantryUC},
			{Tile{6, 21}, pTwo, infantryUC},
			{Tile{11, 21}, pTwo, infantryUC},
			{Tile{16, 21}, pTwo, infantryUC},
			{Tile{21, 21}, pTwo, infantryUC},
			{Tile{26, 21}, pTwo, infantryUC},
		}
	case zoomHorizontal:
		nodeSetupInputs = []nodeSetupInput{
			{Tile{2, 2}, pOne, infantryUC},
			{Tile{6, 2}, pTwo, infantryUC},
		}
	default:
		panic("nodes.setupPositions: unknown setup")
	}

	for _, ni := range nodeSetupInputs {
		node := newNode(ni.player, ni.tile, initialNodeHealth, ni.unitConfig)
		nodes[ni.tile] = &node
	}
}
