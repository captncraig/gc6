package generators

import (
	"github.com/golangchallenge/gc6/mazelib"
)

func DepthFirst() mazelib.MazeI {
	m := mazelib.FullMaze(15, 10)
	return m
}
