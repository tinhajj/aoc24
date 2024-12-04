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
	b, err := os.ReadFile("day_4_input.txt")
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
			if check(grid, i, j) {
				sum++
			}
		}
	}
	fmt.Println(sum)
}

func concat(texts ...string) string {
	return strings.Join(texts, "")
}

func check(grid [][]rune, i, j int) bool {
	words := []string{}
	total := 0

	words = append(words,
		concat(
			get(grid, i-1, j-1),
			get(grid, i, j),
			get(grid, i+1, j+1),
		),
	)
	words = append(words,
		concat(
			get(grid, i+1, j-1),
			get(grid, i, j),
			get(grid, i-1, j+1),
		),
	)

	for _, word := range words {
		if word == "MAS" || word == "SAM" {
			total++
		}
	}
	return total == 2
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
