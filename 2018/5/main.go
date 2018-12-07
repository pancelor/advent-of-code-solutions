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
	line := <-reader.Vals()

	answer = 11298
	for c := 'a'; c <= 'z'; c++ {
		reduced := fpa(annihilateOnce, killPair(c, line))
		res := len(reduced)
		if res < answer {
			answer = res
		}
		debug.Printf("c %c; res %d; answer %d", c, res, answer)
	}

	return
}

func killPair(toKill rune, line string) string {
	res := make([]string, 0)
	last := 0
	for i, c := range(strings.ToLower(line)) {
		if c == toKill {
			res = append(res, line[last:i])
			last = i+1
		}
	}
	res = append(res, line[last:])
	return strings.Join(res, "")
}

func annihilateOnce(s string) string {
	// debug.Println(s)
	for i := 0; i < len(s)-1; i++ {
		if shouldAnnihilate(s[i], s[i+1]) {
			return strings.Join([]string{s[:i], s[i+2:]}, "")
		}
	}
	return s
}

func shouldAnnihilate(c1 byte, c2 byte) bool {
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

// repeatedly calls x = iterate(x) until x stops changing
func fpa(iterate func(string) string, x string) string {
	var last string
	for {
		last = x
		x = iterate(x)
		if last == x {
			return x
		}
	}
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
