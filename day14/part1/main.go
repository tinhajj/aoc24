package main

import (
	"aoc24/scan"
	"fmt"
	"os"
	"strings"
)

var (
	WIDTH, HEIGHT = 101, 103
	SECONDS       = 100

	debug  = true
	sample = false
)

type Vec2 struct {
	X, Y int
}

type Robot struct {
	Loc, Vel Vec2
}

type Quadrant struct {
	Start, Stop Vec2
	RobotCount  int
}

func Abs(i int) int {
	if i == 0 {
		return i
	}
	if i > 0 {
		return i
	}
	return i * -1
}

func RobotMove(r Robot) Robot {
	r.Loc.X = r.Loc.X + r.Vel.X
	r.Loc.Y = r.Loc.Y + r.Vel.Y

	if r.Loc.X >= WIDTH {
		r.Loc.X = r.Loc.X % WIDTH
	}

	if r.Loc.X < 0 {
		r.Loc.X = (r.Loc.X % -WIDTH) + WIDTH
	}

	if r.Loc.Y >= HEIGHT {
		r.Loc.Y = r.Loc.Y % HEIGHT
	}

	if r.Loc.Y < 0 {
		r.Loc.Y = (r.Loc.Y % -HEIGHT) + HEIGHT
	}

	return r
}

func RobotInQuadrant(r Robot, q Quadrant) bool {
	XGood := false
	YGood := false

	if r.Loc.X >= q.Start.X && r.Loc.X <= q.Stop.X {
		XGood = true
	}

	if r.Loc.Y >= q.Start.Y && r.Loc.Y <= q.Stop.Y {
		YGood = true
	}

	return XGood && YGood
}

func main() {
	// NOTE: Might not need a matrix, it would just be a good data structure to print out a visual?
	ROOM := make([][]int, HEIGHT)
	for i := range ROOM {
		ROOM[i] = make([]int, WIDTH)
	}

	var b []byte
	var err error

	if sample {
		b, err = os.ReadFile("day_14_sample_input.txt")
		WIDTH, HEIGHT = 11, 7
	} else {
		b, err = os.ReadFile("day_14_input.txt")
	}
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[:len(lines)-1]

	ROBOTS := []Robot{}

	for _, line := range lines {
		nums := scan.Numbers(line)
		loc := Vec2{Y: nums[1], X: nums[0]}
		vel := Vec2{Y: nums[3], X: nums[2]}
		ROBOTS = append(ROBOTS, Robot{
			Loc: loc,
			Vel: vel,
		})
	}

	for i := 0; i < SECONDS; i++ {
		for i, robot := range ROBOTS {
			ROBOTS[i] = RobotMove(robot)
		}
	}

	quadrants := []Quadrant{
		{
			Start:      Vec2{0, 0},
			Stop:       Vec2{(WIDTH / 2) - 1, (HEIGHT / 2) - 1},
			RobotCount: 0,
		},
		{
			Start:      Vec2{(WIDTH / 2) + 1, 0},
			Stop:       Vec2{WIDTH - 1, (HEIGHT / 2) - 1},
			RobotCount: 0,
		},
		{
			Start:      Vec2{0, (HEIGHT / 2) + 1},
			Stop:       Vec2{(WIDTH / 2) - 1, HEIGHT - 1},
			RobotCount: 0,
		},
		{
			Start:      Vec2{(WIDTH / 2) + 1, (HEIGHT / 2) + 1},
			Stop:       Vec2{WIDTH - 1, HEIGHT - 1},
			RobotCount: 0,
		},
	}

	for _, r := range ROBOTS {
		for j, q := range quadrants {
			if RobotInQuadrant(r, q) {
				quadrants[j].RobotCount++
			}
		}
	}

	for _, q := range quadrants {
		if debug {
			fmt.Printf("%+v\n", q)
		}
	}

	sum := 0
	for _, q := range quadrants {
		if sum == 0 {
			sum = q.RobotCount
			continue
		}
		sum = sum * q.RobotCount
	}
	fmt.Println(sum)
}
