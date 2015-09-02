package solvers

import (
	"fmt"
	"github.com/golangchallenge/gc6/mazelib"
	"math/rand"
)

type mouse struct {
	move MoveFunc
}

func NewMouse(mov MoveFunc) MazeSolver {
	return &mouse{mov}
}

// move a random direction, giving last preference to the direction I just came from
func (m *mouse) Solve(s mazelib.Survey) {
	var e error
	lastDir := ""
	for {
		tentative := ""
		dirs := []string{}
		if !s.Bottom {
			if lastDir != "up" {
				dirs = append(dirs, "down")
			} else {
				tentative = "down"
			}
		}
		if !s.Top {
			if lastDir != "down" {
				dirs = append(dirs, "up")
			} else {
				tentative = "up"
			}
		}
		if !s.Left {
			if lastDir != "right" {
				dirs = append(dirs, "left")
			} else {
				tentative = "left"
			}
		}
		if !s.Right {
			if lastDir != "left" {
				dirs = append(dirs, "right")
			} else {
				tentative = "right"
			}
		}
		if tentative != "" && len(dirs) == 0 {
			dirs = append(dirs, tentative)
		}
		lastDir = dirs[rand.Intn(len(dirs))]
		s, e = m.move(lastDir)
		if e != nil && e.Error() != "" {
			fmt.Println("!!!", e)
			return
		}
	}
}
