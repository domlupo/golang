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

- Research how to pretty log structs in go in a scaleable manner
  - This can be a huge time sink. Be mindful here and a hacky solution may be best.
  - Should I configure logger to have two modes? Human readable and Machine readable?
  - Fix unit logging where I do "unit:refresh:before:". I probably should be using key values pairs instead.
  - Use logger key values to pretty print unit
  - Will improve logging and debugging iteration
- Determine reach behavior:
  - Current behavior seems strange. I am not sure why.
  - I think reach behavers different for unit attacking node and unit attaking other unit
  - Should units priortize nodes or other units?
  - Determine if units should always attack closest unit
  - Units don't currently attack closest unit. This may complete that TODO.
  - vim search for "nearest"
  - Also vim search "TODO" because this solve several
- Add missle:
  - first node unit
  - splash damage?
  - no health?
  - flight?
- Improve flux debugging
  - Add flux square selection
  - Add selected flux logging for detailed tracing
- Create more node setup test scenarios
  - Able to easily reproduce and test precise setups ofr edge case issues
  - Able to randomly generated setups to find hidden edge case issues
- Analyze and iterate on game design and game loop:
  - Check design feasability 
  - Improve game
- clean up bugs and hacks:
  - There are enough to cause bugs while testing game design and game loop

---

### Future Checkpoints

- Research and iterate local pathfinding algorithms/methods
  - More realistic short distance pathing
  - Proper surronding of a target
  - Proper Units bunching into a single spot
  - Helps close attack behavior decision
- Determine behavior for multi hit melee combat 
  - Expose issues with combat
  - Expose issues with timing
  - Expose issues with pathing
  - Will probably lead to bug fix for units not killing correctly
  - May help or hurt efforts towards local pathfinding. I suspect this will be coupled with local pathfinding.
- Tool / CLI which launches multiple windows at once to test multiple scenarios
  - Faster test loop
- save file: 
  - Will expose constraints
  - Will expose platform concerns
  - May expose issues with work so far 
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
  - Initial research shows webRTC is probably better but more complex than websockets
- Polish logger
  - Add logs throughout (global path, local path, attack, move, collision, etc)
  - Add runtime dynamic log sampling rates. Is this even possible? Are the benefits worth it?
  - Paused debugging? e.g. Updating unit during state.go execution can log state at game level.
- Investigate if logging impacts game performance
  - Might be able to fix via build tags to compile different code for release build

---

### Punted Checkpoints

- Update unit pathing algorithm to Theta* with same cost orthgonal and diagonal movement: (https://en.wikipedia.org/wiki/Theta*)
  - Unit pathing is fundamental and will need lots of iteration.
  - More realistic long distance pathing rather than current angled
  - Units will be able to go around obstacles
  - Very likely this is needed for final game
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

- Should there be a game camera:
  - No camera can easily be the best option
  - One option is starcraft / civ where it panes around a larger area
  - One option is mini metro and mini motorways style that slowly zooms out to the full area?
- Should the focus of the game be positioning, logistics and decisions rather than micro?
  - It seems interesting whens units get into long term fights
  - Reminds me of the board game go RISK and Catan
  - Reminds me of the military diagrams that show battles over time
  - May be able to utilize the idea of target change units
  - Early fight for position, long term fight for ?
- Grid is square or hexagonal? A potential issue with hexagonal is more CPU for certain algorithms.
- top down graphics or iso or 3d?
- Players get to pick starting positions?
- Hidden or fully visible grid?
- APM lock on nodes or units so you cannot spam actions? Think about sokpop game clickyland where you are limited by clicks.
- Campaign mode. Specific battles, war or countries? Similar to cities in mini metro and mini motorway.
- Tech tree?
- Special nodes?
- Should nodes only have a single grid offset (offset = 0 is acceptable) or should offset differ? 
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
- Replay feature where you can see a battle recap?
- Should unit behavior be toggable on a global scale? E.g. tanks have fire mode vs push mode

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
    - unit to change targets
    - a unit which pushes other units out of way
  - Air:
    - cargo plane (can transfer production between nodes?)
    - paratrooper
    - bomber
    - stealth bomber
    - fighter
    - stealth fighter
    - helicopter / gun ship
    - drone
    - unit to change targets?
  - City aka Node:
    - missle
    - guided missle
    - nuke
    - anti air gun
    - satelite
    - unit to change targets?

---