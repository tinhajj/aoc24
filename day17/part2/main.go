package main

import (
	"aoc24/scan"
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"syscall"
)

const (
	DEBUG  = false
	SAMPLE = false
)

var (
	InitialProgramSize int
	OperandScratch     [2]int
)

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
	OpLit OpKind = iota
	OpCombo

	OpRegisterA
	OpRegisterB
	OpRegisterC
)

func main() {
	var b []byte
	var err error

	f, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(f)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		pprof.StopCPUProfile()
		os.Exit(1)
	}()

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

	initial := Computer{
		IP:        0,
		RegisterA: scan.Numbers(lines[0])[0],
		RegisterB: scan.Numbers(lines[1])[0],
		RegisterC: scan.Numbers(lines[2])[0],
		Program:   scan.Numbers(lines[4]),
	}

	InitialProgramSize = len(initial.Program)

	var computer Computer
	RegisterA := 1460100000

	for ; !equal(computer.Output, initial.Program); RegisterA++ {
		if DEBUG && false {
			fmt.Println("New Computer")
		}
		if RegisterA%100_000 == 0 {
			fmt.Println("RegisterA:", RegisterA)
		}
		computer = initial
		computer.RegisterA = RegisterA

		if DEBUG && false {
			fmt.Println("Initial Registers")
			fmt.Println(computer.RegisterA, computer.RegisterB, computer.RegisterC, computer.Program)
			fmt.Println()
		}

		for computer.IP < len(computer.Program) {
			HandleInstruction(&computer)
			if !equalPre(computer.Output, initial.Program) {
				continue
			}
		}

		if DEBUG && false {
			fmt.Println("Final Registers")
			fmt.Println(computer.RegisterA, computer.RegisterB, computer.RegisterC, computer.Program)
			fmt.Println()
		}

		if DEBUG && false {
			fmt.Println("Output")
			i := 0
			for ; i < len(computer.Output)-1; i++ {
				fmt.Printf("%d,", computer.Output[i])
			}
			fmt.Printf("%d", computer.Output[i])
			fmt.Println()
		}
	}

	fmt.Println("Answer:", RegisterA-1)
}

func HandleInstruction(computer *Computer) (debug string) {
	sb := strings.Builder{}

	instruction := computer.Program[computer.IP]
	computer.IP++

	switch instruction {
	case 0: // adv
		ops := Operands(computer, OpRegisterA, OpCombo)
		op1, op2 := ops[0], ops[1]

		computer.RegisterA = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("adv RegisterA, %d", op2))
		}
	case 1: // bxl
		ops := Operands(computer, OpRegisterB, OpLit)
		op1, op2 := ops[0], ops[1]

		computer.RegisterB = op1 ^ op2
		if DEBUG {
			sb.WriteString(fmt.Sprintf("bxl RegisterB, %d", op2))
		}
	case 2: // bst
		ops := Operands(computer, OpCombo)
		op1 := ops[0]

		computer.RegisterB = op1 % 8
		if DEBUG {
			sb.WriteString(fmt.Sprintf("bst RegisterB, %d", op1))
		}
	case 3: // jnz
		if computer.RegisterA == 0 {
			Operands(computer, OpLit)

			if DEBUG {
				sb.WriteString("jnz SKIPPED")
			}
		} else {
			ops := Operands(computer, OpLit)
			op1 := ops[0]

			computer.IP = op1
			if DEBUG {
				sb.WriteString(fmt.Sprintf("jnz %d", op1))
			}
		}
	case 4: // bxc
		Operands(computer, OpLit)
		computer.RegisterB = computer.RegisterB ^ computer.RegisterC

		if DEBUG {
			sb.WriteString("bxc RegisterB, RegisterC")
		}
	case 5: // out
		ops := Operands(computer, OpCombo)
		op1 := ops[0]

		computer.Output = append(computer.Output, op1%8)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("out %d", op1))
		}
	case 6: // bdv
		ops := Operands(computer, OpRegisterA, OpCombo)
		op1, op2 := ops[0], ops[1]

		computer.RegisterB = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("bdv RegisterB, %d", op2))
		}
	case 7: // cdv
		ops := Operands(computer, OpRegisterA, OpCombo)
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

func Operands(computer *Computer, opKinds ...OpKind) [2]int {
	i := 0

	for j := 0; j < len(opKinds); i, j = j+1, i+1 {
		opKind := opKinds[j]

		v := computer.Program[computer.IP]

		if opKind == OpCombo || opKind == OpLit {
			computer.IP++
			if opKind == OpCombo {
				if v >= 0 && v <= 3 {
					OperandScratch[i] = v
					continue
				}

				switch v {
				case 4:
					OperandScratch[i] = computer.RegisterA
				case 5:
					OperandScratch[i] = computer.RegisterB
				case 6:
					OperandScratch[i] = computer.RegisterC
				}
			} else {
				OperandScratch[i] = v
			}
		} else {
			switch opKind {
			case OpRegisterA:
				OperandScratch[i] = computer.RegisterA
			case OpRegisterB:
				OperandScratch[i] = computer.RegisterB
			case OpRegisterC:
				OperandScratch[i] = computer.RegisterC
			}
		}
	}

	return OperandScratch
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
