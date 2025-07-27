## Code reflection

These are parts of the code I am unhappy with but they are not worth 
improving. Hopefully this reminds me to avoid these issues in future.

---

The distribute card code ported from Lua to Golang is gross but works. 
This is first time I used Golang so I need to invest time learning it.

The high level game logic is the entire update function. It is only 
150 LOC but it is a bit hard to follow. Breaking into smaller 
functions might improve the code. This would be required for games
with more complex logic. A compromise is comments which guide readers. 

When porting from Lua to Golang I had to add an "active" field to 
card struct. Love2d enegine had an elegant hack which avoided needing 
this field. I like how at the end of the update function it is very 
easy to determine which cards need to be activated or deactivated. 
This made be notice that other code could be moved at the end of game
logic to ensure state is correct. e.g card height.

What seems wrong now is that CardPile has a cardCountAffectsHeight 
field which is used when adding and removing cards to update height. 
This field is required in unit tests which only test game logic, but
never test card height, so it leaks in test which is wrong. Also, in 
a production game it might be better to be able to test game logic 
and game display separately. Even though they are almost always 
coupled. It might be better to remove cardCountAffectHeights from 
cardPile.go and move logic to solitaire.go update function. Overall,
cardCountAffectsHeight seems a bit misplaced.

Unit tests are sparse because:
- The code is ported
- It has extensive manual testing
- I don't expect to refactor this code ever
