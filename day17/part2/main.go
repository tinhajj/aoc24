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
loop and in each loop we divide Register A by 2 ** 3.  So Register A needs to be
big enough to stay above 0 for many divisions.

Beyond that though, I am not sure what to do.
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
	"sync"
	"syscall"
)

const (
	DEBUG  = true
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

	if DEBUG && false {
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

	inScratch := [2]int{}
	outScratch := [2]OpKind{}

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
				d := HandleInstruction(&computer, inScratch, outScratch)
				fmt.Println(d)
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

func HandleInstruction(computer *Computer, outScratch [2]int, inScratch [2]OpKind) (debug string) {
	sb := strings.Builder{}

	instruction := computer.Program[computer.IP]
	computer.IP++

	switch instruction {
	case 0: // adv
		inScratch[0], inScratch[1] = OpRegisterA, OpCombo
		ops, info := Operands(computer, outScratch, inScratch)

		op1, op2 := ops[0], ops[1]

		computer.RegisterA = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("adv RegisterA = %s / 2 ** %s", info[0], info[1]))
		}
	case 1: // bxl
		inScratch[0], inScratch[1] = OpRegisterB, OpLit
		ops, info := Operands(computer, outScratch, inScratch)
		op1, op2 := ops[0], ops[1]

		computer.RegisterB = op1 ^ op2
		if DEBUG {
			sb.WriteString(fmt.Sprintf("bxl RegisterB = %s ^ %s", info[0], info[1]))
		}
	case 2: // bst
		inScratch[0], inScratch[1] = OpCombo, OpNil
		ops, info := Operands(computer, outScratch, inScratch)
		op1 := ops[0]

		computer.RegisterB = op1 % 8
		if DEBUG {
			sb.WriteString(fmt.Sprintf("bst RegisterB = %s %% 8", info[0]))
		}
	case 3: // jnz
		if computer.RegisterA == 0 {
			inScratch[0], inScratch[1] = OpLit, OpNil
			Operands(computer, outScratch, inScratch)

			if DEBUG {
				sb.WriteString("jnz SKIPPED")
			}
		} else {
			inScratch[0], inScratch[1] = OpLit, OpNil
			ops, info := Operands(computer, outScratch, inScratch)
			op1 := ops[0]

			computer.IP = op1
			if DEBUG {
				sb.WriteString(fmt.Sprintf("jnz %s", info[0]))
			}
		}
	case 4: // bxc
		inScratch[0], inScratch[1] = OpLit, OpNil
		Operands(computer, outScratch, inScratch)
		computer.RegisterB = computer.RegisterB ^ computer.RegisterC

		if DEBUG {
			sb.WriteString(fmt.Sprintf("bxc RegisterB = (Register B) ^ (Register C)"))
		}
	case 5: // out
		inScratch[0], inScratch[1] = OpCombo, OpNil
		ops, info := Operands(computer, outScratch, inScratch)
		op1 := ops[0]

		computer.Output = append(computer.Output, op1%8)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("out %s %% 8", info[0]))
		}
	case 6: // bdv
		inScratch[0], inScratch[1] = OpRegisterA, OpCombo
		ops, info := Operands(computer, outScratch, inScratch)
		op1, op2 := ops[0], ops[1]

		computer.RegisterB = int(
			float64(op1) / math.Pow(2.0, float64(op2)),
		)
		if DEBUG {
			sb.WriteString(fmt.Sprintf("bdv RegisterB = %s / 2 ** %s", info[0], info[1]))
		}
	case 7: // cdv
		inScratch[0], inScratch[1] = OpRegisterA, OpCombo
		ops, info := Operands(computer, outScratch, inScratch)
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
