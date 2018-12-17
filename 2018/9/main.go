package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// const instring = "13 players; last marble is worth 7999 points"

func main() {
	nPlayers, nMarbles, err := getInput(os.Stdin)
	// nPlayers, nMarbles, err := getInput(strings.NewReader(instring))
	if err != nil {
		panic(err)
	}
	answer, err := solve(nPlayers, nMarbles)
	fmt.Println(answer)
}

func getInput(in io.Reader) (nPlayers int, nMarbles int, err error) {
	_, err = fmt.Fscanf(in, "%d players; last marble is worth %d points", &nPlayers, &nMarbles)
	return
}

func solve(nPlayers int, nMarbles int) (answer int, err error) {
	score := make([]int, nMarbles)
	g := gameState{
		ring: make([]int, 0),
	}
	for i := 0; i < nMarbles; i++ {
		// log.Println("i", i)
		g.p = i % nPlayers // note that this is actually off-by-one
		if i == 0 {
			g.ring = insert(g.ring, 0, i)
		} else if i%23 != 0 {
			// log.Println("23! i1", g.i)
			g.changeIndex(2)
			if g.i == 0 {
				g.i = len(g.ring)
			}
			// log.Println("23! i2", g.i)
			g.ring = insert(g.ring, g.i, i)
		} else {
			g.changeIndex(-7)
			score[g.p] += i + g.ring[g.i]
			g.ring = arrayDelete(g.ring, g.i)
		}
		// log.Println(g)
	}
	answer = max(score)
	return
}

type gameState struct {
	ring []int
	p    int // current player
	i    int // current index in ring
}

func (g *gameState) changeIndex(delta int) {
	g.i += delta
	for g.i < 0 {
		g.i += len(g.ring)
	}
	for g.i >= len(g.ring) {
		g.i -= len(g.ring)
	}
}

func (g *gameState) copy() gameState {
	c := gameState{
		p: g.p,
		i: g.i,
	}
	copy(c.ring, g.ring)
	return c
}

func (g gameState) String() string {
	var s strings.Builder
	fmt.Fprintf(&s, "[%d]", g.p+1) // off-by-one in RAM; corrected here
	for i := 0; i < len(g.ring); i++ {
		if i == g.i {
			fmt.Fprintf(&s, " (%3d)", g.ring[i])
		} else {
			fmt.Fprintf(&s, "  %3d ", g.ring[i])
		}
	}
	return s.String()
}

func insert(arr []int, i int, x int) []int {
	// log.Printf("insert(len=%d, i=%d)", len(arr), i)
	arr = append(arr, 0)
	copy(arr[i+1:], arr[i:])
	arr[i] = x
	return arr
}

func arrayDelete(arr []int, i int) []int {
	// log.Printf("arrayDelete(len=%d, i=%d)", len(arr), i)
	return append(arr[:i], arr[i+1:]...)
}

func max(arr []int) int {
	if len(arr) == 0 {
		panic("max([]int{}) is undefined")
	}
	res := 0
	for _, x := range arr {
		if x > res {
			res = x
		}
	}
	return res
}
