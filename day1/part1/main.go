package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	b, err := os.ReadFile("day_1_input.txt")
	if err != nil {
		panic(err)
	}
	input := string(b)

	list1, list2 := []int{}, []int{}

	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		re := regexp.MustCompile(`(\d*)   (\d*)$`)
		matches := re.FindAllStringSubmatch(line, -1)

		match1 := matches[0][1]
		match2 := matches[0][2]

		num1, _ := strconv.Atoi(match1)
		num2, _ := strconv.Atoi(match2)

		list1 = append(list1, num1)
		list2 = append(list2, num2)
	}

	sort.Ints(list1)
	sort.Ints(list2)

	totalDistance := 0
	for i := range list1 {
		if list1[i] > list2[i] {
			totalDistance += list1[i] - list2[i]
		} else {
			totalDistance += list2[i] - list1[i]
		}
	}
	fmt.Println(totalDistance)
}
