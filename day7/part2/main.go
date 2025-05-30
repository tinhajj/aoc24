package main

import (
	"aoc24/scan"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Operation int

const (
	ADD      Operation = iota
	MULTIPLY Operation = iota
	CONCAT   Operation = iota
)

func main() {
	b, err := os.ReadFile("day_7_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]

	sum := 0

	for _, line := range lines {
		numbers := scan.Numbers(line)
		goal := numbers[0]
		operands := numbers[1:]

		possibilities := DoOperation(ADD, operands, 0)

		for _, possibility := range possibilities {
			if possibility == goal {
				sum += goal
				break
			}
		}
	}
	fmt.Println(sum)
}

func DoOperation(op Operation, operands []int, sum int) []int {
	if len(operands) == 1 {
		switch op {
		case MULTIPLY:
			return []int{sum * operands[0]}
		case ADD:
			return []int{sum + operands[0]}
		case CONCAT:
			first := strconv.Itoa(sum)
			second := strconv.Itoa(operands[0])
			combined, err := strconv.Atoi(first + second)

			if err != nil {
				return []int{}
			}

			return []int{combined}
		}
	}

	switch op {
	case MULTIPLY:
		result := sum * operands[0]

		option1 := DoOperation(MULTIPLY, operands[1:], result)
		option2 := DoOperation(ADD, operands[1:], result)
		option3 := DoOperation(CONCAT, operands[1:], result)
		return slices.Concat(option1, option2, option3)
	case ADD:
		result := sum + operands[0]

		option1 := DoOperation(MULTIPLY, operands[1:], result)
		option2 := DoOperation(ADD, operands[1:], result)
		option3 := DoOperation(CONCAT, operands[1:], result)
		return slices.Concat(option1, option2, option3)
	case CONCAT:
		first := strconv.Itoa(sum)
		second := strconv.Itoa(operands[0])

		combined, err := strconv.Atoi(first + second)

		if err != nil {
			return []int{}
		}

		option1 := DoOperation(MULTIPLY, operands[1:], combined)
		option2 := DoOperation(ADD, operands[1:], combined)
		option3 := DoOperation(CONCAT, operands[1:], combined)
		return slices.Concat(option1, option2, option3)
	}

	panic("unknown op")
}
