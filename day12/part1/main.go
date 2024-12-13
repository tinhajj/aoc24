package main

import (
	"aoc24/scan"
	st "aoc24/structure"
	"fmt"
	"os"
	"strings"
)

func main() {
	b, err := os.ReadFile("day_12_sample_input.txt")
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

	for k, v := range adjMap {
		fmt.Printf("Vertex: %+v\n", k.Point)

		fmt.Printf("Adjacents: ")
		for _, other := range v {
			fmt.Printf("%+v ", other.Val)
		}
		fmt.Println("")
		fmt.Println("")
	}
}
