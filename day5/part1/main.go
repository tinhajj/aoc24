package main

import (
	"aoc24/scan"
	"fmt"
	"os"
	"strings"
)

func main() {
	b, err := os.ReadFile("day_5_input.txt")
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[0 : len(lines)-1]

	rules := map[int][]int{}

	i := 0
	for ; i < len(lines); i++ {
		numbers := scan.Numbers(lines[i])
		if len(numbers) < 1 {
			i++
			break
		}
		first, second := numbers[0], numbers[1]
		rules[first] = append(rules[first], second)
	}

	updates := [][]int{}
	for ; i < len(lines); i++ {
		numbers := scan.Numbers(lines[i])
		updates = append(updates, numbers)
	}

	sum := 0
update:
	for _, update := range updates {
		for j, page := range update {
			previous := update[:j]
			pageRules := rules[page]

			for _, otherPage := range previous {
				for _, pageRule := range pageRules {
					if otherPage == pageRule {
						continue update
					}
				}
			}
		}
		sum += update[len(update)/2]
	}
	fmt.Println(sum)
}
