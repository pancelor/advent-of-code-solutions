package main

import (
	"fmt"
	"strings"

	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
)

var assert = helpers.Assert
var check = helpers.Check

///////////////////////////////////////////////////////////////////////////////
//
// General Program Setup
//
///////////////////////////////////////////////////////////////////////////////

// NTunnels is the number of tunnels in the maze.
const NTunnels = 56

func main() {
	fmt.Printf("answer:\n%v\n", solve(getInput()))
}

var tunnelNames []string

func getInput() *Maze {
	lines, err := helpers.GetLines()
	check(err)

	l0, lines := lines[0], lines[1:]
	tunnelNames = strings.Split(l0, " ")

	var maze Maze
	for _, l := range lines {
		if l == "" {
			continue
		}
		var row []Tile
		for _, b := range []byte(l) {
			row = append(row, parseTile(b))
		}
		maze.tiles = append(maze.tiles, row)
	}

	assert(nextTunnelID == NTunnels, "didn't use all tunnels")
	return &maze
}

///////////////////////////////////////////////////////////////////////////////
//
// Basic Types
//
///////////////////////////////////////////////////////////////////////////////

// TileType .
type TileType int

const (
	TT_UNKNOWN TileType = iota
	TT_EMPTY
	TT_WALL
	TT_TUNNEL
	TT_OOB
)

// Tile .
type Tile struct {
	tag TileType
	tid TunnelID // optional; exists iff tag == TT_TUNNEL
}

var nextTunnelID = TunnelID(0)

func parseTile(b byte) Tile {
	switch b {
	case '.':
		return Tile{tag: TT_EMPTY}
	case '#':
		return Tile{tag: TT_WALL}
	case ' ':
		return Tile{tag: TT_OOB}
	case 'x':
		tile := Tile{tag: TT_TUNNEL, tid: nextTunnelID}
		nextTunnelID++
		return tile
	default:
		assert(false, "bad tile %s", b)
		return Tile{}
	}
}

func (t Tile) String() string {
	switch t.tag {
	case TT_EMPTY:
		return "."
	case TT_WALL:
		return "#"
	case TT_OOB:
		return " "
	case TT_TUNNEL:
		return t.tid.String()
	default:
		return "?"
	}
}

type point struct {
	x, y int
}

func (p point) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
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

// TunnelID .
type TunnelID int

func (tid TunnelID) String() string {
	assert(helpers.Inbounds(int(tid), 0, len(tunnelNames)), "bad tunnel id")
	return tunnelNames[tid]
}

// Maze represents the static, physical structure of the maze.
// No dynamic solution state is included in this type
type Maze struct {
	tiles [][]Tile
}

func (maze Maze) String() string {
	var s strings.Builder
	for _, row := range maze.tiles {
		for _, tile := range row {
			fmt.Fprintf(&s, "%s", tile)
		}
		fmt.Fprintf(&s, "\n")
	}
	return s.String()
}

func (maze *Maze) inbounds(p point) bool {
	return helpers.Inbounds(p.x, 0, len(maze.tiles[0])) && helpers.Inbounds(p.y, 0, len(maze.tiles))
}

// tunnelAt returns the TunnelID of the tunnel at the given position,
// or panics if no such tunnel exists
func (maze *Maze) tunnelAt(p point) TunnelID {
	tile := maze.tiles[p.y][p.x]
	assert(tile.tag == TT_TUNNEL, "tile at %s isn't a tunnel; it's %v", p, tile)
	return tile.tid
}

// getTid turns a name into a tid
func getTid(name string, n int) TunnelID {
	assert(n == 1 || n == 2, "bad n")
	for i, v := range tunnelNames {
		if v == name {
			n--
		}
		if n == 0 {
			return TunnelID(i)
		}
	}
	assert(false, "no tunnel found with name %s", name)
	return TunnelID(0)
}

// tLoc returns the location of the given TunnelID,
// or panics if no such tunnel exists
func (maze *Maze) tLoc(tid TunnelID) point {
	for r, row := range maze.tiles {
		for c, tile := range row {
			if tile.tag == TT_TUNNEL && tile.tid == tid {
				return point{x: c, y: r}
			}
		}
	}
	assert(false, "TunnelID %s is noplace", tid)
	return point{}
}

///////////////////////////////////////////////////////////////////////////////
//
// TunnelDistances computation
//
///////////////////////////////////////////////////////////////////////////////

// TunnelDistances keeps track of no-tunnel distances between tunnels.
type TunnelDistances map[TunnelID]map[TunnelID]int

func (tds TunnelDistances) String() string {
	var s strings.Builder
	for tidA := TunnelID(0); tidA < NTunnels; tidA++ {
		for tidB := TunnelID(0); tidB < NTunnels; tidB++ {
			if dist, prs := tds[tidA][tidB]; prs {
				fmt.Fprintf(&s, "(%s(%d)->%s(%d):%d), ", tidA, tidA, tidB, tidB, dist)
			}
		}
	}
	return s.String()
}

