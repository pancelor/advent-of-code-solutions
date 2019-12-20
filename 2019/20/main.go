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

func getInput() *Maze {
	lines, err := helpers.GetLines()
	check(err)

	l, lines = lines[0], lines[1:]
	labels := strings.Split(l, " ")
	nextLabel := 0

	var maze Maze
	for r, l := range lines {
		if l == "" {
			continue
		}
		edgeV := r == 0 || r == 110
		var row []Tile
		for c, b := range []byte(l) {
			edgeH := c == 0 || r == 108
			tile := parseTile(b)
			if (edgeH || edgeV) && tile.tag == TT_EMPTY {
				tile = Tile{tag: TT_TUNNEL, labels[nextLabel]}
				nextLabel++
			}
			row = append(row, tile)
		}

		maze.tiles = append(maze.tiles, row)
	}

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
	tid KeyID // optional; exists iff tag == TT_TUNNEL
}

func parseTile(b byte) Tile {
	switch b {
	case '.':
		return Tile{tag: TT_EMPTY}
	case '#':
		return Tile{tag: TT_WALL}
	case ' ':
		return Tile{tag: TT_OOB}
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

// TunnelID .
type TunnelID int

func (tid TunnelID) String() string {
	switch kid {
	case 26:
		return "@"
	case 27:
		return "$"
	case 28:
		return "%"
	case 29:
		return "&"
	default:
		return string(kid + 'a')
	}
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

// tunnelAt returns the KeyID of the key at the given position,
// or panics if no such key exists
func (maze *Maze) tunnelAt(p point) TunnelID {
	tile := maze.tiles[p.y][p.x]
	assert(tile.tag == TT_TUNNEL, "tile at %s isn't a tunnel; it's %v", p, tile)
	return tile.kid
}

// TODO return both?
// tLoc returns the locations of the given TunnelID,
// or panics if no such tunnel exists
func (maze *Maze) tLoc(tid TunnelID) point {
	for r, row := range maze.tiles {
		for c, tile := range row {
			if tile.tag == TT_KEY && tile.kid == kid {
				return point{x: c, y: r}
			}
		}
	}
	assert(false, "TunnelID %s is noplace", kid)
	return point{}
}

///////////////////////////////////////////////////////////////////////////////
//
// KeyDistances computation
//
///////////////////////////////////////////////////////////////////////////////

// KeyDistances keeps track of distances between keys.
type KeyDistances map[KeyID]map[KeyID]KeyDist

func (kds KeyDistances) String() string {
	var s strings.Builder
	for kidA := KeyID(0); kidA < NKeys; kidA++ {
		for kidB := KeyID(0); kidB < NKeys; kidB++ {
			if kd, prs := kds[kidA][kidB]; prs {
				fmt.Fprintf(&s, "[%s->%s]=%s, ", kidA, kidB, kd)
			}
		}
	}
	return s.String()
}

// KeyDist represents the distance between two keys,
// along with the locks that are blocking this path
// and the keys that are "blocking" the path
// (storing the keys helps us reduce the branching factor of our algorithm
// by eagerly picking up keys that are along the path)
// !!! This assumes the maze has no cycles - the whole algorithm falls apart otherwise!!!
type KeyDist struct {
	dist  int
	locks [NKeys]bool
	keys  [NKeys]bool
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

// assumes the maze has no cycles - the whole algorithm falls apart otherwise
func precomputeKeyDistances(maze *Maze) KeyDistances {
	res := make(map[KeyID]map[KeyID]KeyDist)
	for kid := KeyID(0); kid < NKeys; kid++ {
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
	locks [NKeys]bool
	keys  [NKeys]bool
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
	if kid.unused() {
		return res
	}

	visited := makeVisitTracker()
	visited[maze.keyLoc(kid)] = VisitInfo{tag: VT_FRONTEIR}
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

///////////////////////////////////////////////////////////////////////////////
//
// Maze solving
//
///////////////////////////////////////////////////////////////////////////////

// SolveState is a pointer-free value representing the state of an in-progress solution
// This value must have no pointers (e.g. slices) so that it can be used as a hashmap key
// fields:
//   keys: the keys that the solution has collected
//   lastKeys: the last key that each robot collected last; i.e. where each robot is
// Note that this struct does _not_ include the distance that it took
//   to reach this state; if you want that, see SolveStateDist
type SolveState struct {
	keys     [NKeys]bool    // one per alphabet letter + four starts
	lastKeys [NRobots]KeyID // lastKey for each of the four robots
}

func makeSolveState() SolveState {
	state := SolveState{}
	for rid := 0; rid < NRobots; rid++ {
		state = state.collect(rid, KeyID(26+rid), [NKeys]bool{}) // collect psuedo-keys at starting locations
	}
	return state
}

func (state SolveState) String() string {
	var s strings.Builder
	fmt.Fprintf(&s, "(%v)[", state.lastKeys)
	for kid, haveKey := range state.keys {
		if haveKey {
			fmt.Fprintf(&s, "%s", KeyID(kid))
		}
	}
	fmt.Fprintf(&s, "]")
	return s.String()
}

// collect collects the given key (and any in-the-way keys) and returns the new SolveState
// fields:
//   rid: which robot collected the key
//   lastKey: which key to move to
//   otherNewKeys: keys that the robot picked up along the way to lastKey
// since this isn't a pointer method, state is already a copy we can edit at will
func (state SolveState) collect(rid int, lastKey KeyID, otherNewKeys [NKeys]bool) SolveState {
	for kid, val := range otherNewKeys {
		// if val && kid != lastKey && !newState.keys[kid] {
		// 	fmt.Println("force-collecting roadside key", kid)
		// }
		if val {
			state.keys[kid] = true
		}
	}

	state.keys[lastKey] = true // (redundant atm, but better safe)
	state.lastKeys[rid] = lastKey

	return state
}

// numKeys returns the number of (real) keys in this state
func (state *SolveState) numKeys() int {
	count := 0
	for kid, haveKey := range state.keys {
		if KeyID(kid) > HighestKeyID {
			break
		}
		if haveKey {
			count++
		}
	}
	return count
}

// done returns whether the state is fully solved (i.e. all keys are collected)
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

// hasAllKeys returns whether all keys in the given list have been collected
func (state *SolveState) hasAllKeys(kids [NKeys]bool) bool {
	// unset keys we're wondering about
	// (kids is copied by value so we're not doing anything bad by editing it directly)
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

// SolveStateDist is a tuple storing how many steps it took to reach a certain SolveState
// The distance is not a property of the SolveState itself; there may be different ways to
// reach the same SolveState that took different numbers of steps
type SolveStateDist struct {
	state SolveState
	dist  int
}

func (sd SolveStateDist) String() string {
	return fmt.Sprintf("%d %s", sd.dist, sd.state)
}

// solve returns the minimum number of steps it takes to solve the maze
func solve(maze *Maze) int {
	// fmt.Printf("maze:\n%s\n", maze)
	return solveKD(precomputeKeyDistances(maze))
}

func solveKD(keyDists KeyDistances) int {
	// fmt.Printf("keyDists:\n%s\n", keyDists)

	// stateQueue holds the queue of partial solutions
	// that we've tried / want to try
	// (TODO: release old states after processing them; don't store entire queue in memory)
	var stateQueue []SolveStateDist
	stateQueue = append(stateQueue, SolveStateDist{
		state: makeSolveState(),
		dist:  0,
	})

	// bestDists stores the best distance for each SolveState seen so far
	bestDists := make(map[SolveState]int)

	var skipCount int
	for i := 0; i < len(stateQueue); i++ {
		// pop current from stateQueue
		current := stateQueue[i]

		fmt.Println(current.state.numKeys())
		// shouldPrint := i%10000 == 0
		shouldPrint := false
		// shouldPrint := true

		if shouldPrint {
			fmt.Printf("Cycle %d/%d:\n", i+1, len(stateQueue))
			// fmt.Printf("  skipped:     %d\n", skipCount)
			fmt.Printf("  current: %s\n", current)
			fmt.Printf("  numKeys: %d\n", current.state.numKeys())
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
			if shouldPrint && prs {
				fmt.Printf("  new best: %d->%d\n", oldBest, current.dist)
			}
		}

		for _, ak := range current.state.availableKeys(keyDists) {
			new := SolveStateDist{
				state: current.state.collect(ak.rid, ak.KeyID, ak.KeyDist.keys),
				dist:  current.dist + ak.KeyDist.dist,
			}
			if shouldPrint {
				fmt.Printf("    ->: %s\n", new)
			}
			stateQueue = append(stateQueue, new)
		}
	}
	fmt.Printf("Explored %d SolveStates\n", len(stateQueue))

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

// AvailableKey is a tuple saying what key is available to what robot,
// and what the KeyDist is
type AvailableKey struct {
	KeyDist
	KeyID
	rid int
}

// availableKeys returns keys that are available to be gotten
// (and are not already gotten)
func (state *SolveState) availableKeys(kd KeyDistances) (res []AvailableKey) {
	// fmt.Println("state=%s\n", state)
	// fmt.Println("kd=%s\n", kd)
	for rid := 0; rid < NRobots; rid++ {
		submap := kd[state.lastKeys[rid]]
		for kid := KeyID(0); kid < NKeys; kid++ {
			if state.keys[kid] {
				continue
			}
			kdist, prs := submap[kid]
			// fmt.Printf("  KeyDist(%s)=%s\n", kid, kdist)
			if prs && state.hasAllKeys(kdist.locks) {
				res = append(res, AvailableKey{KeyDist: kdist, KeyID: kid, rid: rid})
			}
		}
	}
	return
}
