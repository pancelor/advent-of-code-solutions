package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("debug: ")
	// reader, err := os.Open("./in.txt")
	// if err != nil {gu
	// 	panic(err)
	// }
	data, err := getInput(os.Stdin)
	if err != nil {
		panic(err)
	}
	outf, err := os.Create("./data.js")
	if err != nil {
		panic(err)
	}
	err = dump(data, outf)
	if err != nil {
		panic(err)
	}
}

func getInput(in io.Reader) (data []particle, err error) {
	data = make([]particle, 0)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		text := scanner.Text()
		p := particle{}
		_, err = fmt.Sscanf(text, "position=<%d,%d> velocity=<%d,%d>", &p.X, &p.Y, &p.Dx, &p.Dy)
		if err != nil {
			return
		}
		data = append(data, p)
	}
	err = scanner.Err()
	return
}

func dump(data []particle, f *os.File) error {
	// for i, _ := range data {
	// 	data[i].X += 2
	// }
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = f.WriteString("window.data = ")
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	return err
}

type command struct {
	dt int
}

type particle struct {
	X, Y   int
	Dx, Dy int
}

// func (p *particle) Advance() {
// 	p.x += p.dx
// 	p.y += p.dy
// }

func (p *particle) String() string {
	return fmt.Sprintf("{x=(%3d,%3d),v=(%3d,%3d)}", p.X, p.Y, p.Dx, p.Dy)
}
