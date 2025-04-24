package main

import (
	"aoc24/scan"
	"fmt"
	"math"
	"os"
	"strings"
)

var (
	SAMPLE          = false
	WIDTH           int
	HEIGHT          int
	SIMULATIONCOUNT int
)

type Corruption struct {
	X int
	Y int
}

type Point struct {
	X       int
	Y       int
	Corrupt bool
}

func (p *Point) Add(v Point) Point {
	return Point{
		X: v.X + p.X,
		Y: v.Y + p.Y,
	}
}

func main() {
	var err error
	var b []byte

	if SAMPLE {
		b, err = os.ReadFile("day_18_sample_input.txt")
		HEIGHT = 6
		WIDTH = 6
		SIMULATIONCOUNT = 12
	} else {
		b, err = os.ReadFile("day_18_input.txt")
		HEIGHT = 70
		WIDTH = 70
		SIMULATIONCOUNT = 1024
	}
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(b), "\n")
	lines = lines[:len(lines)-1]

	corruptions := []Corruption{}

	for _, line := range lines {
		numbers := scan.Numbers(line)
		corruptions = append(corruptions, Corruption{X: numbers[0], Y: numbers[1]})
	}

	grid := [][]*Point{}
	adj := map[*Point][]*Point{}

	for y := 0; y <= HEIGHT; y++ {
		grid = append(grid, []*Point{})
		for x := 0; x <= WIDTH; x++ {
			grid[y] = append(grid[y], &Point{X: x, Y: y})
		}
	}

	for y := 0; y <= HEIGHT; y++ {
		for x := 0; x <= WIDTH; x++ {
			current := grid[y][x]

			arounds := Around(current, grid)
			for _, other := range arounds {
				adj[current] = append(adj[current], other)
			}
		}
	}

	for i, corruption := range corruptions {
		grid[corruption.Y][corruption.X].Corrupt = true
		dist, _ := Dijkstra(grid, adj, grid[0][0])

		corner := grid[HEIGHT][WIDTH]

		if dist[corner] == math.MaxInt {
			fmt.Println(corruption)
			break
		}
		fmt.Printf("%d / %d\n", i, len(corruptions))
	}
}

func Dijkstra(grid [][]*Point, adjs map[*Point][]*Point, start *Point) (map[*Point]int, map[*Point]*Point) {
	distance := map[*Point]int{}
	previous := map[*Point]*Point{}

	queue := []*Point{}

	for _, y := range grid {
		for _, p := range y {
			distance[p] = math.MaxInt
			previous[p] = nil
			queue = append(queue, p)
		}
	}
	distance[start] = 0

	for len(queue) > 0 {
		// find min
		i := 0
		min := queue[i]
		for j, q := range queue {
			if distance[q] < distance[min] {
				min = q
				i = j
			}
		}

		// remove min
		queue = append(queue[:i], queue[i+1:]...)

		for _, adj := range adjs[min] {
			inQueue := false
			isCorrupt := adj.Corrupt

			for _, q := range queue {
				if q == adj {
					inQueue = true
				}
			}
			if !inQueue || isCorrupt {
				continue
			}

			alt := distance[min] + 1
			if alt < distance[adj] {
				distance[adj] = alt
				previous[adj] = min
			}
		}
	}

	return distance, previous
}

func Around(p *Point, grid [][]*Point) []*Point {
	result := []*Point{}
	directions := []Point{
		Point{0, -1, false},
		Point{-1, 0, false},
		Point{1, 0, false},
		Point{0, 1, false},
	}

	for _, direction := range directions {
		other := p.Add(direction)

		if other.X < 0 || other.Y < 0 {
			continue
		}

		if other.X > len(grid[0])-1 || other.Y > len(grid)-1 {
			continue
		}

		op := grid[other.Y][other.X]
		result = append(result, op)
	}

	return result
}
