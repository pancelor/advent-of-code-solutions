package main

import (
	"fmt"
	"io"
	debug "log"
	"os"
)

const (
	MAX_X = 300
	MAX_Y = 300
	// SN, XTAR, YTAR = 18, 33, 45
	// SN, XTAR, YTAR = 42, 21, 61
)

type grid [MAX_Y][MAX_X]int

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	debug.SetFlags(0)
	debug.SetPrefix("debug: ")

	data, err := getInput(os.Stdin)
	if err != nil {
		panic(err)
	}
	answer, err := solve(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}

func getInput(in io.Reader) (sn int, err error) {
	_, err = fmt.Scanf("%d", &sn)
	// return SN, nil
	return
}

func hundredsDigit(x int) int {
	y := (x / 100) % 10
	return y
}

func solve(sn int) (answer string, err error) {
	var g grid
	for r := 0; r < MAX_Y; r++ {
		for c := 0; c < MAX_X; c++ {
			x := c + 1
			y := r + 1
			// inefficient raw translation from pseudocode; we could precompute a LOT of this...
			rackId := x + 10
			power := rackId * y
			power += sn
			power *= rackId
			power = hundredsDigit(power)
			power -= 5
			g[r][c] = power
		}
	}

	rmax := 0
	cmax := 0
	max := g.val3by3at(0, 0)
	for r := 0; r < MAX_Y; r++ {
		for c := 0; c < MAX_X; c++ {
			val := g.val3by3at(r, c)
			if val > max {
				max = val
				rmax = r
				cmax = c
			}
		}
	}

	answer = fmt.Sprintf("%d,%d", cmax+1, rmax+1)
	return
}

func (g *grid) val3by3at(r, c int) (val int) {
	// inefficient since we can incrementally compute this as we go
	for dr := 0; dr < 3; dr++ {
		for dc := 0; dc < 3; dc++ {
			val += g.at(r+dr, c+dc)
		}
	}
	return
}

func (g *grid) at(r, c int) int {
	if 0 <= r && 0 <= c && r < MAX_Y && c < MAX_X {
		return g[r][c]
	} else {
		return 0
	}
}

func (g *grid) print(rmin, cmin, rmax, cmax int) {
	fmt.Print("[")
	h := len(g)
	if h > 0 {
		w := len(g[0])
		if w < cmax || h < rmax {
			panic("bad bounds in grid.print")
		}
		for r := rmin; r < rmax; r++ {
			if r != 0 {
				fmt.Print("\n ")
			}
			fmt.Print("[")
			for c := cmin; c < cmax; c++ {
				if c != 0 {
					fmt.Print(" ")
				}
				fmt.Print(g[r][c])
			}
			fmt.Print("]")
		}
	}
	fmt.Print("]\n")
}

func (g *grid) print5by5WithTopLeft(x, y int) {
	r := y - 1
	c := x - 1
	g.print(r, c, r+5, c+5)
}
