package solvers

import (
	"github.com/golangchallenge/gc6/mazelib"
)

// In order to make our solvers usable in non-networked environments like tests or
// deep analysis, we will make their move functions injectable.
// Also alows dom animation if solver used from gopherjs.
type MoveFunc func(direction string) (mazelib.Survey, error)

type MazeSolver interface {
	Solve(mazelib.Survey)
}
