package apparat

type (
	// System describes the CHIP-8 system
	System struct {
		// V is the representation of the registers V0-VE
		V [16]byte
		// I is the index register
		I uint16
		// PC is the program counter
		PC uint16

		Stack *Stack

		Mem Memory

		Dsp Display

		Key Keys
	}

	// Stack is the CHIP-8 stack with 16 levels
	Stack struct {
		s  [16]uint16
		sp byte
	}

	// Memory is the 4K CHIP-8 memory
	Memory [4096]byte

	// Display represents the display state
	Display [32]uint64

	// Keys represents the keyboard state
	Keys byte

	// OpCode represents an opcode
	OpCode uint16
)

// NewSystem initializes a new system
func NewSystem() *System {
	return &System{
		PC:    0x200,
		Stack: &Stack{},
	}
}

// Reset resets the system
func (s *System) Reset() {
	s.V = [16]byte{}
	s.I = 0
	s.PC = 0x200
	s.Stack.Reset()

	s.Mem = Memory([4096]byte{})
	s.Dsp = [32]uint64{}
	s.Key = 0
}

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
	s.sp += 1
}

// Pop pops from the stack
func (s *Stack) Pop() uint16 {
	if s.sp == 0 {
		panic("stack underflow")
	}
	s.sp -= 1
	v := s.s[s.sp]
	return v
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
