package generators

import (
	"math"
	"math/rand"

	"github.com/golangchallenge/gc6/mazelib"
)

var Animate func(m *mazelib.Maze) = nil

type possibility struct {
	dir   string
	coord mazelib.Coordinate
}

//Create a new Depth-First maze with the given bias.
//Valid biases are:
//"H": Prefer horizontal paths over vertical
//"V": Prefer vertical paths over horizontal
//"X": Choose cell furthest from goal. This will create a long winding path from start to finish.
//"O": Chose path closest to goal. This will create a direct path to treasure and the rest will be a disconnected component. Goal being to bait icarus into taking the wrong path and having to backtrack.
//"", or any other value:  Default, no bias. Always picka a random neighbor.
func DepthFirst(width, height int, bias string) *mazelib.Maze {
	m := mazelib.FullMaze(width, height)
	x, y := m.End() //search treasure -> icarus so treasure is usually in a dead end.
	startCoord := mazelib.Coordinate{x, y}
	visited := map[mazelib.Coordinate]bool{}
	visited[startCoord] = true
	current := []mazelib.Coordinate{startCoord}
	goalX, goalY := m.Start()
	for len(current) > 0 {
		possible := []possibility{}
		tip := current[len(current)-1]
		leftCoord := mazelib.Coordinate{tip.X - 1, tip.Y}
		if tip.X > 0 && !visited[leftCoord] {
			possible = append(possible, possibility{"left", leftCoord})
		}
		rightCoord := mazelib.Coordinate{tip.X + 1, tip.Y}
		if tip.X < width-1 && !visited[rightCoord] {
			possible = append(possible, possibility{"right", rightCoord})
		}
		upCoord := mazelib.Coordinate{tip.X, tip.Y - 1}
		if tip.Y > 0 && !visited[upCoord] {
			possible = append(possible, possibility{"up", upCoord})
		}
		downCoord := mazelib.Coordinate{tip.X, tip.Y + 1}
		if tip.Y < height-1 && !visited[downCoord] {
			possible = append(possible, possibility{"down", downCoord})
		}
		if len(possible) == 0 {
			current = current[:len(current)-1]
			continue
		}
		dir := randomDir(possible, tip.X, tip.Y, goalX, goalY, bias)
		newCoord := digInto(dir, tip, m)
		visited[newCoord] = true
		current = append(current, newCoord)
	}
	// Bias O is a direct path to treasure. Remove all other walls from start
	// so that everything not on the path will be connected, and likely fully explored before backtracking
	// to the right route. Also makes choosing correct path less likely.
	if bias == "O" {
		if goalX > 0 {
			digInto("left", mazelib.Coordinate{goalX, goalY}, m)
		}
		if goalX < width-1 {
			digInto("right", mazelib.Coordinate{goalX, goalY}, m)
		}
		if goalY > 0 {
			digInto("up", mazelib.Coordinate{goalX, goalY}, m)
		}
		if goalY < height-1 {
			digInto("down", mazelib.Coordinate{goalX, goalY}, m)
		}
	}
	return m
}

func randomDir(possible []possibility, x, y, avoidX, avoidY int, bias string) string {
	newPossible := possible
	increaseWeight := func(p possibility) {
		newPossible = append(newPossible, p)
		newPossible = append(newPossible, p)
		newPossible = append(newPossible, p)
	}
	if bias == "H" || bias == "V" {
		for _, p := range possible {
			if bias == "V" && (p.dir == "up" || p.dir == "down") {
				increaseWeight(p)
			} else if bias == "H" && (p.dir == "left" || p.dir == "right") {
				increaseWeight(p)
			}
		}
	} else if bias == "X" || bias == "O" {
		// find nearest and farthest neighbor from goal
		maxDist := float64(-1)
		maxAt := 0
		minDist := float64(5000)
		minAt := 0
		for i, p := range possible {
			distx := float64(p.coord.X - avoidX)
			distx *= distx
			disty := float64(p.coord.Y - avoidY)
			disty *= disty
			dist := math.Sqrt(distx + disty)
			if dist > maxDist {
				maxDist = dist
				maxAt = i
			}
			if dist < minDist {
				minDist = dist
				minAt = i
			}
		}
		if bias == "O" {
			return possible[minAt].dir
		}
		return possible[maxAt].dir
	}
	return newPossible[rand.Intn(len(newPossible))].dir
}

func digInto(dir string, current mazelib.Coordinate, m *mazelib.Maze) mazelib.Coordinate {
	var c mazelib.Coordinate
	switch dir {
	case "left":
		c = mazelib.Coordinate{current.X - 1, current.Y}
		roomA, _ := m.GetRoom(current.X, current.Y)
		roomB, _ := m.GetRoom(c.X, c.Y)
		roomA.RmWall(mazelib.W)
		roomB.RmWall(mazelib.E)
	case "right":
		c = mazelib.Coordinate{current.X + 1, current.Y}
		roomA, _ := m.GetRoom(current.X, current.Y)
		roomB, _ := m.GetRoom(c.X, c.Y)
		roomA.RmWall(mazelib.E)
		roomB.RmWall(mazelib.W)
	case "up":
		c = mazelib.Coordinate{current.X, current.Y - 1}
		roomA, _ := m.GetRoom(current.X, current.Y)
		roomB, _ := m.GetRoom(c.X, c.Y)
		roomA.RmWall(mazelib.N)
		roomB.RmWall(mazelib.S)
	case "down":
		c = mazelib.Coordinate{current.X, current.Y + 1}
		roomA, _ := m.GetRoom(current.X, current.Y)
		roomB, _ := m.GetRoom(c.X, c.Y)
		roomA.RmWall(mazelib.S)
		roomB.RmWall(mazelib.N)
	}
	if Animate != nil {
		Animate(m)
	}
	return c
}
