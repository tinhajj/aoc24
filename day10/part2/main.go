package main

import (
	"aoc24/scan"
	st "aoc24/structure"
	"fmt"
	"os"
	"strings"
)

func main() {
	b, err := os.ReadFile("day_10_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")

	digiMatrix, height, width := scan.DigitMatrix(lines, func(s string) int {
		if s == "." {
			return 999
		}
		panic("unexpected")
	})
	vertMatrix := st.VertexMatrixInt(digiMatrix)
	adjMap := map[*st.VertexInt][]*st.VertexInt{}

	for i, vertRow := range vertMatrix {
		for j, vert := range vertRow {
			directions := []st.Point{
				{Y: -1, X: 0},
				{Y: 0, X: 1},
				{Y: 1, X: 0},
				{Y: 0, X: -1},
			}
			adjacents := []*st.VertexInt{}
			for _, direction := range directions {
				ii := i + direction.Y
				jj := j + direction.X

				if ii >= 0 && ii < height && jj >= 0 && jj < width {
					adjacents = append(adjacents, vertMatrix[ii][jj])
				}
			}

			for _, adjacent := range adjacents {
				if adjacent.Val-vert.Val == 1 {
					adjMap[vert] = append(adjMap[vert], adjacent)
				}
			}
		}
	}

	// debug
	for k, adjacents := range adjMap {
		fmt.Printf("Vertex: (Y: %d, X: %d) Val: %d\n", k.Point.Y, k.Point.X, k.Val)
		fmt.Print("\tAdjacents: ")

		for _, adjacent := range adjacents {
			fmt.Printf("{(Y: %d, X: %d) Val: %d}, ", adjacent.Point.Y, adjacent.Point.X, adjacent.Val)
		}
		fmt.Println()
	}

	sum := 0
	for k := range adjMap {
		if k.Val != 0 {
			continue
		}
		sum += bfs(adjMap, k)
	}
	fmt.Println(sum)
}

func bfs(adjMap map[*st.VertexInt][]*st.VertexInt, start *st.VertexInt) int {
	// fmt.Printf("Start: Vertex: (Y: %d, X: %d) Val: %d\n", start.Point.Y, start.Point.X, start.Val)

	var queue []*st.VertexInt
	visited := make(map[*st.VertexInt]bool)

	visited[start] = true
	queue = append(queue, start)

	for len(queue) > 0 {
		if len(queue) > 1 {
			sum := 0
			for _, v := range queue {
				sum += bfs(adjMap, v)
			}
			return sum
		}

		next := queue[0]
		queue = queue[1:]

		if next.Val == 9 {
			return 1
		}

		for _, neighbour := range adjMap[next] {
			if ok := visited[neighbour]; !ok {
				visited[neighbour] = true
				queue = append(queue, neighbour)
			}
		}
	}
	return 0
}
