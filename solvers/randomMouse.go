package solvers

import (
	"github.com/golangchallenge/gc6/mazelib"
	"math/rand"
)

type mouse struct {
	lastDir string
}

func NewMouse() MazeSolver {
	return &mouse{""}
}

// move a random direction, giving last preference to the direction I just came from
func (m *mouse) Step(s mazelib.Survey) string {
	tentative := ""
	dirs := []string{}
	if !s.Bottom {
		if m.lastDir != "up" {
			dirs = append(dirs, "down")
		} else {
			tentative = "down"
		}
	}
	if !s.Top {
		if m.lastDir != "down" {
			dirs = append(dirs, "up")
		} else {
			tentative = "up"
		}
	}
	if !s.Left {
		if m.lastDir != "right" {
			dirs = append(dirs, "left")
		} else {
			tentative = "left"
		}
	}
	if !s.Right {
		if m.lastDir != "left" {
			dirs = append(dirs, "right")
		} else {
			tentative = "right"
		}
	}
	if tentative != "" && len(dirs) == 0 {
		dirs = append(dirs, tentative)
	}
	m.lastDir = dirs[rand.Intn(len(dirs))]
	return m.lastDir
}
