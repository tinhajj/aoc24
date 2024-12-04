package main

import (
	"fmt"
	"os"
	"strings"
)

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
	fmt.Println(grid)
}
