package apparat

import (
	"fmt"
	"math"
)

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

// ExtractXY extracts 3 bytes for the arithmetic instructions
// 0x0XYO
func (o OpCode) ExtractXY() (uint8, uint8, uint8) {
	x := uint8((uint16(o) & 0x0F00) >> 8)
	y := uint8((uint16(o) & 0x00F0) >> 4)
	op := uint8(uint16(o) & 0x000F)
	return x, y, op
}

func (s *System) executeOpcode() {
	op := s.Mem.FetchOpcode(s.PC)
	switch true {
	// 0x00E0: Clear the screen
	// clear
	case 0x00E0 == op:

	// 0x00EE: return from a subroutine
	// ret
	case 0x00EE == op:
		s.PC = s.Stack.Pop()

	// 0x1NNN: JUMP
	// jump 0xNNN
	case op.Instruction() == 0x1:
		s.PC = op.ExtractAddr()

	// 0x2NNN: SUB
	// call 0xNNN
	case op.Instruction() == 0x2:
		s.Stack.Push(s.PC)
		s.PC = op.ExtractAddr()

	// 0x3XNN: if(Vx==NN)
	// skip.eq x NN (skip if equals)
	case op.Instruction() == 0x3:
		v, cmp := op.ExtractVNN()
		if cmp == s.V[v] {
			s.PC += 4
			return
		}
		s.PC += 2

	// 0x4XNN: if (Vx!=NN)
	// skip.ne x NN (skip if not equals)
	case op.Instruction() == 0x4:
		v, cmp := op.ExtractVNN()
		if cmp != s.V[v] {
			s.PC += 4
			return
		}
		s.PC += 2

	// 0x5XY0: if (Vx==Vy)
	// skip.eq x y (skip if register equals)
	case op.Instruction() == 0x5:
		v1 := uint8((uint16(op) & 0x0F00) >> 8)
		v2 := uint8(uint16(op) & 0x00FF)
		if v1 == v2 {
			s.PC += 4
			return
		}
		s.PC += 2

	// 0x6XNN: set Vx = NN
	// load Vx NN
	case op.Instruction() == 0x6:
		v, val := op.ExtractVNN()
		s.V[v] = val
		s.PC += 2

	// 0x7XNN: Vx += NN
	// ADD Vx NN
	case op.Instruction() == 0x7:
		v, val := op.ExtractVNN()
		s.V[v] += val
		s.PC += 2

	case op.Instruction() == 0x8:
		x, y, o := op.ExtractXY()
		switch o {
		// 0x8XY0: Vx=Vy
		// load Vx, Vy
		case 0x0:
			s.V[x] = s.V[y]
			s.PC += 2

		// 0x8XY1: Vx=Vx|Vy
		// or Vx, Vy
		case 0x1:
			s.V[x] |= s.V[y]
			s.PC += 2

		// 0x8XY2: Vx=Vx&Vy
		// and Vx, Vy
		case 0x2:
			s.V[x] &= s.V[y]
			s.PC += 2

		// 0x8XY3: Vx=Vx^Vy
		// xor Vx, Vy
		case 0x3:
			s.V[x] ^= s.V[y]
			s.PC += 2

		// 0x8XY4: Vx += Vy
		// add Vx, Vy
		case 0x4:
			s.V[0xF] = 0
			if math.MaxUint8-s.V[y] < s.V[x] {
				s.V[0xF] = 1
			}
			s.V[x] += s.V[y]
			s.PC += 2

		// 0x8XY5: Vx -= Vy
		// sub Vx, Vy
		case 0x5:
			s.V[0xF] = 1
			if s.V[y] > s.V[x] {
				s.V[0xF] = 0
			}
			s.V[x] -= s.V[y]
			s.PC += 2

		// 0x8XY6: Vx >> 1
		// shr Vx
		case 0x6:
			s.V[0xF] = s.V[x] & 0x01
			s.V[x] = s.V[x] >> 1
			s.PC += 2

		// 0x8XY7: Vx = Vy - Vx
		// dif Vx Vy
		case 0x7:
			s.V[0xF] = 1
			if s.V[x] > s.V[y] {
				s.V[0xF] = 0
			}
			s.V[x] = s.V[y] - s.V[x]

		// 0x8XYE: Vx << 1
		// shl Vx
		case 0xE:
			s.V[0xF] = (s.V[x] & 0x80) >> 7
			s.V[x] = s.V[x] << 1
		}

	default:
		panic(fmt.Sprintf("unknown opcode %X", op))
	}
}
