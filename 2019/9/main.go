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
	in := make(chan int, 10)
	in <- 81
	out, done := run("A", mem, in)

	var res []int
	for {
		select {
		case <-done:
			return res, nil
		case x := <-out:
			res = append(res, x)
		}
	}
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

func paramValue(mem []int, param int, relBase int, mode int) int {
	switch mode {
	case 0:
		ensureInbounds(mem, param)
		return mem[param]
	case 1:
		return param
	case 2:
		ensureInbounds(mem, relBase+param)
		return mem[relBase+param]
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
		} else if s[i] == '2' {
			modes = append(modes, 2)
		} else {
			return 0, Modes{}, fmt.Errorf("Unrecognized mode %q[%d]='%c'", s, i, s[i])
		}
	}
	// fmt.Printf("done\n")
	return opcode, Modes{modes: modes}, nil
}

func run(name string, mem []int, inCh chan int) (chan int, chan struct{}) {
	outCh := make(chan int)
	doneCh := make(chan struct{})
	go func() {
		relBase := 0
		pc := -1
		var halt bool
		for !halt {
			var log strings.Builder
			code := chomp(mem, &pc)
			opcode, modes, err := parseOpcode(code)
			check(err)

			// fmt.Printf("node %s pc=%v code=%d modes=%v\n", name, pc, code, modes)
			// dump(mem, pc)

			fmt.Fprintf(&log, "pc=%d ", pc)
			fmt.Fprintf(&log, "code=%d ", code)
			switch opcode {
			case 1: // add
				a := chomp(mem, &pc)
				b := chomp(mem, &pc)
				c := chomp(mem, &pc)
				av := paramValue(mem, a, relBase, modes.getNext())
				bv := paramValue(mem, b, relBase, modes.getNext())
				ensureInbounds(mem, c)

				var cv int
				switch modes.getNext() {
				case 0:
					cv = c
				case 1:
					assert(false, "immediate mode output param")
				case 2:
					cv = c + relBase
				}

				fmt.Fprintf(&log, "add(%d,%d,%d): %d+%d->%d ", a, b, c, av, bv, cv)
				mem[cv] = av + bv
			case 2: // mult
				a := chomp(mem, &pc)
				b := chomp(mem, &pc)
				c := chomp(mem, &pc)
				av := paramValue(mem, a, relBase, modes.getNext())
				bv := paramValue(mem, b, relBase, modes.getNext())
				ensureInbounds(mem, c)

				var cv int
				switch modes.getNext() {
				case 0:
					cv = c
				case 1:
					assert(false, "immediate mode output param")
				case 2:
					cv = c + relBase
				}

				fmt.Fprintf(&log, "mult(%d,%d,%d): %d+%d->%d ", a, b, c, av, bv, cv)
				mem[cv] = av * bv
			case 3: // input
				a := chomp(mem, &pc)
				ensureInbounds(mem, a)
				i := <-inCh
				// fmt.Printf("%s < %d\n", name, i)

				var av int
				switch modes.getNext() {
				case 0:
					av = a
				case 1:
					assert(false, "immediate mode output param")
				case 2:
					av = a + relBase
				}

				fmt.Fprintf(&log, "input(%d): %d->%d ", a, i, av)
				mem[av] = i
			case 4: // output
				a := chomp(mem, &pc)
				av := paramValue(mem, a, relBase, modes.getNext())

				// fmt.Printf("%s : %d\n", name, av)
				fmt.Fprintf(&log, "output(%d): ->%d ", a, av)

				outCh <- av
			case 5: // jump-if-true
				a := chomp(mem, &pc)
				b := chomp(mem, &pc)
				av := paramValue(mem, a, relBase, modes.getNext())
				bv := paramValue(mem, b, relBase, modes.getNext())

				fmt.Fprintf(&log, "jump-if-true(%d,%d): ", a, b)
				if av != 0 {
					fmt.Fprintf(&log, "%d!=0; jump to %d ", av, bv)
					pc = bv - 1 // pc will increment next chomp
				}
			case 6: // jump-if-false
				a := chomp(mem, &pc)
				b := chomp(mem, &pc)
				av := paramValue(mem, a, relBase, modes.getNext())
				bv := paramValue(mem, b, relBase, modes.getNext())

				fmt.Fprintf(&log, "jump-if-false(%d,%d): ", a, b)
				if av == 0 {
					fmt.Fprintf(&log, "%d==0; jump to %d ", av, bv)
					pc = bv - 1 // pc will increment next chomp
				}
			case 7: // less than
				a := chomp(mem, &pc)
				b := chomp(mem, &pc)
				c := chomp(mem, &pc)
				av := paramValue(mem, a, relBase, modes.getNext())
				bv := paramValue(mem, b, relBase, modes.getNext())

				var cv int
				switch modes.getNext() {
				case 0:
					cv = c
				case 1:
					assert(false, "immediate mode output param")
				case 2:
					cv = c + relBase
				}
				ensureInbounds(mem, cv)

				fmt.Fprintf(&log, "less-than(%d,%d,%d): %d<%d -> %d ", a, b, c, av, bv, cv)
				if av < bv {
					mem[cv] = 1
				} else {
					mem[cv] = 0
				}
			case 8: // equals
				a := chomp(mem, &pc)
				b := chomp(mem, &pc)
				c := chomp(mem, &pc)
				av := paramValue(mem, a, relBase, modes.getNext())
				bv := paramValue(mem, b, relBase, modes.getNext())

				var cv int
				switch modes.getNext() {
				case 0:
					cv = c
				case 1:
					assert(false, "immediate mode output param")
				case 2:
					cv = c + relBase
				}
				ensureInbounds(mem, cv)

				fmt.Fprintf(&log, "equals(%d,%d,%d): %d==%d -> %d", a, b, c, av, bv, cv)
				if av == bv {
					mem[cv] = 1
				} else {
					mem[cv] = 0
				}
			case 9: // adjust relative parameter base
				a := chomp(mem, &pc)
				av := paramValue(mem, a, relBase, modes.getNext())

				fmt.Fprintf(&log, "adjust-rel(%d): %d += %d", a, relBase, av)
				relBase += av
			case 99: // halt
				halt = true
				fmt.Fprintf(&log, "halt()")
				doneCh <- struct{}{}
			default:
				panic(fmt.Errorf("Bad opcode %d at mem[%d]", mem[pc], pc))
			}
			fmt.Println(log.String())
		}
	}()

	return outCh, doneCh
}

// MemSize .
const MemSize = 5000

func dupmem(mem []int) []int {
	res := make([]int, MemSize)
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
