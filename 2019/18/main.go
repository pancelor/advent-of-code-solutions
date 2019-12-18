package main

import (
	"fmt"
	"strings"

	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
)

var assert = helpers.Assert
var check = helpers.Check

type TileType int

const (
	TT_UNKNOWN TileType = iota
	TT_EMPTY
	TT_WALL
	TT_KEY
	TT_LOCK
)

type Tile struct {
	tag   TileType
	label byte // optional; exists iff tag == TT_KEY or TT_LOCK
}

func parseTile(b byte) Tile {
	switch b {
	case '.', '@':
		return Tile{tag: TT_EMPTY}
	case '#':
		return Tile{tag: TT_WALL}
	default:
		if 'a' <= b && b <= 'z' {
			return Tile{tag: TT_KEY, label: b}
		} else {
			assert('A' <= b && b <= 'Z', "unknown char %s", b)
			return Tile{tag: TT_LOCK, label: b}
		}
	}
}

func (t *Tile) String() string {
	switch t.tag {
	case TT_EMPTY:
		return "."
	case TT_WALL:
		return "#"
	case TT_KEY, TT_LOCK:
		return string(t.label)
	default:
		return "?"
	}
}

type point struct {
	x, y int
}

func (p point) addDir(dir int, n int) point {
	newP := p
	switch dir {
	case 0:
		newP.x += n
	case 1:
		newP.y -= n
	case 2:
		newP.x -= n
	case 3:
		newP.y += n
	default:
		assert(false, fmt.Sprintf("bad dir %d", dir))
	}
	return newP
}

func (p point) neighbors() (res []point) {
	for d := 0; d < 4; d++ {
		res = append(res, p.addDir(d, 1))
	}
	return
}

type Maze struct {
	tiles [][]Tile
	start point
}

func (maze *Maze) String() string {
	var s strings.Builder
	for r, row := range maze.tiles {
		for c, tile := range row {
			if maze.start == (point{x: c, y: r}) {
				fmt.Fprintf(&s, "@")
			} else {
				fmt.Fprintf(&s, tile.String())
			}
		}
		fmt.Fprintf(&s, "\n")
	}
	return s.String()
}

type Solver struct {
	keys [26]bool
	pos  point
	*Maze
}

type VisitType int

const (
	VT_UNVISITED VisitType = iota
	VT_FRONTEIR
	VT_VISITED
)

func nextToVisit(tracker map[point]VisitType) (point, bool) {
	for p, vt := range tracker {
		if vt == VT_FRONTEIR {
			return p, true
		}
	}
	return point{}, false
}

func (solver *Solver) canUnlock(t Tile) bool {
	return t.tag == TT_LOCK && solver.keys[t.label-'A']
}

func (solver *Solver) flood() (res []Solver) {
	tracker := make(map[point]VisitType)
	tracker[solver.start] = VT_FRONTEIR
	for {
		next, ok := nextToVisit(tracker)
		if !ok {
			break
		}
		tracker[next] = VT_VISITED

		for _, p := range next.neighbors() {
			tile := solver.tiles[p.y][p.x]
			switch tile.tag {
			case TT_EMPTY:
				tracker[p] = VT_FRONTEIR
			case TT_WALL:
				tracker[p] = VT_VISITED // ignore
			case TT_LOCK:
				if solver.canUnlock(tile) {
					tracker[p] = VT_FRONTEIR
				} else {
					tracker[p] = VT_VISITED // TODO not true
					// working here
				}
			}
		}
	}

	// res = append(res, clone)
	return
}

func solve(maze *Maze) interface{} {
	fmt.Printf("maze:\n%s\n", maze)

	answer := "unimplemented"
	return answer
}

func main() {
	input, err := getInput()
	check(err)
	answer := solve(input)
	fmt.Printf("answer:\n%v\n", answer)
}

func getInput() (*Maze, error) {
	lines, err := helpers.GetLines()
	if err != nil {
		return nil, err
	}

	var maze Maze
	for r, l := range lines {
		if l == "" {
			continue
		}
		var row []Tile
		for c, b := range []byte(l) {
			row = append(row, parseTile(b))
			if b == '@' {
				maze.start = point{x: c, y: r}
			}
		}

		maze.tiles = append(maze.tiles, row)
	}

	return &maze, nil
}
