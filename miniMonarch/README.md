# Mini Monarch

A prototype for a simple war/4x game. Simple like mini metro or mini motorways. 
War game like civ, hearts of iron, 9 kings etc. (also other popular war games on steamdb). 
Think about how this game relates to civilization revolution, catan, RISK, 
starcraft desert strike, totally accurate battle simulator, etc.
Graphics may end up similar to the gnorp analogue.
Is there a mobile equivalent for war/4x games?

## Checkpoints : Purpose

### next checkpoints

- Node highlighting 
  - Quicker testing
- Computer with dynamic targeting:
  - Can create more node setups for testing
- Update unit pathing algorithm to Bresenham:
  - Unit pathing is fundamental and will need lots of iteration.
  - Very likely this is needed for final game
- update pathing algorithm to fill in around obstacles:
  - More realistic pathing
  - Helps close attack behavior decision
- update pathing algorithm to go around obstacles:
  - More realistic pathing
  - Helps close attack behavior decision
- Determine behavior for multi hit melee combat 
  - Expose issues with combat
  - Expose issues with timing
  - Expose issues with pathing
  - Will probably lead to bug fix for units not killing correctly

### future checkpoints

- Tool / CLI which launches multiple windows at once to test multiple scenarios
  - Faster test loop
- Range unit which can attack:
  - First range unit will probably expose issues. 
  - Multiple units will be in needed for final game.
- Analyze and iterate on game design and game loop:
  - Check design feasability 
  - Improve game
- clean up bugs and hacks:
  - There are enough to cause bugs while testing game design and game loop
- determine if FPS lock is needed:
  - fundmental issue that I have not considered
- controls which work on computer, steam deck, phone, etc: 
  - multiple platform means more sales
  - customers can pick favorite device if multiple owned
- release to web:
  - learn how to deploy
  - easier to share with friends
  - test web based issues rather than only mac
- release to phone: 
  - check phone controls
  - check resolution
- local multiplayer:
  - ?
- save file: 
  - ?

### punted checkpoints

- initial map generation:
  - ?
- different terrains: 
  - ?

## Closed decisions

- Save center rather than top left of unit, nodes, etc into state because it makes the engine and game easier to design, test, debug, etc

## issues

- Nothing yet since early prototype stage 

## game design questions

- APM lock on nodes or units so you cannot spam actions? Think about sokpop game were you are limited by clicks
- square, hexagonal or historial battle grounds?
- Campaign mode. Specific battles, war or countries? Similar to cities in mini metro and mini motorway
- hidden or fully visible grid?
- naval units?
- air units?
- tech tree
- era progressions
- buildings
- other civ features
- Do I want to stand out from civ, or be a simpler civ e.g. factorio -> shapez. How does this affect marketing?
- Should players be able to target thier own nodes?
- Should a node be able to target itself?
- grids expand over time to cities, cities can merge together, merged cities are more powerful?
- cities have different units, different unit ratios, or different counts produced?
- traders can go to a node which you own to change what you are producing
- ammo limit for units
- nodes with defense system
- theme is the gnorp analogue?
- game name
- workers or resources that can be collected? This is starting to feel too RTS like.

## misc notes

 - Look into polytopia game
 - Look into other mobile "civ" games like unciv

## archive

- Don't do 3d globe visual for game. Too complex.