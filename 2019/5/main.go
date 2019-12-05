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
	inputs := []int{1}
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
	modes []bool
}

func (m *Modes) get(n int) bool {
	assert(n >= 0, "bad get")
	if n >= len(m.modes) {
		return false
	}
	return m.modes[n]
}

func parseOpcode(code int) (int, Modes, error) {
	s := strconv.Itoa(code)
	opcode := code % 100
	var modes []bool
	for i := len(s) - 3; i >= 0; i-- {
		// fmt.Printf("i=%v\n", i)
		// fmt.Printf("s[i]=%v, '0'=%v, ==?: %v\n", s[i], '0', s[i] == '0')
		if s[i] == '1' {
			modes = append(modes, true)
		} else if s[i] == '0' {
			modes = append(modes, false)
		} else {
			return 0, Modes{}, fmt.Errorf("Unrecognized mode %q[%d]='%c'", s, i, s[i])
		}
	}
	// fmt.Printf("done\n")
	return opcode, Modes{modes}, nil
}

func run(mem []int, inputs []int) ([]int, error) {
	inputsPtr := 0

	var res []int
	pc := -1
	var halt bool
	for !halt {
		code := chomp(mem, &pc)
		opcode, modes, err := parseOpcode(code)
		if err != nil {
			return nil, err
		}

		// TODO
		_ = modes

		// dump(mem, pc)
		switch opcode {
		case 1: // add
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			c := chomp(mem, &pc)
			// a = applyMode(modes.get(0))
			ensureInbounds(mem, a, b, c)
			av := mem[a]
			bv := mem[b]
			mem[c] = av + bv
		case 2: // mult
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			c := chomp(mem, &pc)
			ensureInbounds(mem, a, b, c)
			av := mem[a]
			bv := mem[b]
			mem[c] = av * bv
		case 3: // input
			a := chomp(mem, &pc)
			i := chomp(inputs, &inputsPtr)
			ensureInbounds(mem, a)
			mem[a] = i
		case 4: // output
			a := chomp(mem, &pc)
			ensureInbounds(mem, a)
			res = append(res, mem[a])
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
	assert(modes.get(0) == false, "t1 modes[0]")
	assert(modes.get(1) == true, "t1 modes[1]")
	assert(modes.get(2) == false, "t1 modes[2]")
	assert(modes.get(100) == false, "t1 modes[100]")

	opcode, modes, err = parseOpcode(3002)
	assert(err != nil, "t2 err")

	opcode, modes, err = parseOpcode(42)
	assert(err == nil, "t3 err")

	assert(false, "failfast")
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
