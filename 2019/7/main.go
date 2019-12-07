package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var allPerms = [][5]int{
	[5]int{0, 1, 2, 3, 4},
	[5]int{0, 1, 2, 4, 3},
	[5]int{0, 1, 3, 2, 4},
	[5]int{0, 1, 3, 4, 2},
	[5]int{0, 1, 4, 2, 3},
	[5]int{0, 1, 4, 3, 2},
	[5]int{0, 2, 1, 3, 4},
	[5]int{0, 2, 1, 4, 3},
	[5]int{0, 2, 3, 1, 4},
	[5]int{0, 2, 3, 4, 1},
	[5]int{0, 2, 4, 1, 3},
	[5]int{0, 2, 4, 3, 1},
	[5]int{0, 3, 1, 2, 4},
	[5]int{0, 3, 1, 4, 2},
	[5]int{0, 3, 2, 1, 4},
	[5]int{0, 3, 2, 4, 1},
	[5]int{0, 3, 4, 1, 2},
	[5]int{0, 3, 4, 2, 1},
	[5]int{0, 4, 1, 2, 3},
	[5]int{0, 4, 1, 3, 2},
	[5]int{0, 4, 2, 1, 3},
	[5]int{0, 4, 2, 3, 1},
	[5]int{0, 4, 3, 1, 2},
	[5]int{0, 4, 3, 2, 1},
	[5]int{1, 0, 2, 3, 4},
	[5]int{1, 0, 2, 4, 3},
	[5]int{1, 0, 3, 2, 4},
	[5]int{1, 0, 3, 4, 2},
	[5]int{1, 0, 4, 2, 3},
	[5]int{1, 0, 4, 3, 2},
	[5]int{1, 2, 0, 3, 4},
	[5]int{1, 2, 0, 4, 3},
	[5]int{1, 2, 3, 0, 4},
	[5]int{1, 2, 3, 4, 0},
	[5]int{1, 2, 4, 0, 3},
	[5]int{1, 2, 4, 3, 0},
	[5]int{1, 3, 0, 2, 4},
	[5]int{1, 3, 0, 4, 2},
	[5]int{1, 3, 2, 0, 4},
	[5]int{1, 3, 2, 4, 0},
	[5]int{1, 3, 4, 0, 2},
	[5]int{1, 3, 4, 2, 0},
	[5]int{1, 4, 0, 2, 3},
	[5]int{1, 4, 0, 3, 2},
	[5]int{1, 4, 2, 0, 3},
	[5]int{1, 4, 2, 3, 0},
	[5]int{1, 4, 3, 0, 2},
	[5]int{1, 4, 3, 2, 0},
	[5]int{2, 0, 1, 3, 4},
	[5]int{2, 0, 1, 4, 3},
	[5]int{2, 0, 3, 1, 4},
	[5]int{2, 0, 3, 4, 1},
	[5]int{2, 0, 4, 1, 3},
	[5]int{2, 0, 4, 3, 1},
	[5]int{2, 1, 0, 3, 4},
	[5]int{2, 1, 0, 4, 3},
	[5]int{2, 1, 3, 0, 4},
	[5]int{2, 1, 3, 4, 0},
	[5]int{2, 1, 4, 0, 3},
	[5]int{2, 1, 4, 3, 0},
	[5]int{2, 3, 0, 1, 4},
	[5]int{2, 3, 0, 4, 1},
	[5]int{2, 3, 1, 0, 4},
	[5]int{2, 3, 1, 4, 0},
	[5]int{2, 3, 4, 0, 1},
	[5]int{2, 3, 4, 1, 0},
	[5]int{2, 4, 0, 1, 3},
	[5]int{2, 4, 0, 3, 1},
	[5]int{2, 4, 1, 0, 3},
	[5]int{2, 4, 1, 3, 0},
	[5]int{2, 4, 3, 0, 1},
	[5]int{2, 4, 3, 1, 0},
	[5]int{3, 0, 1, 2, 4},
	[5]int{3, 0, 1, 4, 2},
	[5]int{3, 0, 2, 1, 4},
	[5]int{3, 0, 2, 4, 1},
	[5]int{3, 0, 4, 1, 2},
	[5]int{3, 0, 4, 2, 1},
	[5]int{3, 1, 0, 2, 4},
	[5]int{3, 1, 0, 4, 2},
	[5]int{3, 1, 2, 0, 4},
	[5]int{3, 1, 2, 4, 0},
	[5]int{3, 1, 4, 0, 2},
	[5]int{3, 1, 4, 2, 0},
	[5]int{3, 2, 0, 1, 4},
	[5]int{3, 2, 0, 4, 1},
	[5]int{3, 2, 1, 0, 4},
	[5]int{3, 2, 1, 4, 0},
	[5]int{3, 2, 4, 0, 1},
	[5]int{3, 2, 4, 1, 0},
	[5]int{3, 4, 0, 1, 2},
	[5]int{3, 4, 0, 2, 1},
	[5]int{3, 4, 1, 0, 2},
	[5]int{3, 4, 1, 2, 0},
	[5]int{3, 4, 2, 0, 1},
	[5]int{3, 4, 2, 1, 0},
	[5]int{4, 0, 1, 2, 3},
	[5]int{4, 0, 1, 3, 2},
	[5]int{4, 0, 2, 1, 3},
	[5]int{4, 0, 2, 3, 1},
	[5]int{4, 0, 3, 1, 2},
	[5]int{4, 0, 3, 2, 1},
	[5]int{4, 1, 0, 2, 3},
	[5]int{4, 1, 0, 3, 2},
	[5]int{4, 1, 2, 0, 3},
	[5]int{4, 1, 2, 3, 0},
	[5]int{4, 1, 3, 0, 2},
	[5]int{4, 1, 3, 2, 0},
	[5]int{4, 2, 0, 1, 3},
	[5]int{4, 2, 0, 3, 1},
	[5]int{4, 2, 1, 0, 3},
	[5]int{4, 2, 1, 3, 0},
	[5]int{4, 2, 3, 0, 1},
	[5]int{4, 2, 3, 1, 0},
	[5]int{4, 3, 0, 1, 2},
	[5]int{4, 3, 0, 2, 1},
	[5]int{4, 3, 1, 0, 2},
	[5]int{4, 3, 1, 2, 0},
	[5]int{4, 3, 2, 0, 1},
	[5]int{4, 3, 2, 1, 0},
}

