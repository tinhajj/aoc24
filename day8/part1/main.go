package main

import (
	"fmt"
	"os"
	"strings"
)

type Point struct {
	Y int
	X int
}

func main() {
	b, err := os.ReadFile("day_8_sample_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]

	antennas := map[string][]Point{}

	for y, line := range lines {
		for x, char := range line {
			if string(char) == "." {
				continue
			}
			antennas[string(char)] = append(antennas[string(char)], Point{y, x})
		}
	}
	fmt.Println(antennas)
}
