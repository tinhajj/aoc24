package main

import (
	"aoc24/scan"
	"fmt"
	"os"
	"strings"
)

type point struct {
	x int
	y int
}

type button struct {
	x    int
	y    int
	cost int
}

type machine struct {
	buttonA  button
	buttonB  button
	prizeLoc point
}

var debug = false

func main() {
	b, err := os.ReadFile("day_13_sample_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[:len(lines)-1]

	machines := []machine{}

	for i := 0; i < len(lines); i += 4 {
		var parts []int

		machine := machine{}

		parts = scan.Numbers(lines[i])
		buttonA := button{
			x:    parts[0],
			y:    parts[1],
			cost: 3,
		}
		machine.buttonA = buttonA

		parts = scan.Numbers(lines[i+1])
		buttonB := button{
			x:    parts[0],
			y:    parts[1],
			cost: 1,
		}
		machine.buttonB = buttonB

		parts = scan.Numbers(lines[i+2])
		machine.prizeLoc = point{
			x: parts[0],
			y: parts[1],
		}
		machines = append(machines, machine)
	}

	sum := 0

	for _, machine := range machines {
		if debug {
			fmt.Printf("%+v\n", machine)
		}
		costs := []int{}
		for i := 1; i <= 100; i++ {
			for j := 1; j <= 100; j++ {
				buttonA := machine.buttonA
				buttonB := machine.buttonB

				x := (buttonA.x * i) + (buttonB.x * j)
				y := (buttonA.y * i) + (buttonB.y * j)
				cost := (buttonA.cost * i) + (buttonB.cost * j)

				if x == machine.prizeLoc.x && y == machine.prizeLoc.y {
					costs = append(costs, cost)
				}
			}
		}

		if len(costs) < 1 {
			continue
		}
		if len(costs) == 1 {
			sum += costs[0]
		}
		var least int
		for _, c := range costs {
			if c < least || c == 0 {
				least = c
			}
		}
		sum += least
	}

	fmt.Println(sum)
}
