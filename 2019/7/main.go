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
	[5]int{5, 6, 7, 8, 9},
	[5]int{5, 6, 7, 9, 8},
	[5]int{5, 6, 8, 7, 9},
	[5]int{5, 6, 8, 9, 7},
	[5]int{5, 6, 9, 7, 8},
	[5]int{5, 6, 9, 8, 7},
	[5]int{5, 7, 6, 8, 9},
	[5]int{5, 7, 6, 9, 8},
	[5]int{5, 7, 8, 6, 9},
	[5]int{5, 7, 8, 9, 6},
	[5]int{5, 7, 9, 6, 8},
	[5]int{5, 7, 9, 8, 6},
	[5]int{5, 8, 6, 7, 9},
	[5]int{5, 8, 6, 9, 7},
	[5]int{5, 8, 7, 6, 9},
	[5]int{5, 8, 7, 9, 6},
	[5]int{5, 8, 9, 6, 7},
	[5]int{5, 8, 9, 7, 6},
	[5]int{5, 9, 6, 7, 8},
	[5]int{5, 9, 6, 8, 7},
	[5]int{5, 9, 7, 6, 8},
	[5]int{5, 9, 7, 8, 6},
	[5]int{5, 9, 8, 6, 7},
	[5]int{5, 9, 8, 7, 6},
	[5]int{6, 5, 7, 8, 9},
	[5]int{6, 5, 7, 9, 8},
	[5]int{6, 5, 8, 7, 9},
	[5]int{6, 5, 8, 9, 7},
	[5]int{6, 5, 9, 7, 8},
	[5]int{6, 5, 9, 8, 7},
	[5]int{6, 7, 5, 8, 9},
	[5]int{6, 7, 5, 9, 8},
	[5]int{6, 7, 8, 5, 9},
	[5]int{6, 7, 8, 9, 5},
	[5]int{6, 7, 9, 5, 8},
	[5]int{6, 7, 9, 8, 5},
	[5]int{6, 8, 5, 7, 9},
	[5]int{6, 8, 5, 9, 7},
	[5]int{6, 8, 7, 5, 9},
	[5]int{6, 8, 7, 9, 5},
	[5]int{6, 8, 9, 5, 7},
	[5]int{6, 8, 9, 7, 5},
	[5]int{6, 9, 5, 7, 8},
	[5]int{6, 9, 5, 8, 7},
	[5]int{6, 9, 7, 5, 8},
	[5]int{6, 9, 7, 8, 5},
	[5]int{6, 9, 8, 5, 7},
	[5]int{6, 9, 8, 7, 5},
	[5]int{7, 5, 6, 8, 9},
	[5]int{7, 5, 6, 9, 8},
	[5]int{7, 5, 8, 6, 9},
	[5]int{7, 5, 8, 9, 6},
	[5]int{7, 5, 9, 6, 8},
	[5]int{7, 5, 9, 8, 6},
	[5]int{7, 6, 5, 8, 9},
	[5]int{7, 6, 5, 9, 8},
	[5]int{7, 6, 8, 5, 9},
	[5]int{7, 6, 8, 9, 5},
	[5]int{7, 6, 9, 5, 8},
	[5]int{7, 6, 9, 8, 5},
	[5]int{7, 8, 5, 6, 9},
	[5]int{7, 8, 5, 9, 6},
	[5]int{7, 8, 6, 5, 9},
	[5]int{7, 8, 6, 9, 5},
	[5]int{7, 8, 9, 5, 6},
	[5]int{7, 8, 9, 6, 5},
	[5]int{7, 9, 5, 6, 8},
	[5]int{7, 9, 5, 8, 6},
	[5]int{7, 9, 6, 5, 8},
	[5]int{7, 9, 6, 8, 5},
	[5]int{7, 9, 8, 5, 6},
	[5]int{7, 9, 8, 6, 5},
	[5]int{8, 5, 6, 7, 9},
	[5]int{8, 5, 6, 9, 7},
	[5]int{8, 5, 7, 6, 9},
	[5]int{8, 5, 7, 9, 6},
	[5]int{8, 5, 9, 6, 7},
	[5]int{8, 5, 9, 7, 6},
	[5]int{8, 6, 5, 7, 9},
	[5]int{8, 6, 5, 9, 7},
	[5]int{8, 6, 7, 5, 9},
	[5]int{8, 6, 7, 9, 5},
	[5]int{8, 6, 9, 5, 7},
	[5]int{8, 6, 9, 7, 5},
	[5]int{8, 7, 5, 6, 9},
	[5]int{8, 7, 5, 9, 6},
	[5]int{8, 7, 6, 5, 9},
	[5]int{8, 7, 6, 9, 5},
	[5]int{8, 7, 9, 5, 6},
	[5]int{8, 7, 9, 6, 5},
	[5]int{8, 9, 5, 6, 7},
	[5]int{8, 9, 5, 7, 6},
	[5]int{8, 9, 6, 5, 7},
	[5]int{8, 9, 6, 7, 5},
	[5]int{8, 9, 7, 5, 6},
	[5]int{8, 9, 7, 6, 5},
	[5]int{9, 5, 6, 7, 8},
	[5]int{9, 5, 6, 8, 7},
	[5]int{9, 5, 7, 6, 8},
	[5]int{9, 5, 7, 8, 6},
	[5]int{9, 5, 8, 6, 7},
	[5]int{9, 5, 8, 7, 6},
	[5]int{9, 6, 5, 7, 8},
	[5]int{9, 6, 5, 8, 7},
	[5]int{9, 6, 7, 5, 8},
	[5]int{9, 6, 7, 8, 5},
	[5]int{9, 6, 8, 5, 7},
	[5]int{9, 6, 8, 7, 5},
	[5]int{9, 7, 5, 6, 8},
	[5]int{9, 7, 5, 8, 6},
	[5]int{9, 7, 6, 5, 8},
	[5]int{9, 7, 6, 8, 5},
	[5]int{9, 7, 8, 5, 6},
	[5]int{9, 7, 8, 6, 5},
	[5]int{9, 8, 5, 6, 7},
	[5]int{9, 8, 5, 7, 6},
	[5]int{9, 8, 6, 5, 7},
	[5]int{9, 8, 6, 7, 5},
	[5]int{9, 8, 7, 5, 6},
	[5]int{9, 8, 7, 6, 5},
}

