package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	debug   = false
	sample  = false
	Vec2Nil = Vec2{-1, -1}
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
			switch r {
			case '#':
				row = append(row, "#")
				row = append(row, "#")
			case 'O':
				row = append(row, "[")
				row = append(row, "]")
			case '.':
				row = append(row, ".")
				row = append(row, ".")
			case '@':
				row = append(row, "@")
				row = append(row, ".")
			}
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

	if debug {
		fmt.Println("start")
		print(grid)
	}
	currentPosition := start

	for _, direction := range directions {
		if debug {
			fmt.Println("step", direction)
		}

		ok := canMove(direction, currentPosition, map[Vec2]bool{}, grid)
		if ok {
			currentPosition = move(direction, currentPosition, map[Vec2]bool{}, grid)
		}

		if debug {
			print(grid)
		}
	}

	sum := 0
	for y, row := range grid {
		for x, val := range row {
			if val != "[" {
				continue
			}
			sum += (y * 100) + x
		}
	}
	fmt.Println(sum)
}

func print(grid [][]string) {
	for _, row := range grid {
		fmt.Println(row)
	}
}

func move(direction, position Vec2, memo map[Vec2]bool, grid [][]string) Vec2 {
	_, ok := memo[position]
	if ok {
		return Vec2{}
	}

	currentValue := grid[position.Y][position.X]
	nextPosition := position.Add(direction)

	if grid[nextPosition.Y][nextPosition.X] == "." {
		grid[nextPosition.Y][nextPosition.X] = currentValue
		grid[position.Y][position.X] = "."

		memo[position] = true
		return nextPosition
	}

	// must be "[,]", try and move it before we move this
	nextVal := grid[nextPosition.Y][nextPosition.X]
	var otherSide Vec2

	if nextVal == "[" {
		otherSide = nextPosition.Add(Vec2{1, 0})
	} else {
		otherSide = nextPosition.Add(Vec2{-1, 0})
	}

	move(direction, nextPosition, memo, grid)
	move(direction, otherSide, memo, grid)

	grid[nextPosition.Y][nextPosition.X] = currentValue
	grid[position.Y][position.X] = "."

	memo[position] = true
	memo[nextPosition] = true
	memo[otherSide] = true

	return nextPosition
}

func canMove(direction, position Vec2, memo map[Vec2]bool, grid [][]string) (ableToMove bool) {
	v, ok := memo[position]
	if ok {
		return v
	}

	nextPosition := position.Add(direction)

	if grid[nextPosition.Y][nextPosition.X] == "#" {
		memo[position] = false
		return false
	}

	if grid[nextPosition.Y][nextPosition.X] == "." {
		memo[position] = true
		return true
	}

	// must be "[,]"
	nextVal := grid[nextPosition.Y][nextPosition.X]
	var otherSide Vec2 = Vec2Nil

	if nextVal == "[" && direction.Y != 0 {
		otherSide = nextPosition.Add(Vec2{1, 0})
	} else if nextVal == "]" && direction.Y != 0 {
		otherSide = nextPosition.Add(Vec2{-1, 0})
	}

	var okOther bool = true

	ok = canMove(direction, nextPosition, memo, grid)
	memo[nextPosition] = ok

	if otherSide != Vec2Nil {
		okOther = canMove(direction, otherSide, memo, grid)
		memo[otherSide] = ok
	}

	return ok && okOther
}
