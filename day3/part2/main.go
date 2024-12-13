package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	b, err := os.ReadFile("day_3_input.txt")
	if err != nil {
		panic(err)
	}
	input := string(b)
	lines := strings.Split(input, "\n")

	enabled := true
	sum := 0
	for _, line := range lines {
		memory := line
		re := regexp.MustCompile(`(mul\(\d{1,3},\d{1,3}\)|don't\(\)|do\(\))`)

		matches := re.FindAllStringSubmatch(memory, -1)

		for _, match := range matches {
			instruction := match[0]

			if instruction == "do()" {
				enabled = true
				continue
			}
			if instruction == "don't()" {
				enabled = false
				continue
			}

			opPart := instruction[4 : len(instruction)-1]
			opParts := strings.Split(opPart, ",")

			op1s, op2s := opParts[0], opParts[1]
			op1, err := strconv.Atoi(op1s)
			if err != nil {
				panic(err)
			}
			op2, err := strconv.Atoi(op2s)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%d * %d\n", op1, op2)
			if enabled {
				sum += (op1 * op2)
			}
		}
	}
	fmt.Println(sum)
}
