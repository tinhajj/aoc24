package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	Val int
	Idx int
}

func main() {
	b, err := os.ReadFile("day_9_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]
	line := lines[0]

	diskmap := []int{}
	backwards := []Point{}

	id := 0
	for i := 0; i < len(line); i++ {
		num, err := strconv.Atoi(string(line[i]))
		if err != nil {
			panic(err)
		}

		if i%2 == 0 {
			size := num
			for i := 0; i < size; i++ {
				diskmap = append(diskmap, id)
			}
			id++
		} else {
			freeblocks := num
			for i := 0; i < freeblocks; i++ {
				diskmap = append(diskmap, -1)
			}
		}
	}

	for i := len(diskmap) - 1; i > 0; i-- {
		if diskmap[i] != -1 {
			backwards = append(backwards, Point{Val: diskmap[i], Idx: i})
		}
	}

	for i := 0; i < len(diskmap); i++ {
		if diskmap[i] == -1 && len(backwards) > 0 && backwards[0].Idx > i {
			end := backwards[0]
			diskmap[i] = end.Val
			diskmap[end.Idx] = -1
			backwards = backwards[1:]
		}
	}

	sum := 0
	for i, v := range diskmap {
		if v == -1 {
			continue
		}
		sum += i * v
	}
	fmt.Println(sum)
}
