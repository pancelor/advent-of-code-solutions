package main

import (
	"bufio"
	debug "log"
	"io"
	"fmt"
	"os"
	"./claim"
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
	claims := make([]claim.Claim, 0)
	for v := range reader.Vals() {
		c := claim.New(v)
		claims = append(claims, c)
	}
	outErr = reader.Err()
	if outErr != nil {
		return
	}

	xMax, yMax := getMaxSize(claims)

	// fmt.Println(claims)
	// fmt.Printf("size: %dx%d\n", yMax, xMax)

	g := newGrid(xMax, yMax)
	// printGrid(g)

	for _, c := range claims {
		addClaimToGrid(c, g)
		// printGrid(g)
	}

	reducer := func(acc, x int) int {
		if (x > 1) {
			return acc + 1
		} else {
			return acc
		}
	}
	numTilesWithOverlap := reduceGrid(reducer, 0, g)

	return numTilesWithOverlap, nil
}

func reduceGrid(reducer func(int, int) int, init int, g grid) int {
	res := init
	rMax := len(g)
	if rMax > 0 {
		cMax := len(g[0])
		for r := 0; r < rMax; r++ {
			for c := 0; c < cMax; c++ {
				res = reducer(res, g[r][c])
			}
		}
	}
	return res
}

func addClaimToGrid(c claim.Claim, g grid) {
	rMax := c.Bottom()
	cMax := c.Right()
	for r := c.Y(); r < rMax; r++ {
		for c := c.X(); c < cMax; c++ {
			g[r][c] += 1
		}
	}
}

type grid [][]int

func newGrid(w, h int) grid {
	g := make([][]int, h)
	for i := 0; i < h; i++ {
		g[i] = make([]int, w)
	}
	return g
}

func printGrid(g grid) {
	fmt.Print("[")
	h := len(g)
	if h > 0 {
		w := len(g[0])
		for r := 0; r < h; r++ {
			if r != 0 {
				fmt.Print("\n ")
			}
			fmt.Print("[")
			for c := 0; c < w; c++ {
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

func getMaxSize(claims []claim.Claim) (xMax, yMax int) {
	xMax = max(fmapClaims(func(c claim.Claim) int { return c.Right() }, claims), 1)
	yMax = max(fmapClaims(func(c claim.Claim) int { return c.Bottom() }, claims), 1)
	return
}

func fmapClaims(proj func(c claim.Claim) int, arr []claim.Claim) []int {
	res := make([]int, len(arr))
	for i, in := range arr {
		res[i] = proj(in)
	}
	return res
}

func max(arr []int, def int) int {
	if len(arr) == 0 {
		return def
	}
	res := arr[0]
	for _, e := range arr {
		if e > res {
			res = e
		}
	}
	return res
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
