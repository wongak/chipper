package apparat

import "sync"

type (
	// Stack is the CHIP-8 stack with 16 levels
	Stack struct {
		s  [16]uint16
		sp byte
	}

	// Memory is the 4K CHIP-8 memory
	Memory [4096]byte

	// Display represents the display state
	Display struct {
		m *sync.RWMutex
		d [32]uint64
	}

	// Keys represents the keyboard state
	Keys byte
)

// Reset resets the stack
func (s *Stack) Reset() {
	s.s = [16]uint16{}
	s.sp = 0
}

// Push pushes onto the stack
func (s *Stack) Push(a uint16) {
	if s.sp > 14 {
		panic("stack overflow")
	}
	s.s[s.sp] = a
	s.sp++
}

// Pop pops from the stack
func (s *Stack) Pop() uint16 {
	if s.sp == 0 {
		panic("stack underflow")
	}
	s.sp--
	v := s.s[s.sp]
	return v
}

// NewMemory creates a memory with the preloaded CHIP-8 data
func NewMemory() Memory {
	m := [4096]byte{}
	for i, f := range fontset {
		m[0x50+i] = f
	}
	return Memory(m)
}

// FetchOpcode retrieves an opcode from the given address
//
// An opcode consists of 2 bytes, therefore we need to retrieve
// the bytes from PC and PC+1
func (m *Memory) FetchOpcode(addr uint16) OpCode {
	var opcode uint16
	opcode = uint16(m[addr]) << 8
	opcode = opcode | uint16(m[addr+1])
	return OpCode(opcode)
}

func NewDisplay() *Display {
	return &Display{
		m: &sync.RWMutex{},
		d: [32]uint64{},
	}
}
