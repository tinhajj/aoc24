package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	debug  = true
	sample = false
)

type Vec2 struct {
	X, Y int
}

func (v Vec2) Add(o Vec2) Vec2 {
	return Vec2{
		X: v.X + o.X,
		Y: v.Y + o.Y,
	}
}

func main() {
	var b []byte
	var err error

	if sample {
		b, err = os.ReadFile("day_15_sample_input.txt")
	} else {
		b, err = os.ReadFile("day_15_input.txt")
	}
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[:len(lines)-1]

	grid := [][]string{}
	directions := []Vec2{}

	i := 0
	for ; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			break
		}
		row := []string{}
		for _, r := range line {
			row = append(row, string(r))
		}
		grid = append(grid, row)
	}

	for ; i < len(lines); i++ {
		line := lines[i]
		for _, r := range line {
			var direction Vec2
			switch string(r) {
			case "^":
				direction = Vec2{0, -1}
			case ">":
				direction = Vec2{1, 0}
			case "v":
				direction = Vec2{0, 1}
			case "<":
				direction = Vec2{-1, 0}
			}
			directions = append(directions, direction)
		}
	}

	var start Vec2
	for i, row := range grid {
		for j, val := range row {
			if val == "@" {
				start = Vec2{j, i}
			}
		}
	}

	currentPosition := start
	for _, direction := range directions {
		_, currentPosition = move(direction, currentPosition, grid)
	}

	sum := 0
	for y, row := range grid {
		for x, val := range row {
			if val != "O" {
				continue
			}
			sum += (y * 100) + x
		}
	}
	fmt.Println(sum)
}

func move(direction, position Vec2, grid [][]string) (ableToMove bool, next Vec2) {
	currentVal := grid[position.Y][position.X]
	nextPosition := position.Add(direction)

	if grid[nextPosition.Y][nextPosition.X] == "#" {
		return false, position
	}

	if grid[nextPosition.Y][nextPosition.X] == "." {
		grid[nextPosition.Y][nextPosition.X] = currentVal
		grid[position.Y][position.X] = "."
		return true, nextPosition
	}

	// must be "O", try and move it before we move this
	ok, _ := move(direction, nextPosition, grid)
	if ok {
		grid[nextPosition.Y][nextPosition.X] = currentVal
		grid[position.Y][position.X] = "."
		return true, nextPosition
	}
	return false, position
}
