package computer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
)

var assert = helpers.Assert
var check = helpers.Check

type paramModes struct {
	modes []int
	ptr   int
}

func (m *paramModes) get(n int) int {
	assert(n >= 0, "bad get")
	if n >= len(m.modes) {
		return 0
	}
	return m.modes[n]
}

func (m *paramModes) getNext() int {
	m.ptr++
	return m.get(m.ptr - 1)
}

func parseOpcode(code int) (int, paramModes, error) {
	s := strconv.Itoa(code)
	opcode := code % 100
	var modes []int
	for i := len(s) - 3; i >= 0; i-- {
		// fmt.Printf("i=%v\n", i)
		// fmt.Printf("s[i]=%v, '0'=%v, ==?: %v\n", s[i], '0', s[i] == '0')
		if s[i] == '0' {
			modes = append(modes, 0)
		} else if s[i] == '1' {
			modes = append(modes, 1)
		} else if s[i] == '2' {
			modes = append(modes, 2)
		} else {
			return 0, paramModes{}, fmt.Errorf("Unrecognized mode %q[%d]='%c'", s, i, s[i])
		}
	}
	// fmt.Printf("done\n")
	return opcode, paramModes{modes: modes}, nil
}

func init() {
	opcode, modes, err := parseOpcode(1002)
	assert(err == nil, "t1 err")
	assert(opcode == 2, "t1 opcode")
	assert(modes.getNext() == 0, "t1 modes[0]")
	assert(modes.getNext() == 1, "t1 modes[1]")
	assert(modes.getNext() == 0, "t1 modes[2]")
	assert(modes.getNext() == 0, "t1 modes[3]")
	assert(modes.getNext() == 0, "t1 modes[4]")
	assert(modes.getNext() == 0, "t1 modes[5]")

	opcode, modes, err = parseOpcode(3002)
	assert(err != nil, "t2 err")

	opcode, modes, err = parseOpcode(42)
	assert(err == nil, "t3 err")

	// assert(false, "exit after tests")
}

// Parameter is a generic parameter
type Parameter interface {
	Get() int
	Set(int)
}

// ParamPointer is a pointer to a location in memory
type ParamPointer struct {
	addr int
	cpu  *CPU
}

// Get .
func (p ParamPointer) Get() int {
	helpers.EnsureInbounds(p.cpu.mem, p.addr)
	return p.cpu.mem[p.addr]
}

// Set .
func (p ParamPointer) Set(x int) {
	helpers.EnsureInbounds(p.cpu.mem, p.addr)
	p.cpu.mem[p.addr] = x
}

// String .
func (p ParamPointer) String() string {
	return fmt.Sprintf("[%d]", p.addr)
}

// ParamLiteral is a literal value
type ParamLiteral struct {
	val int
}

// Get .
func (p ParamLiteral) Get() int {
	return p.val
}

// Set .
func (p ParamLiteral) Set(x int) {
	panic("trying to Set a ParamLiteral")
}

// String .
func (p ParamLiteral) String() string {
	return fmt.Sprintf("%d", p.val)
}

// ParamRelativePointer is a pointer to a relative location in memory
type ParamRelativePointer struct {
	addr int
	cpu  *CPU
}

// Get .
func (p ParamRelativePointer) Get() int {
	helpers.EnsureInbounds(p.cpu.mem, p.cpu.relBase+p.addr)
	return p.cpu.mem[p.cpu.relBase+p.addr]
}

// Set .
func (p ParamRelativePointer) Set(x int) {
	helpers.EnsureInbounds(p.cpu.mem, p.cpu.relBase+p.addr)
	p.cpu.mem[p.cpu.relBase+p.addr] = x
}

// String .
func (p ParamRelativePointer) String() string {
	return fmt.Sprintf("[RB+%d]", p.addr)
}

// NextParameter changes cpu state to return the next parameter at the pc
func (cpu *CPU) NextParameter() Parameter {
	v := chomp(cpu.mem, &cpu.pc)
	return MakeParameter(cpu, cpu.modes.getNext(), v)
}

// MakeParameter returns a parameter based on parameter code/mode int
func MakeParameter(cpu *CPU, mode int, val int) Parameter {
	switch mode {
	case 0:
		return ParamPointer{addr: val, cpu: cpu}
	case 1:
		return ParamLiteral{val: val}
	case 2:
		return ParamRelativePointer{addr: val, cpu: cpu}
	}
	assert(false, fmt.Sprintf("bad MakeParameter mode %d", mode))
	return nil
}

// State .
type State int

const (
	CS_NONE State = iota
	CS_WAITING_INPUT
	CS_WAITING_OUTPUT
	CS_DONE
)

