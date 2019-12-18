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
	tag TileType
	kid KeyID // optional; exists iff tag == TT_KEY or TT_LOCK
}

func parseTile(b byte) Tile {
	switch b {
	case '.':
		return Tile{tag: TT_EMPTY}
	case '#':
		return Tile{tag: TT_WALL}
	default:
		if b == '@' {
			return Tile{tag: TT_KEY, kid: KeyID(26)} // special artificial key to make KeyDistances work
		} else if 'a' <= b && b <= 'z' {
			return Tile{tag: TT_KEY, kid: KeyID(b - 'a')}
		} else {
			assert('A' <= b && b <= 'Z', "unknown char %s", b)
			return Tile{tag: TT_LOCK, kid: KeyID(b - 'A')}
		}
	}
}

func (t Tile) String() string {
	switch t.tag {
	case TT_EMPTY:
		return "."
	case TT_WALL:
		return "#"
	case TT_KEY:
		return string(byte(t.kid) + 'a')
	case TT_LOCK:
		return string(byte(t.kid) + 'A')
	default:
		return "?"
	}
}

type point struct {
	x, y int
}

func (p *point) String() string {
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

type Maze struct {
	tiles   [][]Tile
	keyLocs map[KeyID]point // a cache
	start   point
}

func makeMaze() *Maze {
	return &Maze{
		keyLocs: make(map[KeyID]point),
	}
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

type SolveState struct {
	keys    map[KeyID]bool // one per alphabet letter + one for start
	steps   int
	lastKey KeyID
}

func makeSolveState() SolveState {
	return SolveState{
		keys: make(map[KeyID]bool),
	}
}

func (maze *Maze) keyAt(p point) KeyID {
	tile := maze.tiles[p.y][p.x]
	assert(tile.tag == TT_KEY, "tile at %s isn't a key; it's %v", p, tile)
	return tile.kid
}

func (state *SolveState) String() string {
	var s strings.Builder
	fmt.Fprintf(&s, "%d %s [", state.steps, state.lastKey.String())
	for i, b := range state.keys {
		if b {
			fmt.Fprintf(&s, "%s", byte(i+'a'))
		} else {
			fmt.Fprintf(&s, ".")
		}
	}
	fmt.Fprintf(&s, "]")
	return s.String()
}

// store key locations
// solver.dist() gives distance between two points
// solver.available() returns available keys
func solve(maze *Maze) interface{} {
	// fmt.Printf("maze:\n%s\n", maze)

	keyDists := precomputeKeyDistances(maze)
	fmt.Printf("keyDists:\n%s\n", keyDists)

	var stateQueue []SolveState
	s := makeSolveState()
	s.lastKey = 26 // 26 is the start pseudo-key
	stateQueue = append(stateQueue, s)
	for i := 0; i < len(stateQueue); i++ {
		state := stateQueue[i]
		for _, kid := range keyDists.availableKeys(&state) {
			fmt.Printf("Available key: %s\n", kid)
			// TODO
		}
	}
	return nil
}

// availableKeys returns keys that are not yet gotten but are available to be gotten
func (kd KeyDistances) availableKeys(state *SolveState) (res []KeyID) {
	submap := kd[state.lastKey]
	for kid := KeyID(0); kid < 27; kid++ {
		kdist := submap[kid]
		fmt.Println(kid, kdist)
		if state.hasAllKeys(kdist.locks) {
			res = append(res, kid)
		}
	}
	return
}

func (state *SolveState) hasAllKeys(kids []KeyID) bool {
	for _, kid := range kids {
		if !state.keys[kid] {
			return false
		}
	}
	return true
}

// KeyID is 0-25 for an alphabetic key, or 26 for the start pseudo-key
type KeyID int

func (kid KeyID) String() string {
	if kid == 26 {
		return "@"
	}
	return string(kid + 'a')
}

// KeyDistances keeps track of distances between keys.
type KeyDistances map[KeyID]map[KeyID]KeyDist

func (kds KeyDistances) String() string {
	var s strings.Builder
	for kidA := KeyID(0); kidA < 27; kidA++ {
		for kidB := KeyID(0); kidB < 27; kidB++ {
			if kd, prs := kds[kidA][kidB]; prs {
				fmt.Fprintf(&s, "[%s->%s]=%s, ", kidA, kidB, kd)
			}
		}
	}
	return s.String()
}

type KeyDist struct {
	dist  int
	locks []KeyID
}

func (kd KeyDist) String() string {
	var s strings.Builder
	fmt.Fprintf(&s, "%d", kd.dist)
	if len(kd.locks) > 0 {
		fmt.Fprintf(&s, "(")
		for _, lock := range kd.locks {
			fmt.Fprintf(&s, "%s", lock.String())
		}
		fmt.Fprintf(&s, ")")
	}
	return s.String()
}

// assumes the maze has no cycles - whole algorithm falls apart otherwise
func precomputeKeyDistances(maze *Maze) KeyDistances {
	res := make(map[KeyID]map[KeyID]KeyDist)
	for kid := KeyID(0); kid < 27; kid++ {
		res[kid] = precomputeKeyDistancesFrom(maze, kid)
	}
	return res
}

type VisitType int

const (
	VT_UNVISITED VisitType = iota
	VT_VISITED
	VT_FRONTEIR
)

type VisitTracker map[point]VisitInfo

type VisitInfo struct {
	tag   VisitType
	dist  int
	locks []KeyID
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

func precomputeKeyDistancesFrom(maze *Maze, kid KeyID) map[KeyID]KeyDist {
	res := make(map[KeyID]KeyDist)

	visited := makeVisitTracker()
	visited[maze.keyLocs[kid]] = VisitInfo{tag: VT_FRONTEIR}
	for {
		p, info, ok := visited.nextToVisit()
		if !ok {
			break
		}
		locks := info.locks
		dist := info.dist

		tile := maze.tiles[p.y][p.x]
		visited[p] = VisitInfo{tag: VT_VISITED, dist: dist}
		switch tile.tag {
		case TT_EMPTY:
			// nothing special
		case TT_WALL:
			continue // don't add neighbors
		case TT_LOCK:
			locks = append(locks, tile.kid)
		case TT_KEY:
			if tile.kid == kid {
				assert(dist == 0, "loop?") // not an exhaustive loop check
			} else {
				res[tile.kid] = KeyDist{
					dist:  dist,
					locks: locks,
				}
			}
		default:
			assert(false, "unknown tile.tag %s at %s", tile.tag, p)
		}

		for _, np := range p.neighbors() {
			if visited[np].tag == VT_UNVISITED {
				visited[np] = VisitInfo{tag: VT_FRONTEIR, dist: dist + 1, locks: locks}
			}
		}
	}

	return res
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

	maze := makeMaze()
	for r, l := range lines {
		if l == "" {
			continue
		}
		var row []Tile
		for c, b := range []byte(l) {
			tile := parseTile(b)
			row = append(row, tile)
			p := point{x: c, y: r}
			if b == '@' {
				maze.start = p
			}
			if tile.tag == TT_KEY {
				maze.keyLocs[tile.kid] = p
			}
		}

		maze.tiles = append(maze.tiles, row)
	}

	return maze, nil
}
