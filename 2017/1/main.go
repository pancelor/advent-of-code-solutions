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
	for i := 0; i < len(v); i++ {
		curr := v[i]
		var prev byte
		if (i == 0) {
			prev = v[len(v)-1]
		} else {
			prev = v[i-1]
		}
		// debug.Println(i, curr, prev)
		if curr == prev {
			res += int(curr - '0')
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
