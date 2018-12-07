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

func solve(reader valReader) (answer string, outErr error) {
	lines := make([]string, 0)
	for v := range reader.Vals() {
		lines = append(lines, v)
		// debug.Println(v)
	}
	sort.Strings(lines)

	// sleepPatterns[guard_id][min] is the number of times
	//   that guard was asleep at that minute
	sleepPatterns := make(map[int][]int, 60)
	for l := range night.LogGenerator(lines) {
		debug.Printf("%#v\n", l)
		_ = sleepPatterns
		// sleepPatterns[currentGuard][] // TODO
	}

	answer = "unimplemented"
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
