package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

var (
	SAMPLE = false
)

type CheatJump struct {
	Start *Point
	End   *Point
}

type Point struct {
	X   int
	Y   int
	Val string
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
		b, err = os.ReadFile("day_20_sample_input.txt")
	} else {
		b, err = os.ReadFile("day_20_input.txt")
	}
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(b), "\n")
	lines = lines[:len(lines)-1]

	grid := [][]*Point{}
	adjs := map[*Point][]*Point{}

	var start, end *Point
	_, _ = start, end

	for y, line := range lines {
		grid = append(grid, []*Point{})
		for x, c := range line {
			val := string(c)
			p := &Point{X: x, Y: y, Val: val}
			if val == "S" {
				start = p
			}
			if val == "E" {
				end = p
			}
			grid[y] = append(grid[y], p)
		}
	}

	for _, row := range grid {
		for _, p := range row {
			around := Around(p, grid)
			adjs[p] = around
		}
	}

	dist, prev := Dijkstra(grid, adjs, end)
	_, _ = dist, prev

	path := []*Point{start}

	backtrack := prev[start]
	for backtrack != nil {
		path = append(path, backtrack)
		backtrack = prev[backtrack]
	}

	jumpSavings := map[CheatJump]int{}

	for _, p := range path {
		distBefore := dist[p]
		jumps := CheatJumps(p, grid)
		for k, v := range jumps {
			_, ok := jumpSavings[k]
			if ok {
				continue
			}
			distAfter := dist[k.End]
			distSaved := distBefore - distAfter - v
			jumpSavings[k] = distSaved
		}
	}

	counts := 0
	for _, v := range jumpSavings {
		if v >= 100 {
			counts++
		}
	}

	over50 := map[int]int{}
	for _, v := range jumpSavings {
		if v >= 50 {
			over50[v]++
		}
	}

	fmt.Println(counts)
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

			for _, q := range queue {
				if q == adj {
					inQueue = true
				}
			}
			if !inQueue {
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

func OOB(x, y int, grid [][]*Point) bool {
	if y < 0 || x < 0 {
		return true
	}

	if y > len(grid)-1 || x > len(grid[0])-1 {
		return true
	}

	return false
}

func CheatJumps(origin *Point, grid [][]*Point) map[CheatJump]int {
	// a weird way to get all the points within 20 spaces
	jumps := map[CheatJump]int{}

	distance := 20
	xdelta := distance

	for ydelta := 0; ydelta <= distance; ydelta++ {
		howFar := distance

		for x := origin.X - xdelta; x <= origin.X+xdelta; x++ {
			otherX := x
			otherY := origin.Y + ydelta

			currentHowFar := howFar

			if otherX >= origin.X {
				howFar++
			} else {
				howFar--
			}

			if OOB(otherX, otherY, grid) {
				continue
			}

			otherP := grid[otherY][otherX]

			if otherP == origin || otherP.Val == "#" {
				continue
			}

			jumps[CheatJump{Start: origin, End: otherP}] = currentHowFar
		}
		xdelta--
	}
	xdelta = distance
	for ydelta := 0; ydelta >= -distance; ydelta-- {
		howFar := distance

		for x := origin.X - xdelta; x <= origin.X+xdelta; x++ {
			otherX := x
			otherY := origin.Y + ydelta

			currentHowFar := howFar

			if otherX >= origin.X {
				howFar++
			} else {
				howFar--
			}

			if OOB(otherX, otherY, grid) {
				continue
			}

			otherP := grid[otherY][otherX]

			if otherP == origin || otherP.Val == "#" {
				continue
			}
			jumps[CheatJump{Start: origin, End: otherP}] = currentHowFar
		}
		xdelta--
	}

	return jumps
}

func Around(p *Point, grid [][]*Point) []*Point {
	result := []*Point{}
	directions := []Point{
		Point{0, -1, ""},
		Point{-1, 0, ""},
		Point{1, 0, ""},
		Point{0, 1, ""},
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

		if op.Val == "#" {
			continue
		}

		result = append(result, op)
	}

	return result
}
