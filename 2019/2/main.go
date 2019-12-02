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
	targetRes := 19690720
	for noun := 0; noun < len(memTemplate); noun++ {
		for verb := 0; verb < len(memTemplate); verb++ {
			mem := dupmem(memTemplate)
			mem[1] = noun
			mem[2] = verb

			res, err := run(mem)
			if err != nil {
				continue
			}
			if res == targetRes {
				return 100*noun + verb, nil
			}
		}
	}
	return nil, fmt.Errorf("No combination found")
}

func run(mem []int) (int, error) {
	pc := -1
	var halt bool
	for !halt {
		opcode := chomp(mem, &pc)
		// dump(mem, pc)
		switch opcode {
		case 1:
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			c := chomp(mem, &pc)
			ensureInbounds(mem, a, b, c)
			av := mem[a]
			bv := mem[b]
			mem[c] = av + bv
		case 2:
			a := chomp(mem, &pc)
			b := chomp(mem, &pc)
			c := chomp(mem, &pc)
			ensureInbounds(mem, a, b, c)
			av := mem[a]
			bv := mem[b]
			mem[c] = av * bv
		case 99:
			halt = true
		default:
			return 0, fmt.Errorf("Bad opcode %d at mem[%d]", mem[pc], pc)
		}
	}

	return mem[0], nil
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

func main() {
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
