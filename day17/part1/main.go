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

var (
	IP int = 0

	RegisterA int
	RegisterB int
	RegisterC int

	Program []int

	Output []int
)

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

	RegisterA = scan.Numbers(lines[0])[0]
	RegisterB = scan.Numbers(lines[1])[0]
	RegisterC = scan.Numbers(lines[2])[0]

	Program = scan.Numbers(lines[4])

	if DEBUG {
		fmt.Println("Initial Registers")
		fmt.Println(RegisterA, RegisterB, RegisterC, Program)
		fmt.Println()
	}

	for IP < len(Program) {
		HandleInstruction()
	}

	if DEBUG {
		fmt.Println("Final Registers")
		fmt.Println(RegisterA, RegisterB, RegisterC, Program)
		fmt.Println()
	}

	if len(Output) == 0 {
		fmt.Println("No output")
		return
	}

	fmt.Println("Output")
	i := 0
	for ; i < len(Output)-1; i++ {
		fmt.Printf("%d,", Output[i])
	}
	fmt.Printf("%d", Output[i])
}

func HandleInstruction() {
	instruction := Program[IP]
	IP++

	switch instruction {
	case 0: // adv
		ops := Operands(OpRegisterA, OpCombo)
		op1, op2 := ops[0], ops[1]

		RegisterA = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
	case 1: // bxl
		ops := Operands(OpRegisterB, OpLit)
		op1, op2 := ops[0], ops[1]

		RegisterB = op1 ^ op2
	case 2: // bst
		ops := Operands(OpCombo)
		op1 := ops[0]

		RegisterB = op1 % 8
	case 3: // jnz
		if RegisterA == 0 {
			Operands(OpLit)
			return
		} else {
			ops := Operands(OpLit)
			op1 := ops[0]

			IP = op1
		}
	case 4: // bxc
		Operands(OpLit)
		RegisterB = RegisterB ^ RegisterC
	case 5: // out
		ops := Operands(OpCombo)
		op1 := ops[0]

		Output = append(Output, op1%8)
	case 6: // bdv
		ops := Operands(OpRegisterA, OpCombo)
		op1, op2 := ops[0], ops[1]

		RegisterB = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
	case 7: // cdv
		ops := Operands(OpRegisterA, OpCombo)
		op1, op2 := ops[0], ops[1]

		RegisterC = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
	}
}

func Operands(opKinds ...OpKind) []int {
	ops := []int{}

	for _, opKind := range opKinds {
		v := Program[IP]

		if opKind == OpCombo || opKind == OpLit {
			IP++
			if opKind == OpCombo {
				if v >= 0 && v <= 3 {
					ops = append(ops, v)
					continue
				}

				switch v {
				case 4:
					ops = append(ops, RegisterA)
				case 5:
					ops = append(ops, RegisterB)
				case 6:
					ops = append(ops, RegisterC)
				}
			} else {
				ops = append(ops, v)
			}
		} else {
			switch opKind {
			case OpRegisterA:
				ops = append(ops, RegisterA)
			case OpRegisterB:
				ops = append(ops, RegisterB)
			case OpRegisterC:
				ops = append(ops, RegisterC)
			}
		}
	}

	return ops
}
