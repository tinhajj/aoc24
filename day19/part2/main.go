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
	CACHE    = map[string]bool{}
)

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
		ok := solve(test)
		if ok {
			sum++
			fmt.Println(test)
		}
	}
	fmt.Println(sum)
}

func solve(test string) bool {
	solveable, ok := CACHE[test]
	if ok {
		return solveable
	}

	if len(test) == 0 {
		return true
	}

	solveable = false
	for _, pattern := range PATTERNS {
		if strings.HasPrefix(test, pattern) {
			sub := test[len(pattern):]

			ok := solve(sub)
			CACHE[sub] = ok

			if ok {
				solveable = true
			}
		}
	}

	return solveable
}
