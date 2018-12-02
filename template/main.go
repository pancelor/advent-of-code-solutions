package main

import (
	// "fmt"
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

	for v := range reader.Vals() {
		log.Println(v)
	}
	check(reader.Err())
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
