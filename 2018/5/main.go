package main

import (
	"bufio"
	debug "log"
	"io"
	"fmt"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	line string
)

func main() {
	debug.SetFlags(0)
	debug.SetPrefix("debug: ")

	answer, err := solve(newValReader(os.Stdin))
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}

func newMask() []int {
	mask := make([]int, len(line))
	for i := 0; i < len(mask); i++ {
		mask[i] = 1
	}
	return mask
}

func debugLine() {
	fmt.Printf("       [")
	for i, e := range line {
		if i != 0 {
			fmt.Printf(" ")
		}
		fmt.Printf("%c", e)
	}
	fmt.Printf("]\n")
}

func solve(reader valReader) (answer int, outErr error) {
	line = <-reader.Vals()

	// debugLine()

	mask := newMask()
	for annihilateOnePass(mask) {}
	answer = sum(mask)
	debug.Printf("base string: %d", answer)
	for c := 'a'; c <= 'z'; c++ {
		mask := newMask()
		killPair(c, mask)
		for annihilateOnePass(mask) {}
		res := sum(mask)
		if res < answer {
			answer = res
		}
		debug.Printf("c %c; res %d; answer %d", c, res, answer)
	}

	return
}

func killPair(toKill rune, mask []int) {
	for i, c := range(strings.ToLower(line)) {
		if c == toKill {
			mask[i] = 0
		}
	}
}

func annihilateOnePass(mask []int) bool {
	// debug.Println(mask)
	prevI := -1
	changes := false
	for i := 0; i < len(mask); i++ {
		if mask[i] == 1 {
			if prevI == -1 {
				prevI = i
			} else {
				if shouldAnnihilate(prevI, i) {
					changes = true
					mask[prevI] = 0
					mask[i] = 0
					prevI = -1
				} else {
					prevI = i
				}
			}
		} else if mask[i] == 0 {

		} else {
			assert(false)
		}
	}
	return changes
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}

func shouldAnnihilate(i1, i2 int) bool {
	c1 := line[i1]
	c2 := line[i2]
	if c2 < c1 {
		c1, c2 = c2, c1
	}
	// c2 is lowercase, c1 is uppercase
	c1Upper := 'A' <= c1 && c1 <= 'Z'
	c2Lower := 'a' <= c2 && c2 <= 'z'
	diffCorrect := c2-c1 == 'a'-'A'
	// debug.Println(c1, c2, c1Upper, c2Lower, diffCorrect)
	return c1Upper && c2Lower && diffCorrect
}

func sum(a []int) int {
	res := 0
	for _, e := range a {
		res += e
	}
	return res
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
