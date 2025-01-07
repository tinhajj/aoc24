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
	DEBUG  = false
	SAMPLE = false
	GRID   = [][]*Cell{}

	VALIDPATHS = []Path{}
	GOODCELLS  = map[*Cell]bool{}

	DEBUGIMAGEFILE *os.File
	BfsCallCounter = 0
	SkipCounter    = 0
)

type Kind int

const (
	Wall Kind = iota
	Empty
)

type Path struct {
	Visited map[*Cell]bool
	Cost    int
}

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

	if SAMPLE {
		b, err = os.ReadFile("day_16_sample_input.txt")
	} else {
		b, err = os.ReadFile("day_16_input.txt")
	}
	if err != nil {
		panic(err)
	}

	if DEBUG {
		DEBUGIMAGEFILE, err = os.Create("debug.png")
		if err != nil {
			panic(err)
		}
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[:len(lines)-1]

	GRID = [][]*Cell{}
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
		GRID = append(GRID, row)
	}

	for _, row := range GRID {
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
				neighbor := GRID[n.Y][n.X]

				if neighbor.Kind == Empty {
					adj[v] = append(adj[v], neighbor)
				}
			}
		}
	}

	if DEBUG {
		fmt.Println("Starting Grid")
		for _, row := range GRID {
			for _, v := range row {
				fmt.Printf("%s ", v.Value)
			}
			fmt.Println()
		}
	}

	var start *Cell
	for _, row := range GRID {
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

	for _, path := range VALIDPATHS {
		if path.Cost == min {
			for k, v := range path.Visited {
				if v {
					GOODCELLS[k] = true
				}
			}
		}

		if DEBUG {
			fmt.Println("Valid Path:")
			for _, row := range GRID {
				for _, c := range row {
					if path.Visited[c] {
						fmt.Printf("O")
					} else if c.Kind == Wall {
						fmt.Printf("#")
					} else {
						fmt.Printf(".")
					}
				}
				fmt.Println()
			}
		}
	}

	numberOfGoodCells := 0
	for _, v := range GOODCELLS {
		if v {
			numberOfGoodCells++
		}
	}
	fmt.Println(numberOfGoodCells)
}

func bfs(start *Cell, direction Vec2, currentCost int, lowestCosts map[Trajectory]int, visited map[*Cell]bool, deadEnds map[*Cell]bool, adj map[*Cell][]*Cell) (costs []int, deadEnd bool) {
	BfsCallCounter++

	if DEBUG && BfsCallCounter%10_000 == 0 {
		fmt.Println("Skips: ", SkipCounter)

		DEBUGIMAGEFILE.Truncate(0)
		DEBUGIMAGEFILE.Seek(0, 0)

		width, height := len(GRID), len(GRID[0])
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

		for y, row := range GRID {
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
		err := png.Encode(DEBUGIMAGEFILE, img)
		if err != nil {
			panic(err)
		}
		DEBUGIMAGEFILE.Sync()
	}

	queue := []*Cell{start}

	var current *Cell
	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]

		if current.Value == "E" {
			visited[current] = true

			VALIDPATHS = append(VALIDPATHS, Path{
				Visited: visited,
				Cost:    currentCost,
			})

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

		if DEBUG {
			if len(neighbors) > 1 {
				for _, row := range GRID {
				row:
					for _, c := range row {
						for _, n := range neighbors {
							if c == n {
								fmt.Printf("N")
								continue row
							}
						}
						if visited[c] {
							fmt.Printf("O")
						} else if c == current {
							fmt.Printf("@")
						} else if c.Kind == Wall {
							fmt.Printf("#")
						} else {
							fmt.Printf(".")
						}
					}
					fmt.Println()
				}
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

				if newCost <= lowestCosts[Trajectory{Location: neighbor.Vec2, Direction: newDirection}] {
					lowestCosts[Trajectory{Location: current.Vec2, Direction: newDirection}] = newCost
				} else {
					SkipCounter++
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
				if _, ok := lowestCosts[Trajectory{Location: neighbors[0].Vec2, Direction: newDirection}]; !ok {
					lowestCosts[Trajectory{Location: neighbors[0].Vec2, Direction: newDirection}] = math.MaxInt
				}

				if currentCost <= lowestCosts[Trajectory{Location: neighbors[0].Vec2, Direction: newDirection}] {
					lowestCosts[Trajectory{Location: neighbors[0].Vec2, Direction: newDirection}] = currentCost
				} else {
					SkipCounter++
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
