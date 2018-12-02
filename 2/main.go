package main

import (
	"fmt"
	"bufio"
	"io"
	"os"
	"log"
	// "flag"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("debug: ")

	reader := NewValReader(os.Stdin)

	hasDouble := 0
	hasTriple := 0
	for v := range reader.Vals() {
		c := newCounter(v)
		if len(reverseLookup(c, 2)) > 0 {
			hasDouble += 1
		}
		if len(reverseLookup(c, 3)) > 0 {
			hasTriple += 1
		}
	}
	check(reader.Err())
	fmt.Println(hasDouble*hasTriple)
}

type counter map[byte]int

func newCounter(s string) (m counter) {
	m = make(map[byte]int)
	for _, b := range []byte(s) {
		m[b] += 1
	}
	return
}

func reverseLookup(c counter, target int) []byte {
	res := make([]byte, 0)
	for key := range c {
		// log.Printf("key %c val %v\n", key, c[key])
		if c[key] == target {
			res = append(res, key)
		}
	}
	return res
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
