// hard coded test setups which help manual testing
package main

type nodeSetup int

const (
	noNodes nodeSetup = iota
	nodeSurrounded
	zeroZero
	oneOne
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

func setupNodes(s nodeSetup, hc *hexConfig, initialNodeHealth int) map[gridPosition]*node {
	nodes := make(map[gridPosition]*node)
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
	case zeroZero:
		nodeInfos = []nodeInfo{
			{gridPosition{0, 0}, bluePlayer, infantryConfig},
		}
	case oneOne:
		nodeInfos = []nodeInfo{
			{gridPosition{1, 1}, bluePlayer, infantryConfig},
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
			{gridPosition{1, 1}, bluePlayer, infantryConfig},
			{gridPosition{5, 5}, redPlayer, infantryConfig},
		}
	case lowNodeCount:
		nodeInfos = []nodeInfo{
			{gridPosition{1, 1}, bluePlayer, infantryConfig},
			{gridPosition{2, 5}, bluePlayer, infantryConfig},
			{gridPosition{20, 5}, redPlayer, infantryConfig},
		}
	case noNeutral:
		unitConfig := infantryConfig
		nodeInfos = []nodeInfo{
			{gridPosition{1, 1}, bluePlayer, unitConfig},
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
			{gridPosition{1, 1}, bluePlayer, unitConfig},
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
			{gridPosition{1, 1}, bluePlayer, unitConfig},
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
			{gridPosition{1, 1}, bluePlayer, unitConfig},
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
			{gridPosition{1, 1}, bluePlayer, unitConfig},
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
			{gridPosition{1, 1}, bluePlayer, unitConfig},
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
		panic("setup does not exist")
	}

	for _, ni := range nodeInfos {
		node := newNode(ni.player, ni.gridPos, initialNodeHealth, hc, ni.unitConfig)
		nodes[ni.gridPos] = &node
	}

	return nodes
}
