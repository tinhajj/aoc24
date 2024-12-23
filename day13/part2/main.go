/*
A linear algebra problem.  The problem describes a cost
between two buttons which can be pressed to work towards a location on the x
and y axis.  The buttons have a different cost associated with them which on
the surface makes it seem like the problem is about maximizing cost to get
towards the location.  But in linear algebra when you're working with a
system of equations there can either be zero, one or infinite solutions to
the system.  So there isn't really a concept of an optimal solution, there is
just one solution or none in this case.

Take the example:

Button A: X+94, Y+34 (Cost 3)
Button B: X+22, Y+67 (Cost 1)
Prize: X=8400, Y=5400

This translates into a system of equations that looks like this:

94x + 22y = 8400
34x + 67y = 5400

It's possible to use substitution to solve one of these equations for x
For instance the first equation could be rewritten as:

y = (8400 - 94x) / 22

This equation can be used to substitute into the second equation:

34x + [(8400 - 94x) / 22] = 5400

And this equation, with just one variable, can be solved for a constant.  Then
that constant can be plugged into any of the original equations to solve for the
other variable's value.

This whole operation can be 'formulaized' and be done for each machine.

Still, after all of this it does feel strange that there can only be one
solution.  Just based on the problem description it sounds like it should be
possible to have different solutions involving different amounts of presses, but
I guess not.  Visualize the problem as two lines?
*/
package main

import (
	"aoc24/scan"
	"fmt"
	"math"
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
	b, err := os.ReadFile("day_13_input.txt")
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
			x: parts[0] + 10000000000000,
			y: parts[1] + 10000000000000,
		}
		machines = append(machines, machine)
	}

	sum := 0.0
	for _, machine := range machines {
		if debug {
			fmt.Printf("%+v\n", machine)
		}

		lhsI := (machine.prizeLoc.x * machine.buttonB.y) - (machine.prizeLoc.y * machine.buttonB.x)
		rhsI := -((machine.buttonA.y) * machine.buttonB.x) + (machine.buttonA.x * machine.buttonB.y)

		lhsF := float64(lhsI)
		rhsF := float64(rhsI)

		quotientA := lhsF / rhsF

		if quotientA < 0 {
			continue
		}

		if math.Trunc(quotientA) != quotientA {
			continue
		}

		rhsF = float64(machine.prizeLoc.y) - float64(machine.buttonA.y)*quotientA

		quotientB := rhsF / float64(machine.buttonB.y)

		if quotientB < 0 {
			continue
		}

		if math.Trunc(quotientB) != quotientB {
			continue
		}

		ans := (quotientA * 3) + (quotientB * 1)

		sum += ans

		if debug {
			fmt.Println(
				lhsF,
				rhsF,
				quotientA,
				quotientB,
			)
		}
	}
	fmt.Printf("%d\n", int(sum))
}
