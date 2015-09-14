// Copyright © 2015 Steve Francia <spf@spf13.com>.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

// This is a small set of interfaces and utilities designed to help
// with the Go Challenge 6: Daedalus & Icarus

package mazelib

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// Coordinate describes a location in the maze
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Reply from the server to a request
type Reply struct {
	Survey  Survey `json:"survey"`
	Victory bool   `json:"victory"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

// Survey Given a location, survey surrounding locations
// True indicates a wall is present.
type Survey struct {
	Top    bool `json:"top"`
	Right  bool `json:"right"`
	Bottom bool `json:"bottom"`
	Left   bool `json:"left"`
}

const (
	N = 1
	S = 2
	E = 3
	W = 4
)

var ErrVictory error = errors.New("Victory")

// Room contains the minimum informaion about a room in the maze.
type Room struct {
	Treasure bool
	Start    bool
	Visited  bool
	Walls    Survey
}

func (r *Room) AddWall(dir int) {
	switch dir {
	case N:
		r.Walls.Top = true
	case S:
		r.Walls.Bottom = true
	case E:
		r.Walls.Right = true
	case W:
		r.Walls.Left = true
	}
}

func (r *Room) RmWall(dir int) {
	switch dir {
	case N:
		r.Walls.Top = false
	case S:
		r.Walls.Bottom = false
	case E:
		r.Walls.Right = false
	case W:
		r.Walls.Left = false
	}
}

// MazeI Interface
type MazeI interface {
	GetRoom(x, y int) (*Room, error)
	Width() int
	Height() int
	SetStartPoint(x, y int) error
	SetTreasure(x, y int) error
	LookAround() (Survey, error)
	Discover(x, y int) (Survey, error)
	Icarus() (x, y int)
	MoveLeft() error
	MoveRight() error
	MoveUp() error
	MoveDown() error
}

func AvgScores(in []int) int {
	if len(in) == 0 {
		return 0
	}

	var total int = 0

	for _, x := range in {
		total += x
	}
	return total / (len(in))
}

// PrintMaze : Function to Print Maze to Console
func PrintMaze(m MazeI) {
	fmt.Println("_" + strings.Repeat("__", m.Width()))
	for y := 0; y < m.Height(); y++ {
		str := ""
		for x := 0; x < m.Width(); x++ {
			if x == 0 {
				str += "|"
			}
			r, err := m.GetRoom(x, y)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			s, err := m.Discover(x, y)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			if s.Bottom {
				if r.Treasure {
					str += "T̲"
				} else if r.Start {
					str += "S̲"
				} else {
					str += "_"
				}
			} else {
				if r.Treasure {
					str += "T"
				} else if r.Start {
					str += "S"
				} else {
					str += " "
				}
			}

			if s.Right {
				str += "|"
			} else {
				str += " "
			}

		}
		fmt.Println(str)
	}
}

type Maze struct {
	rooms      [][]Room
	start      Coordinate
	end        Coordinate
	icarus     Coordinate
	StepsTaken int
}

// Return a room from the maze
func (m *Maze) GetRoom(x, y int) (*Room, error) {
	if x < 0 || y < 0 || x >= m.Width() || y >= m.Height() {
		return &Room{}, errors.New("room outside of maze boundaries")
	}

	return &m.rooms[y][x], nil
}

func (m *Maze) Width() int  { return len(m.rooms[0]) }
func (m *Maze) Height() int { return len(m.rooms) }

// Return Icarus's current position
func (m *Maze) Icarus() (x, y int) {
	return m.icarus.X, m.icarus.Y
}

func (m *Maze) End() (x, y int) {
	return m.end.X, m.end.Y
}

func (m *Maze) Start() (x, y int) {
	return m.start.X, m.start.Y
}

// Set the location where Icarus will awake
func (m *Maze) SetStartPoint(x, y int) error {
	r, err := m.GetRoom(x, y)

	if err != nil {
		return err
	}

	if r.Treasure {
		return errors.New("can't start in the treasure")
	}

	r.Start = true
	m.start = Coordinate{x, y}
	m.icarus = Coordinate{x, y}
	return nil
}

// Set the location of the treasure for a given maze
func (m *Maze) SetTreasure(x, y int) error {
	r, err := m.GetRoom(x, y)

	if err != nil {
		return err
	}

	if r.Start {
		return errors.New("can't have the treasure at the start")
	}

	r.Treasure = true
	m.end = Coordinate{x, y}
	return nil
}

// Given Icarus's current location, Discover that room
// Will return ErrVictory if Icarus is at the treasure.
func (m *Maze) LookAround() (Survey, error) {
	if m.end.X == m.icarus.X && m.end.Y == m.icarus.Y {
		fmt.Printf("Victory achieved in %d steps \n", m.StepsTaken)
		return Survey{}, ErrVictory
	}

	return m.Discover(m.icarus.X, m.icarus.Y)
}

// Given two points, survey the room.
// Will return error if two points are outside of the maze
func (m *Maze) Discover(x, y int) (Survey, error) {
	if r, err := m.GetRoom(x, y); err != nil {
		return Survey{}, nil
	} else {
		return r.Walls, nil
	}
}

// Moves Icarus's position left one step
// Will not permit moving through walls or out of the maze
func (m *Maze) MoveLeft() error {
	s, e := m.LookAround()
	if e != nil {
		return e
	}
	if s.Left {
		return errors.New("Can't walk through walls")
	}

	x, y := m.Icarus()
	if _, err := m.GetRoom(x-1, y); err != nil {
		return err
	}

	m.icarus = Coordinate{x - 1, y}
	m.StepsTaken++
	return nil
}

// Moves Icarus's position right one step
// Will not permit moving through walls or out of the maze
func (m *Maze) MoveRight() error {
	s, e := m.LookAround()
	if e != nil {
		return e
	}
	if s.Right {
		return errors.New("Can't walk through walls")
	}

	x, y := m.Icarus()
	if _, err := m.GetRoom(x+1, y); err != nil {
		return err
	}

	m.icarus = Coordinate{x + 1, y}
	m.StepsTaken++
	return nil
}

// Moves Icarus's position up one step
// Will not permit moving through walls or out of the maze
func (m *Maze) MoveUp() error {
	s, e := m.LookAround()
	if e != nil {
		return e
	}
	if s.Top {
		return errors.New("Can't walk through walls")
	}

	x, y := m.Icarus()
	if _, err := m.GetRoom(x, y-1); err != nil {
		return err
	}

	m.icarus = Coordinate{x, y - 1}
	m.StepsTaken++
	return nil
}

// Moves Icarus's position down one step
// Will not permit moving through walls or out of the maze
func (m *Maze) MoveDown() error {
	s, e := m.LookAround()
	if e != nil {
		return e
	}
	if s.Bottom {
		return errors.New("Can't walk through walls")
	}

	x, y := m.Icarus()
	if _, err := m.GetRoom(x, y+1); err != nil {
		return err
	}

	m.icarus = Coordinate{x, y + 1}
	m.StepsTaken++
	return nil
}

// Creates a maze without any walls
// Good starting point for additive algorithms
func EmptyMaze(xSize, ySize int) *Maze {
	z := Maze{}
	z.rooms = make([][]Room, ySize)
	for y := 0; y < ySize; y++ {
		z.rooms[y] = make([]Room, xSize)
		for x := 0; x < xSize; x++ {
			z.rooms[y][x] = Room{}
			if x == xSize-1 {
				z.rooms[y][x].AddWall(E)
			}
			if y == ySize-1 {
				z.rooms[y][x].AddWall(S)
			}
			if x == 0 {
				z.rooms[y][x].AddWall(W)
			}
			if y == 0 {
				z.rooms[y][x].AddWall(N)
			}
		}
	}
	z.RandomizeStartAndEnd()
	return &z
}

func (z *Maze) RandomizeStartAndEnd() {
	z.SetStartPoint(rand.Intn(z.Width()), rand.Intn(z.Height()))
	for {
		tX, tY := rand.Intn(z.Width()), rand.Intn(z.Height())
		if tX == z.start.X && tY == z.start.Y {
			continue
		}
		z.SetTreasure(tX, tY)
		break
	}
}

// Creates a maze with all walls
// Good starting point for subtractive algorithms
func FullMaze(xSize, ySize int) *Maze {
	z := EmptyMaze(xSize, ySize)

	for y := 0; y < ySize; y++ {
		for x := 0; x < xSize; x++ {
			z.rooms[y][x].Walls = Survey{true, true, true, true}
		}
	}
	return z
}
