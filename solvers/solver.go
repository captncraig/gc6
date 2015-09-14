package solvers

import (
	"github.com/golangchallenge/gc6/mazelib"
)

type MazeSolver interface {
	Step(mazelib.Survey) string
}
