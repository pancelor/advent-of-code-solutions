package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func solve(memTemplate []int) (interface{}, error) {
	cpu := makeCPU("A")
	cpu.makeMemory(memTemplate)
	cpu.Run()
	cpu.InChan <- 2

	var res []int
	for {
		select {
		case <-cpu.DoneChan:
			return res, nil
		case x := <-cpu.OutChan:
			res = append(res, x)
		}
	}

	return res, nil
}

// Modes represent argument parameter modes
// get(0) == 1 means the first parameter should be interpreted in
//   immediate mode; 0 means position mode
type Modes struct {
	modes []int
	ptr   int
}

func (m *Modes) get(n int) int {
	assert(n >= 0, "bad get")
	if n >= len(m.modes) {
		return 0
	}
	return m.modes[n]
}

func (m *Modes) getNext() int {
	m.ptr++
	return m.get(m.ptr - 1)
}

func makeParameter(cpu *CPU, mode int, val int) Parameter {
	switch mode {
	case 0:
		return ParamPointer{addr: val, cpu: cpu}
	case 1:
		return ParamLiteral{val: val}
	case 2:
		return ParamRelativePointer{addr: val, cpu: cpu}
	}
	assert(false, fmt.Sprintf("bad makeParameter mode %d", mode))
	return nil
}

func parseOpcode(code int) (int, Modes, error) {
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
			return 0, Modes{}, fmt.Errorf("Unrecognized mode %q[%d]='%c'", s, i, s[i])
		}
	}
	// fmt.Printf("done\n")
	return opcode, Modes{modes: modes}, nil
}

// CPU .
type CPU struct {
	Name     string
	InChan   chan int
	OutChan  chan int
	DoneChan chan struct{}

	pc      int
	modes   Modes
	mem     []int
	relBase int
	halted  bool
	err     error
}

func makeCPU(name string) CPU {
	return CPU{
		Name:     name,
		InChan:   make(chan int),
		OutChan:  make(chan int),
		DoneChan: make(chan struct{}),
	}
}

func (cpu *CPU) check(err error) {
	check(err)
}

// NextOpcode changes cpu state to read the next opcode and set cpu.modes
func (cpu *CPU) NextOpcode() int {
	code := chomp(cpu.mem, &cpu.pc)
	opcode, modes, err := parseOpcode(code)
	cpu.check(err)

	cpu.modes = modes
	return opcode
}

// TODO opcode type

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
	ensureInbounds(p.cpu.mem, p.addr)
	return p.cpu.mem[p.addr]
}

// Set .
func (p ParamPointer) Set(x int) {
	ensureInbounds(p.cpu.mem, p.addr)
	p.cpu.mem[p.addr] = x
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

// ParamRelativePointer is a pointer to a relative location in memory
type ParamRelativePointer struct {
	addr int
	cpu  *CPU
}

// Get .
func (p ParamRelativePointer) Get() int {
	ensureInbounds(p.cpu.mem, p.cpu.relBase+p.addr)
	return p.cpu.mem[p.cpu.relBase+p.addr]
}

// Set .
func (p ParamRelativePointer) Set(x int) {
	ensureInbounds(p.cpu.mem, p.cpu.relBase+p.addr)
	p.cpu.mem[p.cpu.relBase+p.addr] = x
}

// NextParameter changes cpu state to return the next parameter at the pc
func (cpu *CPU) NextParameter() Parameter {
	v := chomp(cpu.mem, &cpu.pc)
	return makeParameter(cpu, cpu.modes.getNext(), v)
}

// Run runs the cpu in a goroutine
func (cpu *CPU) Run() {
	go func() {
		for cycles := 0; !cpu.halted; cycles++ {
			if cycles%1000 == 0 {
				fmt.Printf("cycles: %d\n", cycles)
			}
			code := cpu.NextOpcode()

			// fmt.Printf("node %s pc=%v code=%d modes=%v\n", cpu.name, cpu.pc, code, cpu.modes)
			// cpu.dump()

			switch code {
			case 1: // add
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				c := cpu.NextParameter()
				c.Set(a.Get() + b.Get())
			case 2: // mult
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				c := cpu.NextParameter()
				c.Set(a.Get() * b.Get())
			case 3: // input
				a := cpu.NextParameter()
				i := <-cpu.InChan
				// fmt.Printf("%s < %d\n", name, i)
				a.Set(i)
			case 4: // output
				a := cpu.NextParameter()
				cpu.OutChan <- a.Get()
			case 5: // jump-if-true
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				if a.Get() != 0 {
					cpu.pc = b.Get()
				}
			case 6: // jump-if-false
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				if a.Get() == 0 {
					cpu.pc = b.Get()
				}
			case 7: // less than
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				c := cpu.NextParameter()
				if a.Get() < b.Get() {
					c.Set(1)
				} else {
					c.Set(0)
				}
			case 8: // equals
				a := cpu.NextParameter()
				b := cpu.NextParameter()
				c := cpu.NextParameter()
				if a.Get() == b.Get() {
					c.Set(1)
				} else {
					c.Set(0)
				}
			case 9: // adjust relative parameter base
				a := cpu.NextParameter()
				cpu.relBase += a.Get()
			case 99: // halt
				cpu.halted = true
				cpu.DoneChan <- struct{}{}
			default:
				panic("unknown opcode")
			}
		}
	}()
}

// MemSize .
const MemSize = 5000

func (cpu *CPU) makeMemory(mem []int) {
	cpu.mem = make([]int, MemSize)
	copy(cpu.mem, mem)
}

func chomp(mem []int, pc *int) int {
	ensureInbounds(mem, *pc)
	res := mem[*pc]
	*pc++
	return res
}

func ensureInbounds(mem []int, ptr ...int) {
	for _, p := range ptr {
		assert(0 <= p && p < len(mem), fmt.Sprintf("oob: %d", p))
	}
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

func test() {
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

func main() {
	test()

	input, err := getInput()
	if err != nil {
		panic(err)
	}

	answer, err := solve(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("answer:\n%v\n", answer)
}

//
// helpers
//

// for temporary use only
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// for temporary use only
func assert(b bool, msg string) {
	if !b {
		panic(errors.New(msg))
	}
}

var source = os.Stdin

func getInput() ([]int, error) {
	lines, err := getLines()
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

		res = append(res, v)
	}

	return res, nil
}

func getLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
