package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("debug: ")
	reader, err := os.Open("2018/10/in.txt")
	if err != nil {
		panic(err)
	}
	data, err := getInput(reader)
	if err != nil {
		panic(err)
	}
	commands := newCommandStream(os.Stdin)
	answer, err := solve(data, commands)
	fmt.Println(answer)
}

func getInput(in io.Reader) (data []particle, err error) {
	data = make([]particle, 0)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		text := scanner.Text()
		p := particle{}
		_, err = fmt.Sscanf(text, "position=<%d,%d> velocity=<%d,%d>", &p.x, &p.y, &p.dx, &p.dy)
		if err != nil {
			return
		}
		data = append(data, p)
	}
	err = scanner.Err()
	return
}

func solve(data []particle, commands chan command) (answer int, err error) {
	for _, p := range data {
		log.Println(p)
	}
	return
}

type command struct {
	dt int
}

type particle struct {
	x, y   int
	dx, dy int
}

func (p *particle) Advance() {
	p.x += p.dx
	p.y += p.dy
}

func (p *particle) String() string {
	return fmt.Sprintf("{x=(%3d,%3d),v=(%3d,%3d)}", p.x, p.y, p.dx, p.dy)
}