func (s State) String() string {
	switch s {
	case CS_NONE:
		return "CS_NONE"
	case CS_WAITING_INPUT:
		return "CS_WAITING_INPUT"
	case CS_WAITING_OUTPUT:
		return "CS_WAITING_OUTPUT"
	case CS_DONE:
		return "CS_DONE"
	}
	return "CS_?"
}

// CPU .
type CPU struct {
	Name      string
	InChan    chan int
	OutChan   chan int
	StateChan chan State
	Halted    bool

	pc      int
	modes   paramModes
	mem     []int
	relBase int
}

// MakeCPU makes a CPU
func MakeCPU(name string) CPU {
	return CPU{
		Name:      name,
		InChan:    make(chan int, 1),
		OutChan:   make(chan int, 1),
		StateChan: make(chan State, 1),
	}
}

// NextOpcode changes cpu state to read the next opcode and set cpu.modes
func (cpu *CPU) NextOpcode() int {
	code := chomp(cpu.mem, &cpu.pc)
	opcode, modes, err := parseOpcode(code)
	check(err)

	cpu.modes = modes
	return opcode
}

// TODO opcode type

// SendInput sends the input iff the computer's StateChan indicates
// it's ready for it. otherwise, it panics
func (cpu *CPU) SendInput(x int) {
	state := <-cpu.StateChan
	assert(state == CS_WAITING_INPUT, "tried to SendInput but state was %v", state)
	cpu.InChan <- x
}

// RecvOutput receives the output iff the computer's StateChan indicates
// it's ready for it. otherwise, it panics
func (cpu *CPU) RecvOutput() int {
	state := <-cpu.StateChan
	assert(state == CS_WAITING_OUTPUT, "tried to RecvOutput but state was %v", state)
	return <-cpu.OutChan
}

// Run runs the cpu in a goroutine
func (cpu *CPU) Run() {
	go func() {
		for cycles := 0; !cpu.Halted; cycles++ {
			// if cycles%1000 == 0 {
			// 	fmt.Printf("cycles: %d\n", cycles)
			// }
			code := cpu.NextOpcode()

			// fmt.Printf("node %s pc=%v code=%d modes=%v\n", cpu.name, cpu.pc, code, cpu.modes)
			// cpu.dump()

			switch code {
			case 1: // add
				// fmt.Println(" add")
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				c := cpu.NextParameter()
				c.Set(a.Get() + b.Get())
			case 2: // mult
				// fmt.Println(" mult")
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				c := cpu.NextParameter()
				c.Set(a.Get() * b.Get())
			case 3: // input
				// fmt.Println(" input")
				a := cpu.NextParameter()
				cpu.StateChan <- CS_WAITING_INPUT
				i := <-cpu.InChan
				// fmt.Printf("%s < %d\n", cpu.Name, i)
				a.Set(i)
			case 4: // output
				// fmt.Println(" output")
				a := cpu.NextParameter()
				cpu.StateChan <- CS_WAITING_OUTPUT
				cpu.OutChan <- a.Get()
			case 5: // jump-if-true
				// fmt.Println(" jump-if-true")
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				if a.Get() != 0 {
					cpu.pc = b.Get()
				}
			case 6: // jump-if-false
				// fmt.Println(" jump-if-false")
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				if a.Get() == 0 {
					cpu.pc = b.Get()
				}
			case 7: // less than
				// fmt.Println(" less than")
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				c := cpu.NextParameter()
				if a.Get() < b.Get() {
					c.Set(1)
				} else {
					c.Set(0)
				}
			case 8: // equals
				// fmt.Println(" equals")
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				c := cpu.NextParameter()
				if a.Get() == b.Get() {
					c.Set(1)
				} else {
					c.Set(0)
				}
			case 9: // adjust relative parameter base
				// fmt.Println(" adjust relative parameter base")
				a := cpu.NextParameter()
				cpu.relBase += a.Get()
			case 99: // halt
				// fmt.Println(" halt")
				cpu.Halted = true
				cpu.StateChan <- CS_DONE
				close(cpu.StateChan)
				close(cpu.InChan)
				close(cpu.OutChan)
			default:
				panic("unknown opcode")
			}
		}
	}()
}

// MemSize is the length of memory the cpu has
const MemSize = 5000

// SetMemory initializes the cpu memory
func (cpu *CPU) SetMemory(mem []int) {
	cpu.mem = make([]int, MemSize)
	copy(cpu.mem, mem)
}

func chomp(mem []int, pc *int) int {
	helpers.EnsureInbounds(mem, *pc)
	res := mem[*pc]
	*pc++
	return res
}

func (cpu *CPU) dump() {
	fmt.Printf("name=%s mem=[", cpu.Name)
	for i, v := range cpu.mem {
		if i%10 == 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("%3d ", v)
	}
	fmt.Printf("\n]\n")
}

// MockCPU is useful for testing
type MockCPU struct {
	InChan   chan int
	OutChan  chan int
	DoneChan chan struct{}
	Halted   bool
}

