package apparat

import "fmt"

type (
	// OpCode represents an opcode
	OpCode uint16
)

// Instruction returns the instruction byte
func (o OpCode) Instruction() uint8 {
	return uint8((0xF000 & uint16(o)) >> 12)
}

// ExtractAddr extracts 2 bytes for the address instructions
func (o OpCode) ExtractAddr() uint16 {
	return uint16(o) & 0x0FFF
}

// ExtractVNN extracts 2 bytes for the register/byte instructions
func (o OpCode) ExtractVNN() (uint8, uint8) {
	v := uint8((uint16(o) & 0x0F00) >> 8)
	cmp := uint8(uint16(o) & 0x00FF)
	return v, cmp
}

func (s *System) executeOpcode() {
	op := s.Mem.FetchOpcode(s.PC)
	switch true {
	// 0x00E0: Clear the screen
	case 0x00E0 == op:

	// 0x00EE: return from a subroutine
	// RET
	case 0x00EE == op:
		s.PC = s.Stack.Pop()

	// 0x1NNN: JMP
	// JMP 0xNNN
	case op.Instruction() == 0x1:
		s.PC = op.ExtractAddr()

	// 0x2NNN: SUB
	// SUB 0xNNN
	case op.Instruction() == 0x2:
		s.Stack.Push(s.PC)
		s.PC = op.ExtractAddr()

	// 0x3XNN: if(VX==NN)
	// SIE X NN (skip if equals)
	case op.Instruction() == 0x3:
		v, cmp := op.ExtractVNN()
		if cmp == s.V[v] {
			s.PC += 4
			return
		}
		s.PC += 2

	// 0x4XNN: if (VX!=NN)
	// SNE X NN (skip if not equals)
	case op.Instruction() == 0x4:
		v, cmp := op.ExtractVNN()
		if cmp != s.V[v] {
			s.PC += 4
			return
		}
		s.PC += 2

	// 0x5XY0: if (VX==VY)
	// SRE X Y (skip if register equals)
	case op.Instruction() == 0x5:
		v1 := uint8((uint16(op) & 0x0F00) >> 8)
		v2 := uint8(uint16(op) & 0x00FF)
		if v1 == v2 {
			s.PC += 4
			return
		}
		s.PC += 2

	// 0x6XNN: set VX = NN
	// SRG V NN
	case op.Instruction() == 0x6:
		v, val := op.ExtractVNN()
		s.V[v] = val
		s.PC += 2

	default:
		panic(fmt.Sprintf("unknown opcode %X", op))
	}
}
