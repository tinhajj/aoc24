package main

import (
	"aoc24/scan"
	"fmt"
	"os"
	"slices"
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
	rules, updates := getBadUpdates(lines)

	sum := 0

	for _, update := range updates {
		correctUpdate := make([]int, len(update))

		for j, pageNumber := range update {
			others := slices.Concat(update[:j], update[j+1:])

			pageRules := rules[pageNumber]

			matches := 0
			for _, pageRule := range pageRules {
				for _, other := range others {
					if other == pageRule {
						matches++
					}
				}
			}

			correctUpdate[len(update)-1-matches] = pageNumber
		}

		sum += correctUpdate[len(correctUpdate)/2]
	}
	fmt.Println(sum)
}

func getBadUpdates(lines []string) (map[int][]int, [][]int) {
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

	badUpdates := [][]int{}

update:
	for _, update := range updates {
		for j, page := range update {
			previous := update[:j]
			pageRules := rules[page]

			for _, otherPage := range previous {
				for _, pageRule := range pageRules {
					if otherPage == pageRule {
						badUpdates = append(badUpdates, update)
						continue update
					}
				}
			}
		}
	}

	return rules, badUpdates
}
