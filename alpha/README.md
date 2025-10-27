# Alphabet Soup

A portfolio website that mimics alphabet soup.

- Random Letters will float around the screen with steering behavior (boids) and collision.
- Letters will be cohesive to form links (github, linked, etc)

## Unknowns

- Determine if I need to use shaders
- Ebitengine web release

## Theme

Invest in the soup theme?

- A spoon is used to interact with the links
- A spoon that can interact with the soup:
  - swirl, pick up random letters, etc
  - Spoon starts in soup and moves from mouse input
  - Left click lifts spoon from soup, releasing left click (no input) returns spoon to default position into soup
  - Scooping a link is equivalent to clicking a hyperlink
- Letters can sink and surface in soup
- Can spill soup onto side of bowl

## Programming

A spoonful of soup could be modeled as a second "bowl" of soup. I would then need to handle
merge and splits of soup.

Open source so others can use for their portfolio websites