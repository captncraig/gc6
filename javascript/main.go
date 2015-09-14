package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golangchallenge/gc6/generators"
	"github.com/golangchallenge/gc6/mazelib"
	"github.com/golangchallenge/gc6/solvers"
	"honnef.co/go/js/dom"
)

//go:generate gopherjs build main.go
func main() {
	rand.Seed(time.Now().UnixNano())
	setupEvents()
	initialize()
	generators.Animate = AnimateGeneration
}

var currentContext *renderData

func initialize() {
	currentContext = &renderData{}
	currentContext.maze = idToGenerator()
	currentContext.solver = solvers.NewDFS()
	render(currentContext)
}

func idToGenerator() *mazelib.Maze {
	val := dom.GetWindow().Document().GetElementByID("generator").(*dom.HTMLSelectElement).Value
	switch val {
	case "dfs":
		generators.Bias = ""
		return generators.DepthFirst()
	case "dfs-h":
		generators.Bias = "H"
		return generators.DepthFirst()
	case "dfs-v":
		generators.Bias = "V"
		return generators.DepthFirst()
	case "dfs-x":
		generators.Bias = "X"
		return generators.DepthFirst()
	case "dfs-o":
		generators.Bias = "O"
		return generators.DepthFirst()
	case "empty":
		return mazelib.EmptyMaze(15, 10)
	}
	panic("unknown generator")
}

func setupEvents() {
	dom.GetWindow().Document().GetElementByID("init").
		AddEventListener("click", false, func(dom.Event) {
		go initialize()
	})
	dom.GetWindow().Document().GetElementByID("step").
		AddEventListener("click", false, func(dom.Event) {
		step()
	})
	dom.GetWindow().Document().GetElementByID("run").
		AddEventListener("click", false, func(dom.Event) {
		go run()
	})
}

func AnimateGeneration(m *mazelib.Maze) {
	render(&renderData{maze: m})
	time.Sleep(15 * time.Millisecond)
}

type renderData struct {
	maze   *mazelib.Maze
	solver solvers.MazeSolver
	count  int
}

const cellWidth int = 50

func step() error {
	c := currentContext
	surv, err := c.maze.LookAround()
	if err != nil {
		fmt.Println(err)
		return err
	}
	dir := c.solver.Step(surv)
	switch dir {
	case "left":
		err = c.maze.MoveLeft()
	case "right":
		err = c.maze.MoveRight()
	case "down":
		err = c.maze.MoveDown()
	case "up":
		err = c.maze.MoveUp()
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	currentContext.count++
	render(currentContext)
	return nil
}

func run() {
	var err error
	for err == nil {
		err = step()
		time.Sleep(20 * time.Millisecond)
	}
}

func render(c *renderData) {
	dom.GetWindow().Document().GetElementByID("stepCount").(*dom.HTMLSpanElement).SetTextContent(fmt.Sprint(c.count))
	ctx := dom.GetWindow().Document().GetElementByID("dc").(*dom.HTMLCanvasElement).GetContext2d()
	ctx.ClearRect(0, 0, 10000, 10000)
	startX, startY := c.maze.Start()
	endX, endY := c.maze.End()
	curX, curY := c.maze.Icarus()
	for y := 0; y < c.maze.Height(); y++ {
		for x := 0; x < c.maze.Width(); x++ {
			fillCell(ctx, x, y, "white")
			if x == curX && y == curY {
				fillCell(ctx, x, y, "pink")
			} else if x == startX && y == startY {
				fillCell(ctx, x, y, "orange")
			} else if x == endX && y == endY {
				fillCell(ctx, x, y, "yellow")
			}
			drawBorders(c, ctx, x, y)
		}
	}
}

func fillCell(ctx *dom.CanvasRenderingContext2D, x, y int, color string) {
	ctx.FillStyle = color
	ctx.FillRect(x*cellWidth+2, y*cellWidth+2, cellWidth-4, cellWidth-4)
}

func drawBorders(c *renderData, ctx *dom.CanvasRenderingContext2D, x, y int) {
	// 2 wide borders(4 on edge)
	ctx.FillStyle = "black"
	cell, _ := c.maze.GetRoom(x, y)
	if cell.Walls.Top {
		height := 2
		if y == 0 {
			height = 4
		}
		ctx.FillRect(x*cellWidth-2, y*cellWidth, cellWidth+4, height)
	}
	if cell.Walls.Bottom {
		height := 2
		if y == c.maze.Height()-1 {
			height = 4
		}
		ctx.FillRect(x*cellWidth-2, y*cellWidth+(cellWidth-height), cellWidth+4, height)
	}
	if cell.Walls.Left {
		width := 2
		if x == 0 {
			width = 4
		}
		ctx.FillRect(x*cellWidth, y*cellWidth-2, width, cellWidth+4)
	}
	if cell.Walls.Right {
		width := 2
		if x == c.maze.Width()-1 {
			width = 4
		}
		ctx.FillRect(x*cellWidth+(cellWidth-width), y*cellWidth-2, width, cellWidth+4)
	}
}
