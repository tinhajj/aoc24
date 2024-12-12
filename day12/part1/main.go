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
	m, _, _ := scan.RuneMatrix(lines)
	matrix := st.VertexMatrixStr(m)

	for _, row := range matrix {
		for _, vert := range row {
			fmt.Println(vert)
		}
	}
}
