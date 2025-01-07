package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"maps"
	"math"
	"os"
	"strings"
)

var (
	debug  = false
	sample = false
	grid   = [][]*Cell{}

	debugImageFile *os.File
	bfsCallCounter = 0
	skipCounter    = 0
)

type Kind int

const (
	Wall Kind = iota
	Empty
)

type Vec2 struct {
	X, Y int
}

type Cell struct {
	Vec2
	Value string
	Kind  Kind
}

type Trajectory struct {
	Location  Vec2
	Direction Vec2
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
		b, err = os.ReadFile("day_16_sample_input.txt")
	} else {
		b, err = os.ReadFile("day_16_input.txt")
	}
	if err != nil {
		panic(err)
	}

	debugImageFile, err = os.Create("debug.png")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[:len(lines)-1]

	grid = [][]*Cell{}
	adj := map[*Cell][]*Cell{}

	for i, line := range lines {
		row := make([]*Cell, len(line))
		for j, r := range line {
			k := Empty
			if r == '#' {
				k = Wall
			}

			row[j] = &Cell{
				Value: string(r),
				Kind:  k,
				Vec2:  Vec2{j, i},
			}
		}
		grid = append(grid, row)
	}

	for _, row := range grid {
		for _, v := range row {
			if v.Kind == Wall {
				continue
			}
			directions := []Vec2{
				{0, -1},
				{1, 0},
				{0, 1},
				{-1, 0},
			}
			for _, direction := range directions {
				n := v.Vec2.Add(direction)
				neighbor := grid[n.Y][n.X]

				if neighbor.Kind == Empty {
					adj[v] = append(adj[v], neighbor)
				}
			}
		}
	}

	if debug {
		fmt.Println("Starting Grid")
		for _, row := range grid {
			for _, v := range row {
				fmt.Printf("%s ", v.Value)
			}
			fmt.Println()
		}
	}

	var start *Cell
	for _, row := range grid {
		for _, v := range row {
			if v.Value == "S" {
				start = v
			}
		}
	}

	costs, _ := bfs(start, Vec2{1, 0}, 0, map[Trajectory]int{}, map[*Cell]bool{start: true}, map[*Cell]bool{}, adj)
	min := math.MaxInt
	for _, c := range costs {
		if c < min {
			min = c
		}
	}
	fmt.Println(min)
}

func bfs(start *Cell, direction Vec2, currentCost int, lowestCosts map[Trajectory]int, visited map[*Cell]bool, deadEnds map[*Cell]bool, adj map[*Cell][]*Cell) (costs []int, deadEnd bool) {
	bfsCallCounter++

	if bfsCallCounter%10_000 == 0 && debug {
		fmt.Println("Skips: ", skipCounter)

		debugImageFile.Truncate(0)
		debugImageFile.Seek(0, 0)

		width, height := len(grid), len(grid[0])
		img := image.NewRGBA(image.Rect(0, 0, width, height))

		blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
		black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
		yellow := color.RGBA{R: 255, G: 0, B: 255, A: 255}

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				img.Set(x, y, blue)
			}
		}

		for y, row := range grid {
			for x, cell := range row {
				if cell == start {
					img.Set(x, y, blue)
				} else if cell.Value == "E" {
					img.Set(x, y, yellow)
				} else if cell.Kind == Wall {
					img.Set(x, y, black)
				} else if visited[cell] {
					img.Set(x, y, red)
				} else {
					img.Set(x, y, green)
				}
			}
		}
		err := png.Encode(debugImageFile, img)
		if err != nil {
			panic(err)
		}
		debugImageFile.Sync()
	}

	queue := []*Cell{start}

	var current *Cell
	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]

		if current.Value == "E" {
			return []int{currentCost}, false
		}
		visited[current] = true

		n := adj[current]
		neighbors := []*Cell{}

		for _, n := range n {
			_, alreadyVisited := visited[n]
			_, deadEnd := deadEnds[n]
			if !alreadyVisited && !deadEnd {
				neighbors = append(neighbors, n)
			}
		}

		if len(neighbors) > 1 {
			costs := []int{}
			for _, neighbor := range neighbors {
				new := map[*Cell]bool{}
				maps.Copy(new, visited)

				newDirection, expense := calculate(direction, current, neighbor)
				newCost := currentCost + expense

				// Initialize cost
				if _, ok := lowestCosts[Trajectory{Location: neighbor.Vec2, Direction: newDirection}]; !ok {
					lowestCosts[Trajectory{Location: neighbor.Vec2, Direction: newDirection}] = math.MaxInt
				}

				if newCost < lowestCosts[Trajectory{Location: neighbor.Vec2, Direction: newDirection}] {
					lowestCosts[Trajectory{Location: current.Vec2, Direction: newDirection}] = newCost
				} else {
					skipCounter++
					continue
				}

				cs, deadEnd := bfs(neighbor, newDirection, currentCost+expense, lowestCosts, new, deadEnds, adj)
				costs = append(costs, cs...)

				if deadEnd {
					deadEnds[neighbor] = true
				}
			}
			return costs, false
		} else {
			queue = append(queue, neighbors...)
			if len(neighbors) > 0 {
				newDirection, expense := calculate(direction, current, neighbors[0])
				currentCost = expense + currentCost

				// Initialize cost
				if _, ok := lowestCosts[Trajectory{Location: current.Vec2, Direction: newDirection}]; !ok {
					lowestCosts[Trajectory{Location: current.Vec2, Direction: newDirection}] = math.MaxInt
				}

				if currentCost < lowestCosts[Trajectory{Location: current.Vec2, Direction: newDirection}] {
					lowestCosts[Trajectory{Location: current.Vec2, Direction: newDirection}] = currentCost
				} else {
					skipCounter++
					return []int{}, false // give up
				}

				direction = newDirection
			}
		}
	}
	var lonely bool
	if len(adj[current]) < 2 {
		lonely = true
	}
	return []int{}, lonely
}

func calculate(direction Vec2, current *Cell, next *Cell) (newDirection Vec2, cost int) {
	if current.Vec2.Add(direction) == next.Vec2 {
		return direction, 1
	}

	if current.Vec2.X != next.X {
		if current.Vec2.X > next.X {
			return Vec2{-1, 0}, 1001
		} else {
			return Vec2{1, 0}, 1001
		}
	}

	if current.Vec2.Y != next.Y {
		if current.Vec2.Y > next.Y {
			return Vec2{0, -1}, 1001
		} else {
			return Vec2{0, 1}, 1001
		}
	}
	panic(-1)
}
