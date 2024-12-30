package main

import (
	"aoc24/scan"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"strings"
)

var (
	WIDTH, HEIGHT = 101, 103
	IMAGELIMIT    = 100

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
}

func QuadrantImage(robots []Robot, name string) {
	topLeft := image.Point{0, 0}
	bottomRight := image.Point{WIDTH, HEIGHT}

	img := image.NewRGBA(image.Rectangle{topLeft, bottomRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{100, 200, 200, 0xff}

	for _, r := range robots {
		img.Set(r.Loc.X, r.Loc.Y, cyan)
	}

	// Encode as PNG.
	f, _ := os.Create(name)
	png.Encode(f, img)
	f.Close()
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

func SafetyFactor(quadrants []Quadrant, robots []Robot) int {
	counts := map[Quadrant]int{}

	for _, r := range robots {
		for _, q := range quadrants {
			if RobotInQuadrant(r, q) {
				counts[q]++
			}
		}
	}

	sum := 0
	for _, v := range counts {
		if sum == 0 {
			sum = v
			continue
		}
		sum = sum * v
	}
	return sum
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

	quadrants := []Quadrant{
		{
			Start: Vec2{0, 0},
			Stop:  Vec2{(WIDTH / 2) - 1, (HEIGHT / 2) - 1},
		},
		{
			Start: Vec2{(WIDTH / 2) + 1, 0},
			Stop:  Vec2{WIDTH - 1, (HEIGHT / 2) - 1},
		},
		{
			Start: Vec2{0, (HEIGHT / 2) + 1},
			Stop:  Vec2{(WIDTH / 2) - 1, HEIGHT - 1},
		},
		{
			Start: Vec2{(WIDTH / 2) + 1, (HEIGHT / 2) + 1},
			Stop:  Vec2{WIDTH - 1, HEIGHT - 1},
		},
	}

	lowestSafetyFactor := math.MaxInt
	imagesWritten := 0
	fmt.Println("Look for trees in written images")
	for i := 1; true; i++ {
		if imagesWritten > IMAGELIMIT {
			log.Fatalln("written the maximum number of images")
		}

		for i, robot := range ROBOTS {
			ROBOTS[i] = RobotMove(robot)
		}

		sf := SafetyFactor(quadrants, ROBOTS)

		if sf < lowestSafetyFactor {
			lowestSafetyFactor = sf

			fileName := fmt.Sprintf("output/second-%d.png", i)
			QuadrantImage(ROBOTS, fileName)
			fmt.Printf("Image %s written\n", fileName)
			imagesWritten++
		}
	}
}
