package main

import (
	"aoc24/scan"
	st "aoc24/structure"
	"fmt"
	"os"
	"strings"
)

type dimension struct {
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

			found, area, perimeter := bfs(height, width, adjMap, v)
			visited = append(visited, found...)

			regions = append(regions, dimension{
				plant:     letter,
				perimeter: perimeter,
				area:      area,
			})
		}
	}

	sum := 0
	for _, v := range regions {
		if debug {
			fmt.Printf("%s area %d * perimeter %d\n", v.plant, v.area, v.perimeter)
		}
		sum += v.area * v.perimeter
	}
	fmt.Println(sum)
}

func bfs(matrixH, matrixW int, adjMap map[*st.VertexStr][]*st.VertexStr, start *st.VertexStr) (found []*st.VertexStr, a, p int) {
	a = 0
	p = 0

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
		}
		if next.Point.Y == 0 {
			p++
		}
		if next.Point.X == matrixW-1 {
			p++
		}
		if next.Point.Y == matrixH-1 {
			p++
		}

		for _, neighbour := range adjMap[next] {
			if ok := visited[neighbour]; !ok {
				if neighbour.Val != start.Val {
					p++
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

	return found, a, p
}
