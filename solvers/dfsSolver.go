package solvers

import (
	"github.com/golangchallenge/gc6/mazelib"
	"math/rand"
)

type dfs struct {
	//lookup to see if a given coordinate has been visited.
	visited map[mazelib.Coordinate]bool
	//current path. current segment is last element.
	current []*dfsSegment
}

type dfsSegment struct {
	coord mazelib.Coordinate //coordinate of this cell
	dir   string             // the direction I moved to get here. Empty if first cell
}

func NewDFS() MazeSolver {
	zero := mazelib.Coordinate{}
	return &dfs{
		map[mazelib.Coordinate]bool{zero: true},
		[]*dfsSegment{{}},
	}
}

func (d *dfs) Step(s mazelib.Survey) string {
	presentCell := d.current[len(d.current)-1]
	x, y := presentCell.coord.X, presentCell.coord.Y
	possibleDirections := []*dfsSegment{}
	leftCoord := mazelib.Coordinate{x - 1, y}
	rightCoord := mazelib.Coordinate{x + 1, y}
	upCoord := mazelib.Coordinate{x, y - 1}
	downCoord := mazelib.Coordinate{x, y + 1}
	if !s.Left && !d.visited[leftCoord] {
		possibleDirections = append(possibleDirections, &dfsSegment{leftCoord, "left"})
	}
	if !s.Right && !d.visited[rightCoord] {
		possibleDirections = append(possibleDirections, &dfsSegment{rightCoord, "right"})
	}
	if !s.Top && !d.visited[upCoord] {
		possibleDirections = append(possibleDirections, &dfsSegment{upCoord, "up"})
	}
	if !s.Bottom && !d.visited[downCoord] {
		possibleDirections = append(possibleDirections, &dfsSegment{downCoord, "down"})
	}
	if len(possibleDirections) == 0 {
		d.current = d.current[:len(d.current)-1]
		return reverseDir(presentCell.dir)
	}
	chosen := possibleDirections[rand.Intn(len(possibleDirections))]
	d.current = append(d.current, chosen)
	d.visited[chosen.coord] = true
	return chosen.dir
}

func reverseDir(dir string) string {
	switch dir {
	case "left":
		return "right"
	case "right":
		return "left"
	case "up":
		return "down"
	case "down":
		return "up"
	}
	panic("can't reverse unknown dir")
}
