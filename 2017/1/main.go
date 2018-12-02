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

	answer, err := solve(newValReader(os.Stdin))
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}

func solve(reader valReader) (res int, outErr error) {
	v := <-reader.Vals()
	res = 0
	halfLen := len(v)/2
	for i := 0; i < halfLen; i++ {
		curr := v[i]
		pair := v[i+halfLen]
		// debug.Printf("i %d: %c %c", i, curr, pair)
		if curr == pair {
			res += 2*int(curr - '0')
		}
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
