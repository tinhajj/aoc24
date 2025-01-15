/*
Not sure what to do for this one.  The JNZ at the end of the program always goes back to the beginning.
The registers get set to something like this:

B = (A % 8) ^ 7
C = (A / [2 ** [(A % 8) ^ 7]])
B = B ^ 7
A = A / 2 ** 3
B = B ^ C

Out = B % 8
Jump to start if Register A not 0

At the bare minimum Register A would have to be around 35184372088832 in order
to output 16 numbers.  We know this because we can consider the whole program a
loop and in each loop we divide Register A by 2 ** 3 (8).  So Register A needs to be
big enough to stay above 0 for many divisions.

Beyond that though, I am not sure what to do.

The entire program depends on Register A.  And for every loop that is true.

Maybe only some of the bits of Register A matter so we can figure it out in
parts and then combine those bits in the end to see what the answer is?

Also at the start of every loop B and C get set to relative to A.
So their values for the last loop are basically thrown out.
*/
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
	Initial Computer
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

func (c *Computer) Run(scratch [2]int, in [2]OpKind) {
	for c.IP < len(c.Program) {
		HandleInstruction(c, scratch, in)
	}
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

	if DEBUG && false {
		f, _ := os.Create("cpu.prof")
		pprof.StartCPUProfile(f)

		c := make(chan os.Signal, 1)
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

	full := scan.Numbers(lines[4])
	oneLoop := full[:len(full)-2]
	_ = oneLoop

	Initial := Computer{
		IP:        0,
		RegisterA: 0,
		RegisterB: 0,
		RegisterC: 0,
		Program:   full, // take the full program without the jump at the end
	}

	var in [2]OpKind
	var scratch [2]int

	RegisterA := 35184372088832
	RegisterA = 0b1000000000000000000000000000000000000000000000
	RegisterA = 0b11000001_00110100_01011000_00011100_01101010_000000

	fmt.Println("Computer stuff")
	for i := RegisterA; i <= RegisterA+10; i += 1 {
		sample := Initial
		sample.RegisterA = i
		sample.Run(scratch, in)
		fmt.Println(sample.Output)
	}

}

func HandleInstruction(computer *Computer, scratch [2]int, inputs [2]OpKind) (debug string) {
	sb := strings.Builder{}

	instruction := computer.Program[computer.IP]
	computer.IP++

	switch instruction {
	case 0: // adv
		inputs[0], inputs[1] = OpRegisterA, OpCombo
		ops, info := Operands(computer, scratch, inputs)

		op1, op2 := ops[0], ops[1]

		computer.RegisterA = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("adv RegisterA = %s / 2 ** %s", info[0], info[1]))
		}
	case 1: // bxl
		inputs[0], inputs[1] = OpRegisterB, OpLit
		ops, info := Operands(computer, scratch, inputs)
		op1, op2 := ops[0], ops[1]

		computer.RegisterB = op1 ^ op2
		if DEBUG {
			sb.WriteString(fmt.Sprintf("bxl RegisterB = %s ^ %s", info[0], info[1]))
		}
	case 2: // bst
		inputs[0], inputs[1] = OpCombo, OpNil
		ops, info := Operands(computer, scratch, inputs)
		op1 := ops[0]

		computer.RegisterB = op1 % 8
		if DEBUG {
			sb.WriteString(fmt.Sprintf("bst RegisterB = %s %% 8", info[0]))
		}
	case 3: // jnz
		if computer.RegisterA == 0 {
			inputs[0], inputs[1] = OpLit, OpNil
			Operands(computer, scratch, inputs)

			if DEBUG {
				sb.WriteString("jnz SKIPPED")
			}
		} else {
			inputs[0], inputs[1] = OpLit, OpNil
			ops, info := Operands(computer, scratch, inputs)
			op1 := ops[0]

			computer.IP = op1
			if DEBUG {
				sb.WriteString(fmt.Sprintf("jnz %s", info[0]))
			}
		}
	case 4: // bxc
		inputs[0], inputs[1] = OpLit, OpNil
		Operands(computer, scratch, inputs)
		computer.RegisterB = computer.RegisterB ^ computer.RegisterC

		if DEBUG {
			sb.WriteString(fmt.Sprintf("bxc RegisterB = (Register B) ^ (Register C)"))
		}
	case 5: // out
		inputs[0], inputs[1] = OpCombo, OpNil
		ops, info := Operands(computer, scratch, inputs)
		op1 := ops[0]

		computer.Output = append(computer.Output, op1%8)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("out %s %% 8", info[0]))
		}
	case 6: // bdv
		inputs[0], inputs[1] = OpRegisterA, OpCombo
		ops, info := Operands(computer, scratch, inputs)
		op1, op2 := ops[0], ops[1]

		computer.RegisterB = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("bdv RegisterB = %s / 2 ** %s", info[0], info[1]))
		}
	case 7: // cdv
		inputs[0], inputs[1] = OpRegisterA, OpCombo
		ops, info := Operands(computer, scratch, inputs)
		op1, op2 := ops[0], ops[1]

		computer.RegisterC = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("cdv RegisterC = %s / 2 ** %s", info[0], info[1]))
		}
	}

	return sb.String()
}

func Operands(computer *Computer, out [2]int, in [2]OpKind) ([2]int, [2]string) {
	debug := [2]string{}

	for i := 0; i < 2; i++ {
		opKind := in[i]

		switch opKind {
		case OpNil:
			continue
		case OpCombo:
			v := computer.Program[computer.IP]
			computer.IP++
			if v >= 0 && v <= 3 {
				out[i] = v
				debug[i] = fmt.Sprintf("(Combo Literal %d)", v)
				continue
			}

			switch v {
			case 4:
				out[i] = computer.RegisterA
				debug[i] = fmt.Sprintf("(Combo Register A %d)", computer.RegisterA)
			case 5:
				out[i] = computer.RegisterB
				debug[i] = fmt.Sprintf("(Combo Register B %d)", computer.RegisterB)
			case 6:
				out[i] = computer.RegisterC
				debug[i] = fmt.Sprintf("(Combo Register C %d)", computer.RegisterC)
			}
		case OpLit:
			v := computer.Program[computer.IP]
			computer.IP++
			out[i] = v
			debug[i] = fmt.Sprintf("(Literal %d)", v)
		case OpRegisterA:
			out[i] = computer.RegisterA
			debug[i] = fmt.Sprintf("(Register A %d)", computer.RegisterA)
		case OpRegisterB:
			out[i] = computer.RegisterB
			debug[i] = fmt.Sprintf("(Register B %d)", computer.RegisterB)
		case OpRegisterC:
			out[i] = computer.RegisterC
			debug[i] = fmt.Sprintf("(Register C %d)", computer.RegisterC)
		}
	}

	return out, debug
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
