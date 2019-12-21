package main

// Use this file as a template when you have an ascii intcode program to run

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pancelor/advent-of-code-solutions/2019/computer"
	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
)

var assert = helpers.Assert
var check = helpers.Check

func solve(in []int) interface{} {
	assert(len(in) > 0, "you forgot to paste in the input")
	cpu := computer.MakeCPU("ed")
	cpu.SetMemory(in)
	// fmt.Println(cpu.PrintProgram())
	cpu.Run()

	scanner := bufio.NewScanner(os.Stdin)
	var lastOut int
	for {
		state := <-cpu.StateChan
		switch state {
		case computer.CS_WAITING_INPUT:
			fmt.Printf("> ")
			assert(scanner.Scan(), "couldn't scan")
			str := scanner.Text() + "\n"
			// fmt.Printf("sending %q\n", str)
			cpu.SendAsciiInput(str)
		case computer.CS_WAITING_OUTPUT:
			out := <-cpu.OutChan
			lastOut = out
			fmt.Print(string(out))
		case computer.CS_DONE:
			return lastOut
		}
	}
}

func main() {
	fmt.Printf("answer:\n%v\n", solve(input))
}

var input = []int{
	// paste input here
}
