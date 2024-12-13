package main

import (
	"fmt"
	"os"
	"strings"
)

type pair struct {
	i int
	j int
}

func main() {
	b, err := os.ReadFile("day_4_sample_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]
	lineLength := len(lines[0])

	grid := make([][]rune, len(lines))
	for i := range grid {
		grid[i] = make([]rune, lineLength)
	}

	for i, line := range lines {
		for j, r := range line {
			grid[i][j] = r
		}
	}

	sum := 0
	for i, row := range grid {
		for j, _ := range row {
			sum += check(grid, i, j)
		}
	}
	fmt.Println(sum)
}

func concat(texts ...string) string {
	return strings.Join(texts, "")
}

func check(grid [][]rune, i, j int) int {
	words := []string{}
	total := 0

	directions := []pair{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
		{1, 0},
		{1, -1},
		{0, -1},
	}

	for _, direction := range directions {
		words = append(words,
			concat(
				get(grid, i, j),
				get(grid, i+(1*direction.i), j+(1*direction.j)),
				get(grid, i+(2*direction.i), j+(2*direction.j)),
				get(grid, i+(3*direction.i), j+(3*direction.j)),
			),
		)
	}

	for _, word := range words {
		if word == "XMAS" {
			total++
		}
	}
	return total
}

func get(grid [][]rune, i, j int) string {
	if i < 0 || i >= len(grid) {
		return ""
	}
	row := grid[i]

	if j < 0 || j >= len(row) {
		return ""
	}

	return string(row[j])
}
