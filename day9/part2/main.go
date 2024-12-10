package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Block struct {
	Val  int
	Size int
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

	diskmap := []Block{}
	backwards := []Block{}

	BlockString := func(b []Block) string {
		s := ""
		for _, block := range b {
			if block.Size < 1 {
				continue
			}
			for i := 0; i < block.Size; i++ {
				if block.Val == -1 {
					s += fmt.Sprintf(".")
				} else {
					s += fmt.Sprint(block.Val)
				}
			}
		}
		return s
	}
	_ = BlockString

	id := 0
	for i := 0; i < len(line); i++ {
		num, err := strconv.Atoi(string(line[i]))
		if err != nil {
			panic(err)
		}

		if i%2 == 0 {
			size := num
			diskmap = append(diskmap, Block{Val: id, Size: size})
			id++
		} else {
			freeblocks := num

			if freeblocks == 0 {
				continue
			}

			diskmap = append(diskmap, Block{Val: -1, Size: freeblocks})
		}
	}

	for i := len(diskmap) - 1; i > 0; i-- {
		if diskmap[i].Val != -1 {
			og := diskmap[i]
			backwards = append(backwards, og)
		}
	}

	for _, rightFile := range backwards {
		filePosition := 0
		for i, leftFile := range diskmap {
			if rightFile.Val == leftFile.Val {
				filePosition = i
			}
		}

		for j, leftBlock := range diskmap {
			if leftBlock.Val == -1 && leftBlock.Size >= rightFile.Size && j < filePosition {
				diskmap[j].Size = leftBlock.Size - rightFile.Size
				diskmap[filePosition] = Block{Val: -1, Size: rightFile.Size}
				diskmap = slices.Concat(diskmap[:j], []Block{rightFile}, []Block{diskmap[j]}, diskmap[j+1:])

				break
			}
		}
	}

	idx := 0
	sum := 0
	for _, block := range diskmap {
		if block.Val == -1 {
			idx += block.Size
			continue
		}
		for i := 0; i < block.Size; i++ {
			sum += block.Val * idx
			idx++
		}
	}
	fmt.Println(sum)
}
