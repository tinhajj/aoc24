package main

import (
	"aoc24/scan"
	st "aoc24/structure"
	"fmt"
	"os"
	"strings"
)

type direction int

const (
	UP    direction = iota
	DOWN  direction = iota
	LEFT  direction = iota
	RIGHT direction = iota
)

type dimension struct {
	bumps     map[direction]map[*st.VertexStr]bool
	vertices  []*st.VertexStr
	plant     string
	perimeter int
	area      int
}

var debug = false

func main() {
	b, err := os.ReadFile("day_12_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[:len(lines)-1]
	m, height, width := scan.RuneMatrix(lines)
	matrix := st.VertexMatrixStr(m)

	adjMap := map[*st.VertexStr][]*st.VertexStr{}

	for i, row := range matrix {
		for j, vert := range row {
			directions := []st.Point{
				{Y: -1, X: 0},
				{Y: 0, X: 1},
				{Y: 1, X: 0},
				{Y: 0, X: -1},
			}

			adjacents := []*st.VertexStr{}
			for _, direction := range directions {
				ii := i + direction.Y
				jj := j + direction.X

				if ii >= 0 && ii < height && jj >= 0 && jj < width {
					adjacents = append(adjacents, matrix[ii][jj])
				}
			}

			for _, adjacent := range adjacents {
				adjMap[vert] = append(adjMap[vert], adjacent)
			}
		}
	}

	// debug
	for k, v := range adjMap {
		if !debug {
			continue
		}
		fmt.Printf("Vertex: (%+v) %+v\n", k.Point, k.Val)

		fmt.Printf("Adjacents: ")
		for _, other := range v {
			fmt.Printf("%+v ", other.Val)
		}
		fmt.Println("")
		fmt.Println("")
	}

	visited := []*st.VertexStr{}
	regions := []dimension{}

	for _, row := range matrix {
		for _, v := range row {
			seen := false
			for _, visit := range visited {
				if visit == v {
					seen = true
					break
				}
			}
			if seen {
				continue
			}

			letter := v.Val

			bumps, found, area, perimeter := bfs(height, width, adjMap, v)
			visited = append(visited, found...)

			regions = append(regions, dimension{
				bumps:     bumps,
				vertices:  found,
				plant:     letter,
				perimeter: perimeter,
				area:      area,
			})
		}
	}

	sum := 0
	for _, dim := range regions {
		sides := sides(dim, adjMap)
		sum += dim.area * sides

		if debug {
			fmt.Printf("%s area %d * perimeter %d sides: %d\n", dim.plant, dim.area, dim.perimeter, sides)
			fmt.Print("vertices: ")
			for _, v := range dim.vertices {
				fmt.Printf("%s (y: %d, x: %d), ", v.Val, v.Point.Y, v.Point.X)
			}
			fmt.Println()
		}
	}

	fmt.Println(sum)
}

func bfs(matrixH, matrixW int, adjMap map[*st.VertexStr][]*st.VertexStr, start *st.VertexStr) (bumps map[direction]map[*st.VertexStr]bool, found []*st.VertexStr, a, p int) {
	a = 0
	p = 0
	bumps = map[direction]map[*st.VertexStr]bool{UP: {}, DOWN: {}, LEFT: {}, RIGHT: {}}

	var queue []*st.VertexStr
	visited := make(map[*st.VertexStr]bool)

	visited[start] = true
	queue = append(queue, start)

	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]

		if next.Val == start.Val {
			a++
		}
		if next.Point.X == 0 {
			p++
			bumps[LEFT][next] = false
		}
		if next.Point.Y == 0 {
			p++
			bumps[UP][next] = false
		}
		if next.Point.X == matrixW-1 {
			p++
			bumps[RIGHT][next] = false
		}
		if next.Point.Y == matrixH-1 {
			p++
			bumps[DOWN][next] = false
		}

		for _, neighbour := range adjMap[next] {
			if ok := visited[neighbour]; !ok {
				if neighbour.Val != start.Val {
					p++

					if neighbour.Point.X > next.Point.X {
						bumps[RIGHT][next] = false
					}
					if neighbour.Point.X < next.Point.X {
						bumps[LEFT][next] = false
					}
					if neighbour.Point.Y > next.Point.Y {
						bumps[DOWN][next] = false
					}
					if neighbour.Point.Y < next.Point.Y {
						bumps[UP][next] = false
					}

					continue
				}

				visited[neighbour] = true
				queue = append(queue, neighbour)
			}
		}
	}

	found = []*st.VertexStr{}
	for k, _ := range visited {
		if k.Val == start.Val {
			found = append(found, k)
		}
	}

	return bumps, found, a, p
}

func sides(dim dimension, adjMap map[*st.VertexStr][]*st.VertexStr) int {
	//if ok := visited[neighbour]; !ok && in(neighbour, bumps) {
	bumps := dim.bumps
	directions := []direction{UP, DOWN, LEFT, RIGHT}

	sides := 0

	unvisited := func(d direction) *st.VertexStr {
		bumps := bumps[d]
		for k, v := range bumps {
			if !v {
				return k
			}
		}
		return nil
	}

	for _, d := range directions {
		for start := unvisited(d); start != nil; start = unvisited(d) {
			var queue []*st.VertexStr
			visited := make(map[*st.VertexStr]bool)

			visited[start] = true
			queue = append(queue, start)

			for len(queue) > 0 {
				next := queue[0]
				queue = queue[1:]

				for _, neighbour := range adjMap[next] {
					if ok := visited[neighbour]; !ok && in(neighbour, bumps[d]) {
						visited[neighbour] = true
						queue = append(queue, neighbour)
					}
				}
			}

			for k, _ := range visited {
				bumps[d][k] = true
			}
			sides++
		}
	}
	return sides
}

func in(needle *st.VertexStr, haystack map[*st.VertexStr]bool) bool {
	v, ok := haystack[needle]
	if ok && !v {
		return true
	}
	return false
}