// MakeMockCPU makes a mock CPU that will ignore inputs
// and output the given outputs
func MakeMockCPU(outputs []int) *MockCPU {
	mock := MockCPU{
		InChan:   make(chan int, 1),
		OutChan:  make(chan int, 1),
		DoneChan: make(chan struct{}, 1),
	}

	// ignore input
	go func() {
		for _ = range mock.InChan {
		}
	}()

	// produce output
	go func() {
		for _, x := range outputs {
			mock.OutChan <- x
		}
		mock.Halted = true
		mock.DoneChan <- struct{}{}
	}()

	return &mock
}

type JumpLocationType int

const (
	JLT_NONE JumpLocationType = iota
	JLT_QUEUED
	JLT_DONE_BORING
	JLT_DONE
)

type JumpLocations map[int]JumpLocationType

func makeJumpLocations() JumpLocations {
	return make(map[int]JumpLocationType)
}

func (j JumpLocations) enqueue(loc int) {
	switch j[loc] {
	case JLT_NONE:
		j[loc] = JLT_QUEUED
	case JLT_QUEUED:
	case JLT_DONE_BORING:
		j[loc] = JLT_DONE
	case JLT_DONE:
	}
}

func (j JumpLocations) addDoneBoring(loc int) {
	switch j[loc] {
	case JLT_NONE:
		j[loc] = JLT_DONE_BORING
	case JLT_QUEUED:
		j[loc] = JLT_DONE
	case JLT_DONE_BORING:
	case JLT_DONE:
	}
}

func (j JumpLocations) pop() (int, bool) {
	for k, v := range j {
		if v == JLT_QUEUED {
			j[k] = JLT_DONE
			return k, true
		}
	}
	return 0, false
}

func (j JumpLocations) query(loc int) JumpLocationType {
	return j[loc]
}

func (cpu *CPU) PrintProgram() string {
	jumpLocs := makeJumpLocations()
	jumpLocs.enqueue(0)

	type Line struct {
		n    int
		line string
	}
	var lines []Line
	for {
		if loc, ok := jumpLocs.pop(); !ok {
			break
		} else {
			cpu.pc = loc
		}

	RUNNER:
		for cpu.pc < len(cpu.mem) {
			jumpLocs.addDoneBoring(cpu.pc)
			var s strings.Builder
			pc := cpu.pc

			code := cpu.NextOpcode()
			switch code {
			case 1: // add
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				c := cpu.NextParameter()
				fmt.Fprintf(&s, "%s = %s+%s\n", c, a, b)
			case 2: // mult
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				c := cpu.NextParameter()
				fmt.Fprintf(&s, "%s = %s*%s\n", c, a, b)
			case 3: // input
				a := cpu.NextParameter()
				fmt.Fprintf(&s, "%s = input()\n", a)
			case 4: // output
				a := cpu.NextParameter()
				fmt.Fprintf(&s, "output(%s)\n", a)
			case 5: // jump-if-true
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				fmt.Fprintf(&s, "if %s { goto %s }", a, b)
				if lit, ok := b.(ParamLiteral); ok {
					jumpLocs.enqueue(lit.val)
				} else {
					fmt.Fprintf(&s, " // UNFOLLOWABLE")
				}
				fmt.Fprintf(&s, "\n")
			case 6: // jump-if-false
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				fmt.Fprintf(&s, "if !%s { goto %s }", a, b)
				if lit, ok := b.(ParamLiteral); ok {
					jumpLocs.enqueue(lit.val)
				} else {
					fmt.Fprintf(&s, " // UNFOLLOWABLE")
				}
				fmt.Fprintf(&s, "\n")
			case 7: // less than
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				c := cpu.NextParameter()
				fmt.Fprintf(&s, "%s = (%s < %s)\n", c, a, b)
			case 8: // equals
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				c := cpu.NextParameter()
				fmt.Fprintf(&s, "%s = (%s == %s)\n", c, a, b)
			case 9: // adjust relative parameter base
				a := cpu.NextParameter()
				fmt.Fprintf(&s, "RB += %s\n", a)
			case 99: // halt
				fmt.Fprintf(&s, "exit\n")
			default:
				break RUNNER
			}
			lines = append(lines, Line{n: pc, line: s.String()})
		}
	}

	var s strings.Builder
	for _, l := range lines {
		n := l.n
		line := l.line

		doneType := jumpLocs.query(n)
		if doneType == JLT_DONE {
			fmt.Fprintf(&s, "\n%3d*: ", n)
		} else {
			assert(doneType == JLT_DONE_BORING, "should be JLT_DONE_BORING (%d)", doneType)
			fmt.Fprintf(&s, "%3d : ", n)
		}
		fmt.Fprintf(&s, line)
	}
	cpu.pc = 0
	return s.String()
}
