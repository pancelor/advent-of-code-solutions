package main

import (
	"bufio"
	debug "log"
	"io"
	"fmt"
	"os"
	"./night"
	"sort"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	debug.SetFlags(0)
	debug.SetPrefix("debug: ")

	answer, err := solve(newValReader(os.Stdin))
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}

func solve(reader valReader) (answer int, outErr error) {
	lines := make([]string, 0)
	for v := range reader.Vals() {
		lines = append(lines, v)
		// debug.Println(v)
	}
	sort.Strings(lines)

	// sleepPatterns[guard_id][min] is the number of times
	//   that guard was asleep at that minute
	sleepPatterns := make(map[int][]int)
	for state := range night.StateGenerator(lines) {
		// debug.Printf("%#v\n", state)
		if state.Asleep {
			if sleepPatterns[state.Guard] == nil {
				sleepPatterns[state.Guard] = make([]int, 60)
			}
			sleepPatterns[state.Guard][state.Time] += 1
		}
	}

	debugPatterns(sleepPatterns)

	guard, max := argmax(sum, sleepPatterns)
	minute, _ := argmaxInt(sleepPatterns[guard])

	debug.Printf("Guard #%d slept for %d minutes overall; mostly at time %d\n", guard, max, minute)

	answer = guard*minute
	return
}

func sum(a []int) int {
	res := 0
	for _, e := range a {
		res += e
	}
	return res
}

func debugPatterns(pats map[int][]int) {
	fmt.Println("      ID | Minute")
	fmt.Println("         | 000000000011111111112222222222333333333344444444445555555555")
	fmt.Println("         | 012345678901234567890123456789012345678901234567890123456789")
	fmt.Println("-----------------------------------------------------------------------")
	for id, arr := range pats {
		fmt.Printf(" #%6d | ", id)
		for _, sleepCount := range arr {
			if sleepCount >= 100 {
				fmt.Printf("*")
			} else if sleepCount >= 10 {
				fmt.Printf("x")
			} else {
				fmt.Printf("%d", sleepCount)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Println()
}

func argmax(proj func([]int) int, m map[int][]int) (arg, max int) {
	if len(m) == 0 {
		panic("empty array")
	}
	arg = 0
	max = 0 // hacky; not mathematically correct
	for k, v := range m {
		p := proj(v)
		if proj(v) > max {
			arg = k
			max = p
		}
	}
	if arg == 0 {
		panic("probably need to rewrite argmax to play nice")
	}
	return
}

func argmaxInt(arr []int) (arg, max int) {
	if len(arr) == 0 {
		panic("empty array")
	}
	arg = 0
	max = 0 // hacky; not mathematically correct
	for k, v := range arr {
		if v > max {
			arg = k
			max = v
		}
	}
	if arg == 0 {
		panic("probably need to rewrite argmaxInt to play nice")
	}
	return
}

func maxInt(a []int) (max int) {
	if len(a) == 0 {
		panic("empty array")
	}
	max = 0 // hacky; not mathematically correct
	for _, e := range a {
		if e > max {
			max = e
		}
	}
	if max == 0 {
		panic("probably need to rewrite argmax to play nice")
	}
	return
}



// valReader converts an `io.Reader` to a `chan string`
// usage:
//   reader := newValReader(os.Stdin)
//   for val := reader.Vals() {
//     _ = val
//   }
//   if reader.Err() != nil {
//     panic(err)
//   }
type valReader struct {
	in io.Reader
	out chan string
	err error
}

func newValReader(in io.Reader) valReader {
	out := make(chan string)
	return valReader{
		in: in,
		out: out,
	}
}

func (r *valReader) Vals() chan string {
	go func() {
		scanner := bufio.NewScanner(r.in)
		for scanner.Scan() {
			r.out <- scanner.Text()
		}
		if r.check(scanner.Err()) {
			return
		}
		close(r.out)
	}()
	return r.out
}

func (r *valReader) Err() error {
	return r.err
}

func (r *valReader) check(err error) bool {
	if err != nil {
		r.err = err
		close(r.out)
		return true
	}
	return false
}
