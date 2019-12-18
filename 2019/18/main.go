package main

import (
	"fmt"
	"strings"

	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
)

var assert = helpers.Assert
var check = helpers.Check

var HighestKeyID KeyID

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
	keys    [27]bool // one per alphabet letter + one for start
	lastKey KeyID
}

func (state *SolveState) collect(lastKey KeyID, otherNewKeys [27]bool) SolveState {
	newState := *state // copy

	for kid, val := range otherNewKeys {
		// if val && kid != lastKey && !newState.keys[kid] {
		// 	fmt.Println("force-collecting", kid)
		// }
		if val {
			newState.keys[kid] = true
		}
	}

	newState.keys[lastKey] = true // (redundant atm)
	newState.lastKey = lastKey

	return newState
}

func makeSolveState() SolveState {
	return (&SolveState{}).collect(26, [27]bool{}) // 26 is the start pseudo-key
}

func (maze *Maze) keyAt(p point) KeyID {
	tile := maze.tiles[p.y][p.x]
	assert(tile.tag == TT_KEY, "tile at %s isn't a key; it's %v", p, tile)
	return tile.kid
}

func (sd StateDist) String() string {
	return fmt.Sprintf("%d %s", sd.dist, sd.state)
}

func (state SolveState) String() string {
	var s strings.Builder
	fmt.Fprintf(&s, "(%s)[", state.lastKey)
	for kid, haveKey := range state.keys {
		if haveKey {
			fmt.Fprintf(&s, "%s", KeyID(kid))
		}
	}
	fmt.Fprintf(&s, "]")
	return s.String()
}

type StateDist struct {
	state SolveState
	dist  int
}

func solve(maze *Maze) interface{} {
	// fmt.Printf("maze:\n%s\n", maze)

	keyDists := precomputeKeyDistances(maze)
	// fmt.Printf("keyDists:\n%s\n", keyDists)

	var stateQueue []StateDist
	stateQueue = append(stateQueue, StateDist{
		state: makeSolveState(),
		dist:  0,
	})
	bestDists := make(map[SolveState]int)
	var skipCount int
	for i := 0; i < len(stateQueue); i++ {
		current := stateQueue[i]
		shouldPrint := i%10000 == 0
		shouldPrint = false
		if shouldPrint {
			fmt.Printf("Cycle %d/%d:\n", i+1, len(stateQueue))
			// fmt.Printf("  skipped:     %d\n", skipCount)
			fmt.Printf("  current: %s\n", current)
		}
		if oldBest, prs := bestDists[current.state]; prs && oldBest <= current.dist {
			// if we've already analyzed this state from a better starting point (on some earlier iteration)

			if shouldPrint {
				fmt.Printf("  (out of date, skipping (%d<=%d))\n", oldBest, current.dist)
			}
			skipCount++
			continue
		} else {
			bestDists[current.state] = current.dist
			if shouldPrint {
				fmt.Printf("  new best: %d->%d\n", oldBest, current.dist)
			}
		}

		for _, ktemp := range keyDists.availableKeys(&current.state) {
			new := StateDist{
				state: current.state.collect(ktemp.KeyID, ktemp.KeyDist.keys),
				dist:  current.dist + ktemp.KeyDist.dist,
			}
			if shouldPrint {
				fmt.Printf("    ->: %s\n", new)
			}
			stateQueue = append(stateQueue, new)
		}
	}

	best := 1000000
	for state, dist := range bestDists {
		if !state.done() {
			continue
		}
		// fmt.Println("candidate", state, dist)
		if dist < best {
			best = dist
		}
	}
	return best
}

type KeyTemp struct {
	KeyDist
	KeyID
}

// availableKeys returns keys that are available to be gotten
// note: this will often return keys we've already gotten (useful as checkpoints)
func (kd KeyDistances) availableKeys(state *SolveState) (res []KeyTemp) {
	submap := kd[state.lastKey]
	for kid := KeyID(0); kid < 27; kid++ {
		kdist, prs := submap[kid]
		// fmt.Printf("  KeyDist(%s)=%s\n", kid, kdist)
		if prs && state.hasAllKeys(kdist.locks) {
			res = append(res, KeyTemp{kdist, kid})
		}
	}
	return
}

func (state *SolveState) done() bool {
	for kid, haveKey := range state.keys {
		if KeyID(kid) > HighestKeyID {
			break
		}
		if !haveKey {
			return false
		}
	}
	return true
}

func (state *SolveState) hasAllKeys(kids [27]bool) bool {
	// unset keys we have
	for kid, haveKey := range state.keys {
		if haveKey {
			kids[KeyID(kid)] = false
		}
	}

	// are there any keys still set?
	for _, v := range kids {
		if v {
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
	locks [27]bool
	keys  [27]bool
}

func (kd KeyDist) String() string {
	var s strings.Builder
	fmt.Fprintf(&s, "%d", kd.dist)
	if len(kd.locks) > 0 {
		fmt.Fprintf(&s, "(")
		for kid, val := range kd.locks {
			if val {
				fmt.Fprintf(&s, "%s", strings.ToUpper(KeyID(kid).String()))
			}
		}
		fmt.Fprintf(&s, ")")
	}
	if len(kd.keys) > 0 {
		fmt.Fprintf(&s, "(")
		for kid, val := range kd.keys {
			if val {
				fmt.Fprintf(&s, "%s", KeyID(kid).String())
			}
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
	locks [27]bool
	keys  [27]bool
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
		keys := info.keys
		dist := info.dist

		tile := maze.tiles[p.y][p.x]
		visited[p] = VisitInfo{tag: VT_VISITED, dist: dist}
		switch tile.tag {
		case TT_EMPTY:
			// nothing special
		case TT_WALL:
			continue // don't add neighbors
		case TT_LOCK:
			locks[tile.kid] = true
		case TT_KEY:
			keys[tile.kid] = true
			if tile.kid == kid {
				assert(dist == 0, "loop?") // note: this is not an exhaustive loop check
			} else {
				res[tile.kid] = KeyDist{
					dist:  dist,
					locks: locks,
					keys:  keys,
				}
			}
		default:
			assert(false, "unknown tile.tag %s at %s", tile.tag, p)
		}

		for _, np := range p.neighbors() {
			if visited[np].tag == VT_UNVISITED {
				visited[np] = VisitInfo{tag: VT_FRONTEIR, dist: dist + 1, locks: locks, keys: keys}
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
				if tile.kid != 26 && tile.kid > HighestKeyID {
					HighestKeyID = tile.kid
				}
			}
		}

		maze.tiles = append(maze.tiles, row)
	}

	return maze, nil
}
