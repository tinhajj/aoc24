package main

import (
	"fmt"
	"maps"
	"math"
	"os"
	"strings"
)

var (
	debug  = true
	sample = false
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

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[:len(lines)-1]

	grid := [][]*Cell{}
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
		fmt.Println("grid")
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

	costs := bfs(start, Vec2{1, 0}, 0, math.MaxInt, map[*Cell]bool{start: true}, adj)
	min := math.MaxInt
	for _, c := range costs {
		if c < min {
			min = c
		}
	}
	fmt.Println(min)
}

func bfs(start *Cell, direction Vec2, currentCost int, lowestCost int, visited map[*Cell]bool, adj map[*Cell][]*Cell) []int {
	queue := []*Cell{start}

	var current *Cell
	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]

		if current.Value == "E" {
			return []int{currentCost}
		}
		visited[current] = true

		n := adj[current]
		neighbors := []*Cell{}

		for _, n := range n {
			_, ok := visited[n]
			if !ok {
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

				if newCost > lowestCost {
					continue // give up
				}

				costs = append(costs, bfs(neighbor, newDirection, currentCost+expense, lowestCost, new, adj)...)
				for _, c := range costs {
					if c < lowestCost {
						lowestCost = c
					}
				}
			}
			return costs
		} else {
			queue = append(queue, neighbors...)
			if len(neighbors) > 0 {
				newDirection, expense := calculate(direction, current, neighbors[0])

				currentCost = expense + currentCost
				if currentCost > lowestCost {
					return []int{} // give up
				}
				direction = newDirection
			}
		}
	}
	return []int{}
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
