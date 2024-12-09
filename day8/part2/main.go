package main

import (
	"fmt"
	"os"
	"strings"
)

type Point struct {
	Y int
	X int
}

func main() {
	b, err := os.ReadFile("day_8_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]

	groups := map[string][]Point{}

	for y, line := range lines {
		for x, char := range line {
			if string(char) == "." {
				continue
			}
			groups[string(char)] = append(groups[string(char)], Point{Y: y, X: x})
		}
	}

	mapHeight := len(lines)
	mapWidth := len(lines[0])

	antinodes := []Point{}

	for _, antennas := range groups {
		for i := 0; i < len(antennas)-1; i++ {
			for j := i + 1; j < len(antennas); j++ {
				first := antennas[i]
				second := antennas[j]

				antinodes = append(antinodes, first)
				antinodes = append(antinodes, second)

				distY := second.Y - first.Y
				distX := second.X - first.X

				antinode1 := Point{Y: second.Y + distY, X: second.X + distX}
				for antinode1.X >= 0 && antinode1.X < mapWidth && antinode1.Y >= 0 && antinode1.Y < mapHeight {
					antinodes = append(antinodes, antinode1)
					antinode1 = Point{Y: antinode1.Y + distY, X: antinode1.X + distX}
				}

				distX = distX * -1
				distY = distY * -1

				antinode2 := Point{Y: first.Y + distY, X: first.X + distX}
				for antinode2.X >= 0 && antinode2.X < mapWidth && antinode2.Y >= 0 && antinode2.Y < mapHeight {
					antinodes = append(antinodes, antinode2)
					antinode2 = Point{Y: antinode2.Y + distY, X: antinode2.X + distX}
				}
			}
		}
	}

	seen := map[Point]int{}
	for _, antinode := range antinodes {
		seen[antinode]++
	}

	fmt.Println(len(seen))
}
