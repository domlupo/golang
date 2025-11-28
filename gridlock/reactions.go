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
	return u.missileHeatSeekingCount != u.missileHeatSeekingCountMax ||
		u.flareCount != u.flareCountMax ||
		u.bombGravityCount != u.bombGravityCountMax ||
		u.missileRadarActiveCount != u.missileRadarActiveCountMax ||
		u.chaffCount != u.chaffCountMax
}

func (u *unit) refillPayload() {
	u.missileHeatSeekingCount = u.missileHeatSeekingCountMax
	u.flareCount = u.flareCountMax
	u.bombGravityCount = u.bombGravityCountMax
	u.missileRadarActiveCount = u.missileRadarActiveCountMax
	u.chaffCount = u.chaffCountMax
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
					return bomber.Pixel.equals(bomber.sourcePosition) && bomber.isPayloadNotFull()
				},
				do: func(bomber *unit) {
					bomber.refillPayload()
				},
				delay: bomberConfig.reloadPayloadDelay,
			},
			reactFlare(bomberConfig.fireOrdnanceDelay),
			reactChaff(bomberConfig.fireOrdnanceDelay),
		}
	case fighter:
		return []reaction{
			{
				when: func(fighter *unit) bool {
					// TODO: bresenham is not guaranteed to get to final position. I have not seen
					// aircraft not reach their final destination but missiles have. Need to ensure
					// this is not error prone or make it more conservative.
					return fighter.Pixel.equals(fighter.sourcePosition) && fighter.isPayloadNotFull()
				},
				do: func(fighter *unit) {
					fighter.refillPayload()
					// TODO: when a missile terminates, it should remove the target from skip targets
					// rather than do this during restock. This should be an upgrade.
					// Also, multiple ordnance might make this complex. Just because a single ordnance skips
					// a target does not mean others should or can this be avoided?
					fighter.skipTargets = []id{}
				},
				delay: fighterConfig.reloadPayloadDelay,
			},
			{
				when: func(fighter *unit) bool { return fighter.nearest != nil },
				do: func(fighter *unit) {
					var uc unitConfig
					// TODO: determine fighter missile preference? Configurable by player?
					// Arbitrarily set radar to have higher priority than heat seeking.
					switch {
					case fighter.missileRadarActiveCount > 0 && !missileRadarActiveConfig.ignoreTarget(fighter.nearest):
						fighter.missileRadarActiveCount--
						uc = missileRadarActiveConfig
					case fighter.missileHeatSeekingCount > 0 && !missileHeatSeekingConfig.ignoreTarget(fighter.nearest):
						fighter.missileHeatSeekingCount--
						uc = missileHeatSeekingConfig
					default:
						return
					}
					spawn := newUnit(
						uc,
						position{layer: uc.produceLayer, Pixel: fighter.Pixel},
						fighter.player,
						[]position{fighter.nearest.position},
						[]action{},
						GridPosition{},
						fighter.selected)
					spawn.sourceUnit = fighter.id
					fighter.requests = append(fighter.requests, request{spawn: &spawn})
				},
				delay: fighterConfig.fireOrdnanceDelay,
			},
			reactChaff(fighterConfig.fireOrdnanceDelay),
			reactFlare(fighterConfig.fireOrdnanceDelay),
		}
	default:
		return []reaction{}
	}
}

func reactChaff(fireOrdnanceDelay uint64) reaction {
	return reaction{
		when: func(u *unit) bool {
			return u.nearest != nil && u.nearest.name == missileRadarActive && u.chaffCount > 0
		},
		do: func(u *unit) {
			chaff := newUnit(
				chaffConfig,
				position{layer: chaffConfig.produceLayer, Pixel: u.Pixel},
				u.player,
				[]position{u.position},
				[]action{},
				GridPosition{},
				u.selected)
			u.requests = append(u.requests, request{spawn: &chaff})
			u.chaffCount--
		},
		delay: fireOrdnanceDelay,
	}
}

func reactFlare(fireOrdnanceDelay uint64) reaction {
	return reaction{
		when: func(u *unit) bool {
			return u.nearest != nil && u.nearest.name == missileHeatSeeking && u.flareCount > 0
		},
		do: func(u *unit) {
			flare := newUnit(
				flareConfig,
				position{layer: flareConfig.produceLayer, Pixel: u.Pixel},
				u.player,
				[]position{u.nearest.position},
				[]action{},
				GridPosition{},
				u.selected)
			u.requests = append(u.requests, request{spawn: &flare})
			u.flareCount--
		},
		delay: fireOrdnanceDelay,
	}
}
