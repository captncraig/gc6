package main

import (
	"github.com/golangchallenge/gc6/mazelib"
	"honnef.co/go/js/dom"
	"math/rand"
	"time"
)

//go:generate gopherjs build main.go
func main() {
	rand.Seed(time.Now().UnixNano())
	setupEvents()
	m := mazelib.EmptyMaze(15, 10)
	render(renderData{m})
}

func setupEvents() {
	initButton := dom.GetWindow().Document().GetElementByID("init")
	initButton.AddEventListener("click", false, func(dom.Event) {
		m := mazelib.EmptyMaze(15, 10)
		render(renderData{m})
	})
}

type renderData struct{ maze *mazelib.Maze }

const cellWidth int = 50

func render(c renderData) {
	ctx := dom.GetWindow().Document().GetElementByID("dc").(*dom.HTMLCanvasElement).GetContext2d()
	startX, startY := c.maze.Start()
	endX, endY := c.maze.End()
	for y := 0; y < c.maze.Height(); y++ {
		for x := 0; x < c.maze.Width(); x++ {
			fillCell(ctx, x, y, "white")
			if x == startX && y == startY {
				fillCell(ctx, x, y, "orange")
			}
			if x == endX && y == endY {
				fillCell(ctx, x, y, "yellow")
			}
			drawBorders(c, ctx, x, y)
		}
	}
}

func fillCell(ctx *dom.CanvasRenderingContext2D, x, y int, color string) {
	ctx.FillStyle = color
	ctx.FillRect(x*cellWidth, y*cellWidth, cellWidth, cellWidth)
}

func drawBorders(c renderData, ctx *dom.CanvasRenderingContext2D, x, y int) {
	// 2 wide borders(4 on edge)
	ctx.FillStyle = "black"
	cell, _ := c.maze.GetRoom(x, y)
	if cell.Walls.Top {
		height := 2
		if y == 0 {
			height = 4
		}
		ctx.FillRect(x*cellWidth, y*cellWidth, cellWidth, height)
	}
	if cell.Walls.Bottom {
		height := 2
		if y == c.maze.Height()-1 {
			height = 4
		}
		ctx.FillRect(x*cellWidth, y*cellWidth+(cellWidth-height), cellWidth, height)
	}
	if cell.Walls.Left {
		width := 2
		if x == 0 {
			width = 4
		}
		ctx.FillRect(x*cellWidth, y*cellWidth, width, cellWidth)
	}
	if cell.Walls.Right {
		width := 2
		if x == c.maze.Width()-1 {
			width = 4
		}
		ctx.FillRect(x*cellWidth+(cellWidth-width), y*cellWidth, width, cellWidth)
	}
}
