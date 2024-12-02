package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

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

			// this checks all possible reports, but it seems like too much of a
			// brute force solution to be the most opitmal way to do it.
			bools := []bool{}
			for i := 0; i < len(report); i++ {
				newReport := slices.Concat(report[:i], report[i+1:])
				bools = append(bools, isReportSafe(newReport, false))
			}
			for _, b := range bools {
				if b {
					return true
				}
			}
			return false
		}
	}
	return true
}
