package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"syscall"
)

var (
	SAMPLE   = false
	PATTERNS = []string{}
	CACHE    = map[string]Solution{}
)

type Solution struct {
	Solveable bool
	Ways      int
}

func main() {
	var err error
	var b []byte

	if false {
		f, err := os.Create("cpu.pprof")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		pprof.StartCPUProfile(f)

		// Setup signal handling
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigs
			pprof.StopCPUProfile()
			f.Close()
			os.Exit(0)
		}()
	}

	if SAMPLE {
		b, err = os.ReadFile("day_19_sample_input.txt")
	} else {
		b, err = os.ReadFile("day_19_input.txt")
	}
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(b), "\n")
	lines = lines[:len(lines)-1]

	tests := []string{}

	for _, chunk := range strings.Split(lines[0], ",") {
		PATTERNS = append(PATTERNS, strings.TrimSpace(chunk))
	}

	for _, test := range lines[2:] {
		tests = append(tests, test)
	}

	sum := 0
	for _, test := range tests {
		solution := solve(test)
		if solution.Solveable {
			sum += solution.Ways
		}
	}
	fmt.Println(sum)
}

func solve(test string) Solution {
	solution, ok := CACHE[test]
	if ok {
		return solution
	}

	if len(test) == 0 {
		return Solution{Solveable: true, Ways: 1}
	}

	solution = Solution{Solveable: false, Ways: 0}

	for _, pattern := range PATTERNS {
		if strings.HasPrefix(test, pattern) {
			subPattern := test[len(pattern):]
			subSolution := solve(subPattern)
			CACHE[subPattern] = subSolution

			if subSolution.Solveable {
				solution.Solveable = true
				solution.Ways += subSolution.Ways
			}
		}
	}
	CACHE[test] = solution

	return solution
}
