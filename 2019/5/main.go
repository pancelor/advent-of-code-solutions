package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func solve(memTemplate []int) (interface{}, error) {
	mem := dupmem(memTemplate)
	inputs := []int{5}
	res, err := run(mem, inputs)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Modes represent argument parameter modes
// get(0) == true means the first parameter should be interpreted in
//   immediate mode; false means position mode
type Modes struct {
	modes []int
}

func (m *Modes) get(n int) int {
	assert(n >= 0, "bad get")
	if n >= len(m.modes) {
		return 0
	}
	return m.modes[n]
}

func parseOpcode(code int) (int, Modes, error) {
	s := strconv.Itoa(code)
	opcode := code % 100
	var modes []int
	for i := len(s) - 3; i >= 0; i-- {
		// fmt.Printf("i=%v\n", i)
		// fmt.Printf("s[i]=%v, '0'=%v, ==?: %v\n", s[i], '0', s[i] == '0')
		if s[i] == '1' {
			modes = append(modes, 1)
		} else if s[i] == '0' {
			modes = append(modes, 0)
		} else {
			return 0, Modes{}, fmt.Errorf("Unrecognized mode %q[%d]='%c'", s, i, s[i])
		}
	}
	// fmt.Printf("done\n")
	return opcode, Modes{modes}, nil
}

func run(mem []int, inputs []int) ([]int, error) {
	inputsPtr := -1

	var res []int
	pc := -1
	var halt bool
	for !halt {
		code := chomp(mem, &pc)
		opcode, modes, err := parseOpcode(code)
		if err != nil {
			return nil, err
		}

		// dump(mem, pc)
		switch opcode {
		case 1: // add
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			c := chomp(mem, &pc)
			av := paramValue(mem, a, modes.get(0))
			bv := paramValue(mem, b, modes.get(1))
			ensureInbounds(mem, c)
			assert(modes.get(2) == 0, "immediate mode output param")

			mem[c] = av + bv
		case 2: // mult
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			c := chomp(mem, &pc)
			av := paramValue(mem, a, modes.get(0))
			bv := paramValue(mem, b, modes.get(1))
			ensureInbounds(mem, c)
			assert(modes.get(2) == 0, "immediate mode output param")

			mem[c] = av * bv
		case 3: // input
			a := chomp(mem, &pc)
			ensureInbounds(mem, a)
			assert(modes.get(0) == 0, "immediate mode output param")
			i := chomp(inputs, &inputsPtr)

			mem[a] = i
		case 4: // output
			a := chomp(mem, &pc)
			av := paramValue(mem, a, modes.get(0))

			res = append(res, av)
		case 5: // jump-if-true
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			av := paramValue(mem, a, modes.get(0))
			bv := paramValue(mem, b, modes.get(1))
			if av != 0 {
				pc = bv - 1 // pc will increment next chomp
			}
		case 6: // jump-if-false
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			av := paramValue(mem, a, modes.get(0))
			bv := paramValue(mem, b, modes.get(1))
			if av == 0 {
				pc = bv - 1 // pc will increment next chomp
			}
		case 7: // less than
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			c := chomp(mem, &pc)
			av := paramValue(mem, a, modes.get(0))
			bv := paramValue(mem, b, modes.get(1))
			assert(modes.get(2) == 0, "immediate mode output param")
			ensureInbounds(mem, c)

			if av < bv {
				mem[c] = 1
			} else {
				mem[c] = 0
			}
		case 8: // equals
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			c := chomp(mem, &pc)
			av := paramValue(mem, a, modes.get(0))
			bv := paramValue(mem, b, modes.get(1))
			assert(modes.get(2) == 0, "immediate mode output param")
			ensureInbounds(mem, c)

			if av == bv {
				mem[c] = 1
			} else {
				mem[c] = 0
			}
		case 99: // halt
			halt = true
		default:
			return nil, fmt.Errorf("Bad opcode %d at mem[%d]", mem[pc], pc)
		}
	}

	return res, nil
}

func dupmem(mem []int) []int {
	res := make([]int, len(mem))
	copy(res, mem)
	return res
}

func chomp(mem []int, pc *int) int {
	*pc++
	ensureInbounds(mem, *pc)
	return mem[*pc]
}

func paramValue(mem []int, param int, mode int) int {
	switch mode {
	case 0:
		ensureInbounds(mem, param)
		return mem[param]
	case 1:
		return param
	}
	assert(false, "bad paramValue")
	return 0
}

func ensureInbounds(mem []int, ptr ...int) {
	for _, p := range ptr {
		assert(0 <= p && p < len(mem), "oob")
	}
}

func dump(mem []int, pc int) {
	fmt.Printf("pc=%d, mem=[", pc)
	for i, v := range mem {
		if i%10 == 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("%3d ", v)
	}
	fmt.Printf("\n]\n")
}

func test() {
	opcode, modes, err := parseOpcode(1002)
	assert(err == nil, "t1 err")
	assert(opcode == 2, "t1 opcode")
	assert(modes.get(0) == 0, "t1 modes[0]")
	assert(modes.get(1) == 1, "t1 modes[1]")
	assert(modes.get(2) == 0, "t1 modes[2]")
	assert(modes.get(100) == 0, "t1 modes[100]")

	opcode, modes, err = parseOpcode(3002)
	assert(err != nil, "t2 err")

	opcode, modes, err = parseOpcode(42)
	assert(err == nil, "t3 err")

	// assert(false, "exit after tests")
}

func main() {
	test()

	input, err := getInput()
	if err != nil {
		panic(err)
	}

	answer, err := solve(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("answer:\n%v\n", answer)
}

//
// helpers
//

// for temporary use only
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// for temporary use only
func assert(b bool, msg string) {
	if !b {
		panic(errors.New(msg))
	}
}

var source = os.Stdin

func getInput() ([]int, error) {
	lines, err := getLines()
	if err != nil {
		return nil, err
	}

	var res []int
	for _, l := range strings.Split(lines[0], ",") {
		if l == "" {
			continue
		}

		v, err := strconv.Atoi(l)
		if err != nil {
			return nil, err
		}

		res = append(res, v)
	}

	return res, nil
}

func getLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
