package main

import (
	"fmt"
	"os"
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
		if isReportSafe(report) {
			safeReports++
		}
	}
	fmt.Println(safeReports)
}

func isReportSafe(report []int) bool {
	min := 1
	max := 3

	difference := report[0] - report[1]
	if difference < 0 {
		min = -3
		max = -1
	}

	for i := 0; i < len(report)-1; i++ {
		difference := report[i] - report[i+1]

		if difference >= min && difference <= max {
		} else {
			return false
		}
	}
	return true
}
