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

	memory := make(map[int]bool)
	var total int
	loop: for {
		log.Println("again")
		stream, err := os.Open(*infile)
		check(err)
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			val, err := strconv.Atoi(scanner.Text())
			check(err)
			seen := memory[total]
			// log.Println(val, total, seen)
			if seen {
				fmt.Println(total)
				break loop
			}
			memory[total] = true
			total += val
		}
	}
}
