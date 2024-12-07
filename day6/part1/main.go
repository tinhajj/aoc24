package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Point struct {
	X int
	Y int
}

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

func main() {
	//b, err := os.ReadFile("day_6_sample_input.txt")
	b, err := os.ReadFile("day_6_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]

	area := make([][]string, len(lines)) // [y, x]
	visits := make([][]bool, len(lines)) // [y, x]
	start := Point{}

	rows := map[int][]Point{} // keep track of the blockers by row
	cols := map[int][]Point{} // keep track of the blockers by column

	for i, line := range lines {
		for j, r := range line {
			area[i] = append(area[i], string(r))
			visits[i] = append(visits[i], false)
			if string(r) == "^" {
				start = Point{j, i}
				visits[i][j] = true
			}
			if string(r) == "#" {
				rows[i] = append(rows[i], Point{i, j})
				cols[j] = append(cols[j], Point{i, j})
			}
		}
	}

	for k, v := range rows {
		sort.Slice(rows[k], func(i, j int) bool {
			return v[i].X < v[j].X
		})
	}
	for k, v := range cols {
		sort.Slice(cols[k], func(i, j int) bool {
			return v[i].Y < v[j].Y
		})
	}

	inbounds := true
	direction := UP
	pos := start

outer:
	for inbounds {
		switch direction {
		case UP:
			next := Point{pos.X, pos.Y - 1}
			if oob(area, next) {
				break outer
			}
			if area[next.Y][next.X] == "#" {
				direction = RIGHT
				continue
			}
			pos = next
			visits[next.Y][next.X] = true
		case DOWN:
			next := Point{pos.X, pos.Y + 1}
			if oob(area, next) {
				break outer
			}
			if area[next.Y][next.X] == "#" {
				direction = LEFT
				continue
			}
			pos = next
			visits[next.Y][next.X] = true
		case LEFT:
			next := Point{pos.X - 1, pos.Y}
			if oob(area, next) {
				break outer
			}
			if area[next.Y][next.X] == "#" {
				direction = UP
				continue
			}
			pos = next
			visits[next.Y][next.X] = true
		case RIGHT:
			next := Point{pos.X + 1, pos.Y}
			if oob(area, next) {
				break outer
			}
			if area[next.Y][next.X] == "#" {
				direction = DOWN
				continue
			}
			pos = next
			visits[next.Y][next.X] = true
		}
	}

	sum := 0
	for _, row := range visits {
		for _, val := range row {
			if val == true {
				sum++
			}
		}
	}

	fmt.Println("sum", sum)
}

func oob(area [][]string, p Point) bool {
	if p.Y >= len(area) || p.Y < 0 {
		return true
	}
	if p.X >= len(area[p.Y]) || p.X < 0 {
		return true
	}
	return false
}
