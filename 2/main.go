package main

import (
	"fmt"
	"bufio"
	"io"
	"os"
	debug "log"
	// "flag"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	debug.SetFlags(0)
	debug.SetPrefix("debug: ")

	reader := NewValReader(os.Stdin)
	seen := make(map[string]bool)
	for v := range reader.Vals() {
		mutations := allMutations(v)
		// debug.Println(mutations)
		for _, m := range mutations {
			if seen[m] {
				fmt.Println(m)
				return
			}
		}
		for _, m := range mutations {
			seen[m] = true
		}
	}
	check(reader.Err())
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



type valReader struct {
	stream io.Reader
	c chan string
	err error
}

func NewValReader(stream io.Reader) valReader {
	c := make(chan string)
	return valReader{
		stream: stream,
		c: c,
	}
}

func (r *valReader) check(err error) bool {
	if err != nil {
		r.err = err
		close(r.c)
		return true
	}
	return false
}

func (r *valReader) Vals() chan string {
	go func() {
		scanner := bufio.NewScanner(r.stream)
		for scanner.Scan() {
			r.c <- scanner.Text()
		}
		if r.check(scanner.Err()) {
			return
		}
		close(r.c)
	}()
	return r.c
}

func (r *valReader) Err() error {
	return r.err
}
