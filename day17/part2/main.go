package main

import (
	"aoc24/scan"
	"fmt"
	"math"
	"os"
	"strings"
)

var (
	DEBUG  = false
	SAMPLE = false
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

	var computer Computer
	RegisterA := 1433100000

	for ; !equal(computer.Output, initial.Program); RegisterA++ {
		if RegisterA%100_000 == 0 {
			fmt.Println("RegisterA:", RegisterA)
		}
		computer = initial
		computer.RegisterA = RegisterA

		if DEBUG {
			fmt.Println("Initial Registers")
			fmt.Println(computer.RegisterA, computer.RegisterB, computer.RegisterC, computer.Program)
			fmt.Println()
		}

		for computer.IP < len(computer.Program) {
			HandleInstruction(&computer)
		}

		if DEBUG {
			fmt.Println("Final Registers")
			fmt.Println(computer.RegisterA, computer.RegisterB, computer.RegisterC, computer.Program)
			fmt.Println()
		}

		if len(computer.Output) == 0 {
			fmt.Println("No output")
			return
		}

		if DEBUG {
			fmt.Println("Output")
			i := 0
			for ; i < len(computer.Output)-1; i++ {
				fmt.Printf("%d,", computer.Output[i])
			}
			fmt.Printf("%d", computer.Output[i])
			fmt.Println()
		}
	}

	fmt.Println(RegisterA - 1)
}

func HandleInstruction(computer *Computer) {
	instruction := computer.Program[computer.IP]
	computer.IP++

	switch instruction {
	case 0: // adv
		ops := Operands(computer, OpRegisterA, OpCombo)
		op1, op2 := ops[0], ops[1]

		computer.RegisterA = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
	case 1: // bxl
		ops := Operands(computer, OpRegisterB, OpLit)
		op1, op2 := ops[0], ops[1]

		computer.RegisterB = op1 ^ op2
	case 2: // bst
		ops := Operands(computer, OpCombo)
		op1 := ops[0]

		computer.RegisterB = op1 % 8
	case 3: // jnz
		if computer.RegisterA == 0 {
			Operands(computer, OpLit)
			return
		} else {
			ops := Operands(computer, OpLit)
			op1 := ops[0]

			computer.IP = op1
		}
	case 4: // bxc
		Operands(computer, OpLit)
		computer.RegisterB = computer.RegisterB ^ computer.RegisterC
	case 5: // out
		ops := Operands(computer, OpCombo)
		op1 := ops[0]

		computer.Output = append(computer.Output, op1%8)
	case 6: // bdv
		ops := Operands(computer, OpRegisterA, OpCombo)
		op1, op2 := ops[0], ops[1]

		computer.RegisterB = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
	case 7: // cdv
		ops := Operands(computer, OpRegisterA, OpCombo)
		op1, op2 := ops[0], ops[1]

		computer.RegisterC = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
	}
}

func Operands(computer *Computer, opKinds ...OpKind) []int {
	ops := []int{}

	for _, opKind := range opKinds {
		v := computer.Program[computer.IP]

		if opKind == OpCombo || opKind == OpLit {
			computer.IP++
			if opKind == OpCombo {
				if v >= 0 && v <= 3 {
					ops = append(ops, v)
					continue
				}

				switch v {
				case 4:
					ops = append(ops, computer.RegisterA)
				case 5:
					ops = append(ops, computer.RegisterB)
				case 6:
					ops = append(ops, computer.RegisterC)
				}
			} else {
				ops = append(ops, v)
			}
		} else {
			switch opKind {
			case OpRegisterA:
				ops = append(ops, computer.RegisterA)
			case OpRegisterB:
				ops = append(ops, computer.RegisterB)
			case OpRegisterC:
				ops = append(ops, computer.RegisterC)
			}
		}
	}

	return ops
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
