package main

import (
	"bufio"
	debug "log"
	"io"
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	debug.SetFlags(0)
	debug.SetPrefix("debug: ")

	answer, err := solve(NewValReader(os.Stdin))
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}

func solve(reader valReader) (answer string, outErr error) {
	seen := make(map[string]bool)
	for v := range reader.Vals() {
		mutations := allMutations(v)
		// debug.Println(mutations)
		for _, m := range mutations {
			if seen[m] {
				answer = m
				return
			}
		}
		for _, m := range mutations {
			seen[m] = true
		}
	}
	if err := reader.Err(); err != nil {
		outErr = err
	}
	return
}

func allMutations(s string) []string {
	bytes := []byte(s)
	mutations := make([]string, len(bytes))
	for i := 0; i < len(bytes); i++ {
		// splice out character i
		spliced := make([]byte, len(bytes)-1)
		for j := 0; j < len(spliced); j++ {
			offset := 0
			if j >= i {
				offset = 1
			}
			spliced[j] = bytes[j + offset]
		}
		mutations[i] = string(spliced)
	}
	return mutations
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

func NewValReader(in io.Reader) valReader {
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
