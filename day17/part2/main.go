package main

import (
	"aoc24/scan"
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"sync"
	"syscall"
)

const (
	DEBUG  = false
	SAMPLE = false
)

var (
	Initial            Computer
	InitialProgramSize int
)

type Scope struct {
	Start int
	End   int
}

type Computer struct {
	IP int

	RegisterA int
	RegisterB int
	RegisterC int

	Program []int

	Output []int
}

type OpKind int

const (
	OpNil OpKind = iota
	OpLit
	OpCombo

	OpRegisterA
	OpRegisterB
	OpRegisterC
)

func main() {
	var b []byte
	var err error

	if DEBUG {
		f, _ := os.Create("cpu.prof")
		pprof.StartCPUProfile(f)

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			pprof.StopCPUProfile()
			os.Exit(1)
		}()
	}

	if SAMPLE {
		b, err = os.ReadFile("day_17_sample_input.txt")
	} else {
		b, err = os.ReadFile("day_17_input.txt")
	}
	if err != nil {
		panic(err)
	}

	input := string(b)
	lines := strings.Split(input, "\n")
	lines = lines[:len(lines)-1]

	Initial := Computer{
		IP:        0,
		RegisterA: scan.Numbers(lines[0])[0],
		RegisterB: scan.Numbers(lines[1])[0],
		RegisterC: scan.Numbers(lines[2])[0],
		Program:   scan.Numbers(lines[4]),
	}

	InitialProgramSize = len(Initial.Program)

	RegisterA := 0

	workerCount := 1
	workC := make(chan Scope, workerCount)
	doneC := make(chan int, workerCount)

	wg := &sync.WaitGroup{}
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go Worker(i, Initial, workC, doneC, wg)
	}

	var answer int
	go func() {
		for {
			select {
			case answer = <-doneC:
				close(workC)
				return
			default:
				size := 10_000_000
				scope := Scope{
					Start: RegisterA,
					End:   RegisterA + size,
				}
				RegisterA += size

				workC <- scope
			}
		}
	}()

	wg.Wait()

	fmt.Println("Answer:", answer)

}

func Worker(id int, initial Computer, work chan Scope, done chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	opValScratch := [2]int{}
	opKindScratch := [2]OpKind{}

	for scope := range work {
		for i := scope.Start; i <= scope.End; i++ {
			regA := i

			if regA%10_000_000 == 0 {
				fmt.Printf("Worker %d, Scope: %d\n", id, scope)
			}

			computer := initial
			computer.RegisterA = regA

			fmt.Println("New Computer")
			for computer.IP < len(computer.Program) {
				HandleInstruction(&computer, opValScratch, opKindScratch)
				if !equalPre(computer.Output, initial.Program) {
					continue
				}
			}

			if equal(computer.Output, initial.Program) {
				done <- regA
				return
			}
		}
	}
}

func HandleInstruction(computer *Computer, opValScratch [2]int, opKindScratch [2]OpKind) (debug string) {
	sb := strings.Builder{}

	instruction := computer.Program[computer.IP]
	computer.IP++

	switch instruction {
	case 0: // adv
		opKindScratch[0], opKindScratch[1] = OpRegisterA, OpCombo
		ops := Operands(computer, opValScratch, opKindScratch)

		op1, op2 := ops[0], ops[1]

		computer.RegisterA = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("adv RegisterA, %d", op2))
		}
	case 1: // bxl
		opKindScratch[0], opKindScratch[1] = OpRegisterB, OpLit
		ops := Operands(computer, opValScratch, opKindScratch)
		op1, op2 := ops[0], ops[1]

		computer.RegisterB = op1 ^ op2
		if DEBUG {
			sb.WriteString(fmt.Sprintf("bxl RegisterB, %d", op2))
		}
	case 2: // bst
		opKindScratch[0], opKindScratch[1] = OpCombo, OpNil
		ops := Operands(computer, opValScratch, opKindScratch)
		op1 := ops[0]

		computer.RegisterB = op1 % 8
		if DEBUG {
			sb.WriteString(fmt.Sprintf("bst RegisterB, %d", op1))
		}
	case 3: // jnz
		if computer.RegisterA == 0 {
			opKindScratch[0], opKindScratch[1] = OpLit, OpNil
			Operands(computer, opValScratch, opKindScratch)

			if DEBUG {
				sb.WriteString("jnz SKIPPED")
			}
		} else {
			opKindScratch[0], opKindScratch[1] = OpLit, OpNil
			ops := Operands(computer, opValScratch, opKindScratch)
			op1 := ops[0]

			computer.IP = op1
			if DEBUG {
				sb.WriteString(fmt.Sprintf("jnz %d", op1))
			}
		}
	case 4: // bxc
		opKindScratch[0], opKindScratch[1] = OpLit, OpNil
		Operands(computer, opValScratch, opKindScratch)
		computer.RegisterB = computer.RegisterB ^ computer.RegisterC

		if DEBUG {
			sb.WriteString("bxc RegisterB, RegisterC")
		}
	case 5: // out
		opKindScratch[0], opKindScratch[1] = OpCombo, OpNil
		ops := Operands(computer, opValScratch, opKindScratch)
		op1 := ops[0]

		computer.Output = append(computer.Output, op1%8)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("out %d", op1))
		}
	case 6: // bdv
		opKindScratch[0], opKindScratch[1] = OpRegisterA, OpCombo
		ops := Operands(computer, opValScratch, opKindScratch)
		op1, op2 := ops[0], ops[1]

		computer.RegisterB = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("bdv RegisterB, %d", op2))
		}
	case 7: // cdv
		opKindScratch[0], opKindScratch[1] = OpRegisterA, OpCombo
		ops := Operands(computer, opValScratch, opKindScratch)
		op1, op2 := ops[0], ops[1]

		computer.RegisterC = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("cdv RegisterC, %d", op2))
		}
	}

	return sb.String()
}

func Operands(computer *Computer, scratch [2]int, opKinds [2]OpKind) [2]int {
	for i := 0; i < 2; i++ {
		opKind := opKinds[i]

		switch opKind {
		case OpCombo:
			v := computer.Program[computer.IP]
			computer.IP++
			if v >= 0 && v <= 3 {
				scratch[i] = v
				continue
			}

			switch v {
			case 4:
				scratch[i] = computer.RegisterA
			case 5:
				scratch[i] = computer.RegisterB
			case 6:
				scratch[i] = computer.RegisterC
			}
		case OpLit:
			v := computer.Program[computer.IP]
			computer.IP++
			scratch[i] = v
		case OpRegisterA:
			scratch[i] = computer.RegisterA
		case OpRegisterB:
			scratch[i] = computer.RegisterB
		case OpRegisterC:
			scratch[i] = computer.RegisterC
		}
	}

	return scratch
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func equalPre(a, b []int) bool {
	if len(a) > len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
