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

var StoneNil Stone = Stone{Value: -1, Left: 0}

func main() {
	b, err := os.ReadFile("day_11_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")

	nums := scan.Numbers(lines[0])
	stones := []Stone{}

	for i := len(nums) - 1; i >= 0; i-- {
		stones = append(stones, Stone{nums[i], 75})
	}

	memo := map[int]bool{}

	sum := 0
	var currentStone Stone
	for len(stones) > 0 {
		currentStone, stones = stones[len(stones)-1], stones[:len(stones)-1]

		if currentStone == StoneNil {
			continue
		}

		memo[currentStone.Value] = true

		if currentStone.Left == 0 {
			sum++
			continue
		}

		s1, s2 := rule(currentStone)
		stones = append(stones, s2)
		stones = append(stones, s1)
	}

	fmt.Println(sum)
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
