// Unit actions for different unit types which occur when a condition is met.

package main

// A unit action which should occur when a condition is met.
type reaction struct {
	when  func(*unit) bool
	do    func(*unit)
	next  uint64
	delay uint64
}

func (u unit) isPayloadNotFull() bool {
	return u.missileCount != u.missileCountMax ||
		u.missileHeatCount != u.missileHeatCountMax ||
		u.flareCount != u.flareCountMax ||
		u.bombCount != u.bombCountMax
}

func (u *unit) refillPayload() {
	u.missileCount = u.missileCountMax
	u.missileHeatCount = u.missileHeatCountMax
	u.flareCount = u.flareCountMax
	u.bombCount = u.bombCountMax
}

// Returns the reactions for the passed unit name.
func reactions(n name) []reaction {
	switch n {
	case bomber:
		return []reaction{
			{
				when: func(bomber *unit) bool {
					// TODO: bresenham is not guaranteed to get to final position. I have not seen
					// aircraft not reach their final destination but missiles have. Need to ensure
					// this is not error prone or make it more conservative.
					return bomber.PixelPosition.equals(bomber.sourcePosition) && bomber.isPayloadNotFull()
				},
				do: func(bomber *unit) {
					bomber.refillPayload()
				},
				delay: bomberConfig.reloadPayloadDelay,
			},
			{
				when: func(bomber *unit) bool { return bomber.nearest != nil },
				do: func(bomber *unit) {
					if bomber.nearest.name == missileHeatSeeking {
						flare := newUnit(
							flareConfig,
							position{layer: flareConfig.produceLayer, PixelPosition: bomber.PixelPosition},
							bomber.player,
							[]position{bomber.nearest.position},
							[]action{},
							GridPosition{},
							bomber.selected)
						bomber.requests = append(bomber.requests, request{spawn: &flare})
						bomber.flareCount--
					}

				},
				delay: bomberConfig.fireOrdnanceDelay,
			}}
	case fighter:
		return []reaction{
			{
				when: func(fighter *unit) bool {
					// TODO: bresenham is not guaranteed to get to final position. I have not seen
					// aircraft not reach their final destination but missiles have. Need to ensure
					// this is not error prone or make it more conservative.
					return fighter.PixelPosition.equals(fighter.sourcePosition) && fighter.isPayloadNotFull()
				},
				do: func(fighter *unit) {
					fighter.refillPayload()
					// TODO: when a missile terminates, it should remove the target from skip targets
					// rather than do this during restock. This should be an upgrade.
					fighter.skipTargets = []id{}
				},
				delay: fighterConfig.reloadPayloadDelay,
			},
			{
				when: func(fighter *unit) bool { return fighter.nearest != nil },
				do: func(fighter *unit) {
					var spawn unit
					switch {
					case fighter.nearest.name == bomb && fighter.missileCount > 0:
						fighter.missileCount--
						spawn = newUnit(
							missileConfig,
							position{layer: missileConfig.produceLayer, PixelPosition: fighter.PixelPosition},
							fighter.player,
							[]position{fighter.nearest.position},
							[]action{},
							GridPosition{},
							fighter.selected)
					case fighter.nearest.name == bomber && fighter.missileHeatCount > 0:
						fighter.missileHeatCount--
						spawn = newUnit(
							missileHeatSeekingConfig,
							position{layer: missileConfig.produceLayer, PixelPosition: fighter.PixelPosition},
							fighter.player,
							[]position{fighter.nearest.position},
							[]action{},
							GridPosition{},
							fighter.selected)
					default:
						return
					}
					spawn.sourceUnit = fighter.id
					fighter.requests = append(fighter.requests, request{spawn: &spawn})
				},
				delay: fighterConfig.fireOrdnanceDelay,
			}}
	default:
		return []reaction{}
	}
}
