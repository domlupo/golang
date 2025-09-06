# Mini Monarch

A prototype for a simple war/4x game. Simple like mini metro or mini motorways. 
War game like civ, hearts of iron, 9 kings etc. (also other popular war games on steamdb). 
Think about how this game relates to civilization revolution, catan, RISK, 
starcraft desert strike, totally accurate battle simulator, etc.
Graphics may end up similar to the gnorp analogue.
Is there a mobile equivalent for war/4x games?

---

## Checkpoints : Purposes

### Next Checkpoints

- Computer with dynamic targeting:
  - Can create additional precise node setups for testing
  - Can create randomly generated node setups for testing
  - Pathfinding algorithm practice
- Range unit which can attack:
  - First range unit will probably expose issues. 
  - Multiple units will be in needed for final game.
- Update unit pathing algorithm to Theta* with same cost orthgonal and diagonal movement: (https://en.wikipedia.org/wiki/Theta*)
  - Unit pathing is fundamental and will need lots of iteration.
  - More realistic long distance pathing rather than current angled
  - Units will be able to go around obstacles
  - Very likely this is needed for final game
- Determine behavior for multi hit melee combat 
  - Expose issues with combat
  - Expose issues with timing
  - Expose issues with pathing
  - Will probably lead to bug fix for units not killing correctly
  - May help or hurt efforts towards local pathfinding. I suspect this will be coupled with local pathfinding.
- save file: 
  - Will expose constraints
  - Will expose platform concerns
  - May expose issues with work so far 
- Analyze and iterate on game design and game loop:
  - Check design feasability 
  - Improve game
- clean up bugs and hacks:
  - There are enough to cause bugs while testing game design and game loop

---

### Future Checkpoints

- Once Theta* is utilized, research how to update pathing algorithm for short distances
  - More realistic short distance pathing
  - Proper surronding of a target
  - Proper Units bunching into a single spot
  - Helps close attack behavior decision
- Tool / CLI which launches multiple windows at once to test multiple scenarios
  - Faster test loop
- Determine if FPS lock is needed:
  - fundmental issue that I have not considered
- Controls which work on computer, steam deck, phone, etc: 
  - multiple platform means more sales
  - customers can pick favorite device if multiple owned
- Release to web:
  - learn how to deploy
  - easier to share with friends
  - test web based issues rather than only mac
- Release to phone: 
  - check phone controls
  - check resolution
- local multiplayer:
  - ?

---

### Punted Checkpoints

- Research more global pathfinding algorithms e.g. AP Theta*
  - ?
- Research more local pathfinding algorithms
  - ?
- initial map generation:
  - ?
- different terrains: 
  - Careful, this may be hard to implement due to pathfinding algorithms assuming same cost movement

---

## Commands

```find . -type f -name "*.go" -exec wc -l {} +```

```go tool pprof -http=:6060 cpu.prof```

---

## Closed Decisions

1. No worker unit or resource collection
    - Adds complexity
    - Nodes can produce cargo units which more aligned with multi use nodes
    - No guaranteed benefit
2. No civilization era progressions:
    - Directly adds design and balancing complexity
    - Indirectly adds complexity and scope throughout game e.g. more units
    - No guaranteed benefit
3. Save center rather than top left of unit, nodes, etc into state:
   - Easier to design game/engine
   - Easier to test game/engine
   - Easier to debug game/engine
4. Same cost Orgthogonal and Cardinal movement:
   - Simplifies path finding
   - Lowers CPU cost
   - Lowers overall game time
5. Don't do 3d globe visual for game:
   - Very complex
   - No guaranteed benefit
   - Performance concerns

---

## Issues and Concerns

- Nothing yet since early prototype stage 

---

## Game Design Questions

- Grid is square or hexagonal?
- top down graphics or iso or 3d?
- Players get to pick starting positions?
- Hidden or fully visible grid?
- APM lock on nodes or units so you cannot spam actions? Think about sokpop game clickyland where you are limited by clicks.
- Campaign mode. Specific battles, war or countries? Similar to cities in mini metro and mini motorway.
- Tech tree?
- Special nodes?
- Node expansion overtime?
- Grids expand over time to cities, cities can merge together, merged cities are more powerful?
- Buildings?
- Do I want to stand out from civ, or be a simpler civ e.g. factorio -> shapez. How does this affect marketing?
- Should players be able to target thier own nodes?
- Should a node be able to target itself?
- Cities have different units, different unit ratios, or different counts produced? e.g. Catan
- Traders or cargo can go to a node which you own to transfer production of city production
- Ammo limit for units?
- Nodes with defense system
- How is node unit production handled?
- Theme is the gnorp analogue?
- Game name?

---

## Misc Notes

- Research civ concepts and features e.g. culture victory, trade, etc
- Research polytopia game
- Research other mobile platform "civ" games like unciv

---

### Unit List Brainstorm

- Conclusions:
  - Modern era is better than medieval era because:
    - Variety
    - Familiarity
  - Naval units are unneeded and huge scope creep.

- Questions:
  - Which units will counter other units?
  - Should units have upgrades?
  - Should certain "civs" have better units?
  - Should planes need to fly back to airport to refuel?
  - Should planes need to fly back to airport to rearm?
  - Should planes need to fly back to airport to repair?
  - Should plane matience be much quicker than building a new plane?
  - How can other units be unique and fit with node mechanic?
  - Should I do purpose exercise with units like I do with checkpoints

- Potential Units:
  - Land:
    - marine (cross water?)
    - infantry
    - sniper
    - anti tank infantry
    - tank
    - artillery 
    - mobile SAM
    - troop transport
    - cargo truck (can transfer production between nodes?)
  - Air:
    - cargo plane (can transfer production between nodes?)
    - paratrooper
    - bomber
    - stealth bomber
    - fighter
    - stealth fighter
    - helicopter / gun ship
    - drone
  - City aka Node:
    - missle
    - guided missle
    - nuke
    - anti air gun
    - satelite

---