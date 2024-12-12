package main

import (
	"aoc24/scan"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Stone struct {
	Value int
	Left  int
}

var StoneNil = Stone{Value: -1, Left: 0}
var Memo = map[Stone]int{}

func main() {
	b, err := os.ReadFile("day_11_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")

	nums := scan.Numbers(lines[0])
	stones := []Stone{}

	for _, num := range nums {
		stones = append(stones, Stone{num, 75})
	}

	sum := 0
	for _, stone := range stones {
		sum += solve(stone)
	}
	fmt.Println(sum)
}

func solve(stone Stone) int {
	if stone == StoneNil {
		return 0
	}
	if stone.Left == 0 {
		return 1
	}

	_, ok := Memo[stone]
	if ok {
		return Memo[stone]
	}

	stone1, stone2 := rule(stone)

	sum1 := solve(stone1)
	sum2 := solve(stone2)

	Memo[stone1] = sum1
	Memo[stone2] = sum2

	return sum1 + sum2
}

func rule(stone Stone) (Stone, Stone) {
	if stone.Value == 0 {
		return Stone{Value: 1, Left: stone.Left - 1}, StoneNil
	}

	plain := strconv.Itoa(stone.Value)
	if len(plain)%2 == 0 {
		mid := len(plain) / 2
		first, err := strconv.Atoi(plain[:mid])
		if err != nil {
			panic(err)
		}
		second, err := strconv.Atoi(plain[mid:])
		if err != nil {
			panic(err)
		}
		return Stone{Value: first, Left: stone.Left - 1}, Stone{Value: second, Left: stone.Left - 1}
	}

	return Stone{Value: stone.Value * 2024, Left: stone.Left - 1}, StoneNil
}
