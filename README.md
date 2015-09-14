### The Go Challenge 6

#### Icarus

My Icarus is a fairly straightforward depth-first search implementation. It will never choose a visited neighbor, and will backtrack if it hits a dead-end.

#### Daedalus

My maze generator assumes the opponent will use a similar backtracking dfs approach to what I did. It attempts to force icarus to backtrack as much as possible by:

1. Making as short a path as possible to the treasure. 
2. Attempting to make it unlikely icarus will choose the correct path from the starting cell.

The hope is that icarus will choose a wrong path initially, and then have to backtrack the entire map before finally trying the correct path. The best case is where icarus starts right next to the treasure, but chooses one of the other 3 directions, which are all connected to each other. HE must then traverse the entire wrong component of the map before backtracking to the start and finding the treasure.

If we were doing perfect-knowledge solving, this would be an extremely poor strategy, but with icarus only seeing his current room, it works ok.

### Cool Visualization

I used this as an opportunity to play with gopherjs. I made a simple page to visualize my various generator options and solvers, complete with animation. Demo can be viewed at http://captncraig.github.io/gc6/javascript/test.html