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

		Mem Memory

		Dsp Display

		Key Keys
	}

	// Stack is the CHIP-8 stack with 16 levels
	Stack struct {
		s  [16]uint16
		sp uint16
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

func (m *Memory) FetchOpcode(addr uint16) OpCode {
	var opcode uint16
	opcode = uint16(m[addr]) << 8
	opcode = opcode | uint16(m[addr+1])
	return OpCode(opcode)
}
