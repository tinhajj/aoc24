package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

// some notes about this day 2, it took an hour to solve and part 2 took a while
// because I made an assumption that you can be greedy and just remove one of
// the first two levels that causes the report to be bad.  that doesn't seem to
// be the case for some reason and I ended up brute forcing to get the answer
// quicker.
//
// missed the fact that the first level can be removed to create a report that is
// increasing instead of decreasing and vice versa. this means you only have to
// check 3 different reports not all of them
func main() {
	b, err := os.ReadFile("day_2_input.txt")
	if err != nil {
		panic(err)
	}
	input := string(b)

	safeReports := 0
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		report := []int{}
		for _, part := range parts {
			level, err := strconv.Atoi(part)
			if err != nil {
				panic(err)
			}
			report = append(report, level)
		}
		if isReportSafe(report, true) {
			safeReports++
		}
	}
	fmt.Println(safeReports)
}

func isReportSafe(report []int, tolerant bool) bool {
	min := 1
	max := 3

	difference := report[0] - report[1]
	if difference < 0 {
		min = -3
		max = -1
	}

	for i := 0; i < len(report)-1; i++ {
		j := i + 1
		difference := report[i] - report[j]

		if difference >= min && difference <= max {
		} else {
			if !tolerant {
				return false
			}

			newReport1 := slices.Concat(report[:i], report[i+1:])
			newReport2 := slices.Concat(report[:j], report[j+1:])
			newReport3 := report[1:]

			return isReportSafe(newReport1, false) || isReportSafe(newReport2, false) || isReportSafe(newReport3, false)
		}
	}
	return true
}