func solve(memTemplate []int) (interface{}, error) {
	best := 0
	for _, settings := range allPerms {
		assert(len(settings) == 5, "bad settings len")

		// TEMP
		// TEMP
		// TEMP
		settings = [5]int{9, 7, 8, 5, 6}
		// TEMP
		// TEMP
		// TEMP

		fmt.Printf("\n\n\nsettings=%v\n", settings)

		startCh := make(chan int)

		aCh, aDone := run("A", dupmem(memTemplate), startCh)
		bCh, bDone := run("B", dupmem(memTemplate), aCh)
		cCh, cDone := run("C", dupmem(memTemplate), bCh)
		dCh, dDone := run("D", dupmem(memTemplate), cCh)
		eCh, eDone := run("E", dupmem(memTemplate), dCh)

		_ = aDone
		_ = bDone
		_ = cDone
		_ = dDone

		// prime the pump
		startCh <- settings[0]
		aCh <- settings[1]
		bCh <- settings[2]
		cCh <- settings[3]
		dCh <- settings[4]

		// go
		startCh <- 0

		go func() {
			for x := range eCh {
				// fmt.Printf("---got %d from eCh---\n", x)
				startCh <- x
			}
		}()

		<-eDone
		res := <-startCh
		fmt.Printf("(settings=%v) res=%v\n", settings, res)
		if res > best {
			best = res
		}

		// TEMP
		// TEMP
		// TEMP
		break
		// TEMP
		// TEMP
		// TEMP
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

func run(name string, mem []int, inCh chan int) (chan int, chan struct{}) {
	outCh := make(chan int)
	doneCh := make(chan struct{})
	go func() {
		pc := -1
		var halt bool
		for !halt {
			code := chomp(mem, &pc)
			opcode, modes, err := parseOpcode(code)
			check(err)

			// fmt.Printf("node %s pc=%v\n", name, pc)
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
				i := <-inCh
				fmt.Printf("%s < %d\n", name, i)

				mem[a] = i
			case 4: // output
				a := chomp(mem, &pc)
				av := paramValue(mem, a, modes.getNext())

				fmt.Printf("%s : %d\n", name, av)
				outCh <- av
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
				doneCh <- struct{}{}
			default:
				panic(fmt.Errorf("Bad opcode %d at mem[%d]", mem[pc], pc))
			}
		}
	}()

	return outCh, doneCh
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
