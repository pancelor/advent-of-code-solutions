package main

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
	"log"
	"flag"
)

var infile = flag.String("in", "in.txt", "input file")

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("debug: ")

	reader := NewValReader(*infile)

	seen := make(map[int]bool)
	seen[0] = true
	cumsums := make([]int, 0)

	var total int
	for val := range reader.Vals() {
		total += val
		if seen[total] {
			fmt.Println(total)
			return
		}
		seen[total] = true
		cumsums = append(cumsums, total)
	}
	check(reader.Err())
	// log.Println(cumsums)

	for it := 1; true; it++ {
		log.Println("it", it)
		for _, val := range cumsums {
			laterVal := val + it*total
			// log.Printf("val %v, lval %v\n", val, laterVal)
			if seen[laterVal] {
				fmt.Println(laterVal)
				return
			}
		}
	}
}

type valReader struct {
	filename string
	c chan int
	err error
}

func NewValReader(filename string) valReader {
	c := make(chan int)
	return valReader{
		filename: filename,
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

func (r *valReader) Vals() chan int {
	go func() {
		stream, err := os.Open(r.filename)
		if r.check(err) {
			return
		}

		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			val, err := strconv.Atoi(scanner.Text())
			if r.check(err) {
				return
			}
			r.c <- val
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