// assumes the maze has no cycles - the whole algorithm falls apart otherwise
func (maze *Maze) precomputeTunnelDistances() TunnelDistances {
	res := make(map[TunnelID]map[TunnelID]int)
	for tid := TunnelID(0); tid < NTunnels; tid++ {
		fmt.Printf("Precalculating tid %d\n", tid)
		res[tid] = maze.precomputeTunnelDistancesFrom(tid)
		pID, ok := tid.pairID()
		if ok {
			res[tid][pID] = 1
		}
	}
	return res
}

func (tid TunnelID) pairID() (TunnelID, bool) {
	s := tid.String()
	if s == "AA" || s == "ZZ" {
		return TunnelID(0), false
	}
	t1 := getTid(s, 1)
	t2 := getTid(s, 2)
	if tid == t1 {
		return t2, true
	} else if tid == t2 {
		return t1, true
	} else {
		assert(false, "no pair")
		return TunnelID(0), false
	}
}

type VisitType int

const (
	VT_UNVISITED VisitType = iota
	VT_VISITED
	VT_FRONTEIR
)

type VisitTracker map[point]VisitInfo

type VisitInfo struct {
	tag  VisitType
	dist int
}

func makeVisitTracker() VisitTracker {
	return make(map[point]VisitInfo)
}

func (vt VisitTracker) nextToVisit() (point, VisitInfo, bool) {
	for p, info := range vt {
		if info.tag == VT_FRONTEIR {
			return p, info, true
		}
	}
	return point{}, VisitInfo{}, false
}

func (vt VisitTracker) numLeft() int {
	count := 0
	for _, info := range vt {
		if info.tag == VT_FRONTEIR {
			count++
		}
	}
	return count
}

func (maze *Maze) precomputeTunnelDistancesFrom(tid TunnelID) map[TunnelID]int {
	res := make(map[TunnelID]int)
	// if tid.unused() {
	// 	return res
	// }

	visited := makeVisitTracker()
	p0 := maze.tLoc(tid)
	visited[p0] = VisitInfo{tag: VT_FRONTEIR}
	for {
		p, info, ok := visited.nextToVisit()
		if !ok {
			break
		}
		// fmt.Println(p)
		dist := info.dist
		visited[p] = VisitInfo{tag: VT_VISITED, dist: dist}
		if !maze.inbounds(p) {
			continue
		}

		// fmt.Println("num left:", visited.numLeft())
		// fmt.Printf("%s, %v, %d\n", p, info, dist)
		tile := maze.tiles[p.y][p.x]
		switch tile.tag {
		case TT_EMPTY:
			// nothing special
		case TT_WALL:
			continue // don't add neighbors
		case TT_OOB:
			continue
		case TT_TUNNEL:
			if tile.tid == tid {
				assert(dist == 0, "loop?") // note: this is not an exhaustive loop check
			} else {
				res[tile.tid] = dist
			}
		default:
			assert(false, "unknown tile.tag %s at %s", tile.tag, p)
		}

		for _, np := range p.neighbors() {
			if visited[np].tag == VT_UNVISITED {
				visited[np] = VisitInfo{tag: VT_FRONTEIR, dist: dist + 1}
			}
		}
	}

	return res
}

///////////////////////////////////////////////////////////////////////////////
//
// Maze solving
//
///////////////////////////////////////////////////////////////////////////////

// solve returns the minimum number of steps it takes to solve the maze
func solve(maze *Maze) int {
	fmt.Printf("maze:\n%s\n", maze)
	return maze.precomputeTunnelDistances().solve()
}

type TTD struct {
	t1   TunnelID
	t2   TunnelID
	dist int
}

func (tds TunnelDistances) update(updates []TTD) TunnelDistances {
	for _, ttd := range updates {
		t1 := ttd.t1
		t2 := ttd.t2
		dist := ttd.dist
		// fmt.Println("update", t1, t2, dist)
		assert(tds[t1][t2] == tds[t2][t1], "tds is asymmetric")
		tds[t1][t2] = dist
		tds[t2][t1] = dist
	}
	return tds
}

func (tds TunnelDistances) solve() int {
	fmt.Printf("tds:\n%s\n", tds)
	tidAA := getTid("AA", 1)
	tidZZ := getTid("ZZ", 1)
	for {
		var updates []TTD
		for t1, submap := range tds {
			for t2, d12 := range submap {
				for t3, d23 := range tds[t2] {
					if t1 == t3 {
						continue
					}
					dist := d12 + d23 + 1
					if t1 == tidAA {
						fmt.Println(t1, t2, t3, dist, "<?", tds[t1][t3])
					}
					if oldDist, prs := tds[t1][t3]; !prs || dist < oldDist {
						updates = append(updates, TTD{t1, t3, dist})
					}
				}
			}
		}
		if len(updates) == 0 {
			break
		}
		tds = tds.update(updates)
	}

	return tds[tidAA][tidZZ]
}

// func (tds TunnelDistances) distsTo(tid TunnelID) map[TunnelID]int {
// 	for k, v :=
// }