func solve(memTemplate []int) (interface{}, error) {
	best := 0
	for _, settings := range allPerms {
		assert(len(settings) == 5, "bad settings len")
		aOut, err := run(dupmem(memTemplate), []int{settings[0], 0})
		check(err)
		assert(len(aOut) == 1, "bad aOut len")
		bOut, err := run(dupmem(memTemplate), []int{settings[1], aOut[0]})
		check(err)
		assert(len(bOut) == 1, "bad bOut len")
		cOut, err := run(dupmem(memTemplate), []int{settings[2], bOut[0]})
		check(err)
		assert(len(cOut) == 1, "bad cOut len")
		dOut, err := run(dupmem(memTemplate), []int{settings[3], cOut[0]})
		check(err)
		assert(len(dOut) == 1, "bad dOut len")
		eOut, err := run(dupmem(memTemplate), []int{settings[4], dOut[0]})
		check(err)
		assert(len(eOut) == 1, "bad eOut len")
		res := eOut[0]
		if res > best {
			best = res
		}
	}

	return best, nil
}

// Modes represent argument parameter modes
// get(0) == 1 means the first parameter should be interpreted in
//   immediate mode; 0 means position mode
type Modes struct {
	modes []int
	ptr   int
}

func (m *Modes) get(n int) int {
	assert(n >= 0, "bad get")
	if n >= len(m.modes) {
		return 0
	}
	return m.modes[n]
}

func (m *Modes) getNext() int {
	m.ptr++
	return m.get(m.ptr - 1)
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
	return opcode, Modes{modes: modes}, nil
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
			av := paramValue(mem, a, modes.getNext())
			bv := paramValue(mem, b, modes.getNext())
			ensureInbounds(mem, c)
			assert(modes.getNext() == 0, "immediate mode output param")

			mem[c] = av + bv
		case 2: // mult
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			c := chomp(mem, &pc)
			av := paramValue(mem, a, modes.getNext())
			bv := paramValue(mem, b, modes.getNext())
			ensureInbounds(mem, c)
			assert(modes.getNext() == 0, "immediate mode output param")

			mem[c] = av * bv
		case 3: // input
			a := chomp(mem, &pc)
			ensureInbounds(mem, a)
			assert(modes.getNext() == 0, "immediate mode output param")
			i := chomp(inputs, &inputsPtr)

			mem[a] = i
		case 4: // output
			a := chomp(mem, &pc)
			av := paramValue(mem, a, modes.getNext())

			res = append(res, av)
		case 5: // jump-if-true
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			av := paramValue(mem, a, modes.getNext())
			bv := paramValue(mem, b, modes.getNext())
			if av != 0 {
				pc = bv - 1 // pc will increment next chomp
			}
		case 6: // jump-if-false
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			av := paramValue(mem, a, modes.getNext())
			bv := paramValue(mem, b, modes.getNext())
			if av == 0 {
				pc = bv - 1 // pc will increment next chomp
			}
		case 7: // less than
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			c := chomp(mem, &pc)
			av := paramValue(mem, a, modes.getNext())
			bv := paramValue(mem, b, modes.getNext())
			assert(modes.getNext() == 0, "immediate mode output param")
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
			av := paramValue(mem, a, modes.getNext())
			bv := paramValue(mem, b, modes.getNext())
			assert(modes.getNext() == 0, "immediate mode output param")
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
	assert(modes.getNext() == 0, "t1 modes[0]")
	assert(modes.getNext() == 1, "t1 modes[1]")
	assert(modes.getNext() == 0, "t1 modes[2]")
	assert(modes.getNext() == 0, "t1 modes[3]")
	assert(modes.getNext() == 0, "t1 modes[4]")
	assert(modes.getNext() == 0, "t1 modes[5]")

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
