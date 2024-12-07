package main

import (
	"aoc24/scan"
	"fmt"
	"os"
	"strings"
)

type Operation int

const (
	ADD      Operation = iota
	MULTIPLY Operation = iota
)

func main() {
	b, err := os.ReadFile("day_7_sample_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]

	for _, line := range lines {
		numbers := scan.Numbers(line)
		result := numbers[0]
		operands := numbers[1:]

		totalOperations := len(operands) - 1

		for i := 0; i < totalOperations; i++ {
		}

		fmt.Println(result, operands)
	}
}
