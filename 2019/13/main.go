package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pancelor/advent-of-code-solutions/2019/computer"
	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
)

var assert = helpers.Assert
var check = helpers.Check

// Input is the type of input we pass to solve()
type Input []int

type TileType int

const (
	TT_EMPTY      TileType = iota // No game object appears in this tile.
	TT_WALL                       // Walls are indestructible barriers.
	TT_BLOCK                      // Blocks can be broken by the ball.
	TT_HORIZONTAL                 // The paddle is indestructible.
	TT_BALL                       // The ball moves diagonally and bounces off objects.
)

func (t *TileType) String() string {
	switch *t {
	case TT_EMPTY:
		return "."
	case TT_WALL:
		return "#"
	case TT_BLOCK:
		return "O"
	case TT_HORIZONTAL:
		return "_"
	case TT_BALL:
		return "*"
	}
	return "?"
}

type Screen [20][38]TileType

func (s *Screen) sizeCheck(x, y int) {
	assert(helpers.Inbounds(x, 0, 400) && helpers.Inbounds(y, 0, 400), "bad screen coords (%d, %d)", x, y)
}

func (s *Screen) set(x, y int, t TileType) {
	s.sizeCheck(x, y)
	s[y][x] = t
}

func (s *Screen) get(x, y int) TileType {
	s.sizeCheck(x, y)
	return s[y][x]
}

func (s *Screen) String() string {
	var b strings.Builder
	for y := 0; y < len(s); y++ {
		for x := 0; x < len(s[0]); x++ {
			fmt.Fprintf(&b, s[y][x].String())
		}
		fmt.Fprintf(&b, "\n")
	}
	return b.String()
}

func (s *Screen) Count(t TileType) int {
	total := 0
	for y := 0; y < len(s); y++ {
		for x := 0; x < len(s[0]); x++ {
			if s[y][x] == t {
				total++
			}
		}
	}
	return total
}

func solve(in Input) interface{} {
	cpu := computer.MakeCPU("grenadier")
	in[0] = 2
	cpu.SetMemory(in)

	// fmt.Println(cpu.PrintProgram())
	// return nil

	cpu.Run()
	score := 0

	var screen Screen
	doPrint := false
	ballX := 0
	paddleX := 0
	for {
		switch state := <-cpu.StateChan; state {
		case computer.CS_WAITING_INPUT:
			// doPrint = true
			// time.Sleep(1 * time.Millisecond)
			val := 0
			if ballX < paddleX {
				val = -1
			} else if ballX > paddleX {
				val = 1
			}
			cpu.InChan <- val
		case computer.CS_WAITING_OUTPUT:
			x := <-cpu.OutChan
			y := cpu.RecvOutput()
			z := cpu.RecvOutput()
			if x == -1 && y == 0 {
				score = z
			} else {
				t := TileType(z)
				// fmt.Printf("x,y,t=%d,%d,%s\n", x, y, t.String())
				screen.set(x, y, t)
				if t == TT_BALL {
					ballX = x
				}
				if t == TT_HORIZONTAL {
					paddleX = x
				}
			}
			if doPrint {
				fmt.Printf("score: %d\n", score)
				fmt.Printf("%s\n", screen.String())
			}
		case computer.CS_DONE:
			return score
		default:
			assert(false, "unexpected state %s", state)
		}
	}
}

func init() {
	// tests go here
	assert(true, "t1")

	// assert(false, "exit after tests")
}

func main() {
	input, err := getInput()
	check(err)
	answer := solve(input)
	fmt.Printf("answer:\n%v\n", answer)
}

func getInput() (Input, error) {
	lines, err := helpers.GetLines()
	if err != nil {
		return nil, err
	}

	var res []int
	for _, l := range strings.Split(lines[0], ",") {
		if l == "" {
			continue
		}

		v, err := strconv.Atoi(l)
		if err != nil {
			return nil, err
		}

		// var dirCode string
		// var dist int
		// _, err := fmt.Sscanf(token, "%1s%d", &dirCode, &dist)
		// if err != nil {
		// 	return res, err
		// }

		res = append(res, v)
	}

	return res, nil
}
