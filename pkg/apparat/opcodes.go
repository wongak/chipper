package apparat

import (
	"fmt"
	"math"
)

type (
	// Executer can execute an operation
	Executer interface {
		Execute(s *System)
		String() string
	}

	// OpCode represents an opcode
	OpCode uint16

	opCLS struct {
		OpCode
	}
	opRET struct {
		OpCode
	}
	opJP struct {
		OpCode
	}
	opCALL struct {
		OpCode
	}
	opSE struct {
		OpCode
	}
	opSNE struct {
		OpCode
	}
	opLD struct {
		OpCode
	}
	opADD struct {
		OpCode
	}
	opOR struct {
		OpCode
	}
	opAND struct {
		OpCode
	}
	opXOR struct {
		OpCode
	}
	opSUB struct {
		OpCode
	}
	opSHR struct {
		OpCode
	}
	opSUBN struct {
		OpCode
	}
	opSHL struct {
		OpCode
	}
	opRND struct {
		OpCode
	}
	opDRW struct {
		OpCode
	}
	opSKP struct {
		OpCode
	}
	opSKNP struct {
		OpCode
	}
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

// ExtractXYN extracts 3 bytes for the arithmetic instructions
// 0x0XYO
func (o OpCode) ExtractXYN() (uint8, uint8, uint8) {
	x := uint8((uint16(o) & 0x0F00) >> 8)
	y := uint8((uint16(o) & 0x00F0) >> 4)
	n := uint8(uint16(o) & 0x000F)
	return x, y, n
}

// Executer returns an Executer for the opcode
func (o OpCode) Executer() (Executer, error) {
	switch true {
	// 0x00E0: Clear the screen
	// CLS
	case 0x00E0 == o:
		return opCLS{o}, nil
	// 0x00EE: return from a subroutine
	// ret
	case 0x00EE == o:
		return opRET{o}, nil
	// 0x1NNN: JUMP
	// JP 0xNNN
	case o.Instruction() == 0x1:
		return opJP{o}, nil
	// 0x2NNN: SUB
	// CALL 0xNNN
	case o.Instruction() == 0x2:
		return opCALL{o}, nil
	// 0x3XNN: if(Vx==NN)
	// SE Vx, NN (skip if equals)
	case o.Instruction() == 0x3:
		return opSE{o}, nil
	// 0x4XNN: if (Vx!=NN)
	// SNE Vx, NN (skip if not equals)
	case o.Instruction() == 0x4:
		return opSNE{o}, nil
	// 0x5XY0: if (Vx==Vy)
	// SE Vx, Vy (skip if register equals)
	case o.Instruction() == 0x5:
		return opSE{o}, nil
	// 0x6XNN: set Vx = NN
	// LD Vx, NN
	case o.Instruction() == 0x6:
		return opLD{o}, nil
	// 0x7XNN: Vx += NN
	// ADD Vx, NN
	case o.Instruction() == 0x7:
		return opADD{o}, nil
	case o.Instruction() == 0x8:
		_, _, op := o.ExtractXYN()
		switch op {
		// 0x8XY0: Vx=Vy
		// LD Vx, Vy
		case 0x0:
			return opLD{o}, nil
		// 0x8XY1: Vx=Vx|Vy
		// OR Vx, Vy
		case 0x1:
			return opOR{o}, nil
		// 0x8XY2: Vx=Vx&Vy
		// AND Vx, Vy
		case 0x2:
			return opAND{o}, nil
		// 0x8XY3: Vx=Vx^Vy
		// XOR Vx, Vy
		case 0x3:
			return opXOR{o}, nil
		// 0x8XY4: Vx += Vy
		// ADD Vx, Vy
		case 0x4:
			return opADD{o}, nil
		// 0x8XY5: Vx -= Vy
		// SUB Vx, Vy
		case 0x5:
			return opSUB{o}, nil
		// 0x8XY6: Vx >> 1
		// SHR Vx
		case 0x6:
			return opSHR{o}, nil
		// 0x8XY7: Vx = Vy - Vx
		// SUBN Vx, Vy
		case 0x7:
			return opSUBN{o}, nil
		// 0x8XYE: Vx << 1
		// SHL Vx
		case 0xE:
			return opSHL{o}, nil
		default:
			panic(fmt.Sprintf("unknown opcode %X", o))
		}
	// 0x9XY0: if Vx != Vy
	// SNE Vx, Vy
	case o.Instruction() == 0x9:
		return opSNE{o}, nil
	// 0xANNN: set I = NNN
	// LD I, NNN
	case o.Instruction() == 0xA:
		return opLD{o}, nil
	// 0xBNNN: jump V0 + NNN
	// JP V0, 0xNNN
	case o.Instruction() == 0xB:
		return opJP{o}, nil
	// 0xCXNN: rand Vx & NN
	// RND Vx 0xNN
	case o.Instruction() == 0xC:
		return opRND{o}, nil
	// 0xDXYN: draw(Vx, Vy, N)
	// DRW Vx Vy N
	case o.Instruction() == 0xD:
		return opDRW{o}, nil
	// the 0xEXNN group
	case o.Instruction() == 0xE:
		_, n := o.ExtractVNN()
		switch n {
		// 0xEX9E: if(key()==Vx)
		// SKP Vx
		case 0x9E:
			return opSKP{o}, nil
		// 0xEXA1: if(key()!=Vx)
		// skip.ne Vx, key
		case 0xA1:
			return opSKNP{o}, nil
		default:
			panic(fmt.Sprintf("unknown opcode %X", o))
		}
	// the 0xFXNN group
	case o.Instruction() == 0xF:
		_, n := o.ExtractVNN()
		switch n {
		// 0xFX07: Vx = get_delay()
		// LD Vx, DT
		case 0x07:
			return opLD{o}, nil
		// 0xFX0A: Vx = get_key()
		// LD Vx, K
		case 0x0A:
			return opLD{o}, nil
		// 0xFX15: delay_timer(Vx)
		// LD DT, Vx
		case 0x15:
			return opLD{o}, nil
		// 0xFX18: sound_timer(Vx)
		// LD ST, Vx
		case 0x18:
			return opLD{o}, nil
		// 0xFX1E: I += Vx
		// ADD I, Vx
		case 0x1E:
			return opADD{o}, nil
		// 0xFX29: I = sprite_addr[Vx]
		// LD F, Vx
		case 0x29:
			return opLD{o}, nil
		// 0xFX33: set_BCD(Vx)
		// LD B, Vx
		case 0x33:
			return opLD{o}, nil
		// 0xFX55: reg_dump(Vx,&I)
		// LD [I], Vx
		case 0x55:
			return opLD{o}, nil
		// 0xFX65: reg_load(Vx,&I)
		// LD Vx, [I]
		case 0x65:
			return opLD{o}, nil
		default:
			return nil, fmt.Errorf("unknown opcode %X", o)
		}
	default:
		return nil, fmt.Errorf("unknown opcode %X", o)
	}
}

func (o opCLS) Execute(s *System) {
	s.Dsp.Clear()
	s.PC += 2
}
func (o opCLS) String() string {
	return "CLS\n"
}

func (o opRET) Execute(s *System) {
	s.PC = s.Stack.Pop()
	s.PC += 2
}
func (o opRET) String() string {
	return "RET\n"
}

func (o opJP) Execute(s *System) {
	switch o.Instruction() {
	case 0x1:
		s.PC = o.ExtractAddr()
	case 0xB:
		s.PC = o.ExtractAddr()
		s.PC += uint16(s.V[0])
	default:
		panic("invalid JP")
	}
}
func (o opJP) String() string {
	switch o.Instruction() {
	case 0x1:
		return fmt.Sprintf("JP 0x%X\n", o.ExtractAddr())
	case 0xB:
		return fmt.Sprintf("JP V0, 0x%X\n", o.ExtractAddr())
	default:
		panic("invalid JP")
	}
}

func (o opCALL) Execute(s *System) {
	s.Stack.Push(s.PC)
	s.PC = o.ExtractAddr()
}
func (o opCALL) String() string {
	return fmt.Sprintf("CALL 0x%X\n", o.ExtractAddr())
}

func (o opSE) Execute(s *System) {
	switch o.Instruction() {
	case 0x3:
		v, cmp := o.ExtractVNN()
		if cmp == s.V[v] {
			s.PC += 4
			return
		}
		s.PC += 2
	case 0x5:
		x, y, _ := o.ExtractXYN()
		if s.V[x] == s.V[y] {
			s.PC += 4
			return
		}
		s.PC += 2
	default:
		panic("invalid SE")
	}
}
func (o opSE) String() string {
	switch o.Instruction() {
	case 0x3:
		v, cmp := o.ExtractVNN()
		return fmt.Sprintf("SE V%X, 0x%X\n", v, cmp)
	case 0x5:
		x, y, _ := o.ExtractXYN()
		return fmt.Sprintf("SE V%X, V%X\n", x, y)
	default:
		panic("invalid SE")
	}
}

func (o opSNE) Execute(s *System) {
	switch o.Instruction() {
	case 0x4:
		v, cmp := o.ExtractVNN()
		if cmp != s.V[v] {
			s.PC += 4
			return
		}
		s.PC += 2
	case 0x9:
		x, y, _ := o.ExtractXYN()
		if s.V[x] != s.V[y] {
			s.PC += 4
			return
		}
		s.PC += 2
	default:
		panic("invalid SNE")
	}
}
func (o opSNE) String() string {
	switch o.Instruction() {
	case 0x4:
		v, cmp := o.ExtractVNN()
		return fmt.Sprintf("SNE V%X, 0x%X\n", v, cmp)
	case 0x9:
		x, y, _ := o.ExtractXYN()
		return fmt.Sprintf("SNE V%X, V%X\n", x, y)
	default:
		panic("invalid SNE")
	}
}

func (o opLD) Execute(s *System) {
	switch o.Instruction() {
	case 0x6:
		v, val := o.ExtractVNN()
		s.V[v] = val
		s.PC += 2
	case 0x8:
		x, y, op := o.ExtractXYN()
		switch op {
		case 0x0:
			s.V[x] = s.V[y]
			s.PC += 2
		default:
			panic("invalid LD")
		}
	case 0xA:
		s.I = o.ExtractAddr()
		s.PC += 2
	case 0xF:
		x, n := o.ExtractVNN()
		switch n {
		case 0x07:
			s.V[x] = s.Timers.Delay
			s.PC += 2
		case 0x0A:
			if s.Key.HasState() {
				s.V[x] = s.Key.State()
				s.PC += 2
			}
		case 0x15:
			s.Timers.Delay = s.V[x]
			s.PC += 2
		case 0x18:
			s.Timers.Sound = s.V[x]
			s.PC += 2
		// LD F, Vx
		case 0x29:
			s.I = fontsetStartAddress + (fontSize * uint16(s.V[x]))
			s.PC += 2
		// LD B, Vx
		case 0x33:
			v := s.V[x]
			s.Mem[s.I] = v / 100
			s.Mem[s.I+1] = (v / 10) % 10
			s.Mem[s.I+2] = (v % 100) % 10
			s.PC += 2
		// LD [I], Vx
		case 0x55:
			copy(s.Mem[s.I:], s.V[:x+1])
			s.PC += 2
		// LD Vx, [I]
		case 0x65:
			copy(s.V[:x+1], s.Mem[s.I:])
			s.PC += 2
		default:
			panic("invalid LD")
		}
	default:
		panic("invalid LD")
	}
}
func (o opLD) String() string {
	switch o.Instruction() {
	case 0x6:
		v, val := o.ExtractVNN()
		return fmt.Sprintf("LD V%X, 0x%X\n", v, val)
	case 0x8:
		x, y, op := o.ExtractXYN()
		switch op {
		case 0x0:
			return fmt.Sprintf("LD V%X, V%X\n", x, y)
		default:
			panic("invalid LD")
		}
	case 0xA:
		addr := o.ExtractAddr()
		return fmt.Sprintf("LD I, 0x%X\n", addr)
	case 0xF:
		x, n := o.ExtractVNN()
		switch n {
		case 0x07:
			return fmt.Sprintf("LD V%X, DT\n", x)
		case 0x0A:
			return fmt.Sprintf("LD V%X, K\n", x)
		case 0x15:
			return fmt.Sprintf("LD DT, V%X\n", x)
		case 0x18:
			return fmt.Sprintf("LD ST, V%X\n", x)
		case 0x29:
			return fmt.Sprintf("LD F, V%X\n", x)
		case 0x33:
			return fmt.Sprintf("LD B, V%X\n", x)
		case 0x55:
			return fmt.Sprintf("LD [I], V%X\n", x)
		case 0x65:
			return fmt.Sprintf("LD V%X, [I]\n", x)
		default:
			panic("invalid LD")
		}
	default:
		panic("invalid LD")
	}
}

func (o opADD) Execute(s *System) {
	switch o.Instruction() {
	case 0x7:
		v, val := o.ExtractVNN()
		s.V[v] += val
		s.PC += 2
	case 0x8:
		x, y, _ := o.ExtractXYN()
		s.V[0xF] = 0
		if math.MaxUint8-s.V[y] < s.V[x] {
			s.V[0xF] = 1
		}
		s.V[x] += s.V[y]
		s.PC += 2
	case 0xF:
		x, _ := o.ExtractVNN()
		// undocumented feature, overflow
		// will set VF to 1
		if math.MaxUint16-uint16(s.V[x]) < s.I {
			s.V[0xF] = 1
		}
		s.I += uint16(s.V[x])
		s.PC += 2
	default:
		panic("invalid ADD")
	}
}
func (o opADD) String() string {
	switch o.Instruction() {
	case 0x7:
		v, val := o.ExtractVNN()
		return fmt.Sprintf("ADD V%X, 0x%X\n", v, val)
	case 0x8:
		x, y, _ := o.ExtractXYN()
		return fmt.Sprintf("ADD V%X, V%X\n", x, y)
	case 0xF:
		x, _ := o.ExtractVNN()
		return fmt.Sprintf("ADD I, V%X\n", x)
	default:
		panic("invalid ADD")
	}
}

func (o opOR) Execute(s *System) {
	x, y, _ := o.ExtractXYN()
	s.V[x] |= s.V[y]
	s.PC += 2
}
func (o opOR) String() string {
	x, y, _ := o.ExtractXYN()
	return fmt.Sprintf("OR V%X, V%X\n", x, y)
}

func (o opAND) Execute(s *System) {
	x, y, _ := o.ExtractXYN()
	s.V[x] &= s.V[y]
	s.PC += 2
}
func (o opAND) String() string {
	x, y, _ := o.ExtractXYN()
	return fmt.Sprintf("AND V%X, V%X\n", x, y)
}

func (o opXOR) Execute(s *System) {
	x, y, _ := o.ExtractXYN()
	s.V[x] ^= s.V[y]
	s.PC += 2
}
func (o opXOR) String() string {
	x, y, _ := o.ExtractXYN()
	return fmt.Sprintf("XOR V%X, V%X\n", x, y)
}

func (o opSUB) Execute(s *System) {
	x, y, _ := o.ExtractXYN()
	s.V[0xF] = 1
	if s.V[y] > s.V[x] {
		s.V[0xF] = 0
	}
	s.V[x] -= s.V[y]
	s.PC += 2
}
func (o opSUB) String() string {
	x, y, _ := o.ExtractXYN()
	return fmt.Sprintf("SUB V%X, V%X\n", x, y)
}

func (o opSHR) Execute(s *System) {
	x, _, _ := o.ExtractXYN()
	s.V[0xF] = s.V[x] & 0x01
	s.V[x] = s.V[x] >> 1
	s.PC += 2
}
func (o opSHR) String() string {
	x, _, _ := o.ExtractXYN()
	return fmt.Sprintf("SHR V%X\n", x)
}

func (o opSUBN) Execute(s *System) {
	x, y, _ := o.ExtractXYN()
	s.V[0xF] = 1
	if s.V[x] > s.V[y] {
		s.V[0xF] = 0
	}
	s.V[x] = s.V[y] - s.V[x]
	s.PC += 2
}
func (o opSUBN) String() string {
	x, y, _ := o.ExtractXYN()
	return fmt.Sprintf("SUBN V%X, V%X\n", x, y)
}

func (o opSHL) Execute(s *System) {
	x, _, _ := o.ExtractXYN()
	s.V[0xF] = (s.V[x] & 0x80) >> 7
	s.V[x] = s.V[x] << 1
	s.PC += 2
}
func (o opSHL) String() string {
	x, _, _ := o.ExtractXYN()
	return fmt.Sprintf("SHL V%X\n", x)
}

func (o opRND) Execute(s *System) {
	x, n := o.ExtractVNN()
	s.rndSource.Read(s.V[x : x+1])
	s.V[x] = s.V[x] & n
	s.PC += 2
}
func (o opRND) String() string {
	x, n := o.ExtractVNN()
	return fmt.Sprintf("RND V%X, 0x%X\n", x, n)
}

func (o opDRW) Execute(s *System) {
	x, y, n := o.ExtractXYN()
	sprite := s.Mem[s.I : (s.I+(uint16(n)*8))+1]
	s.V[0xF] = s.Dsp.Draw(s.V[x], s.V[y], n, sprite)
	s.PC += 2
}
func (o opDRW) String() string {
	x, y, n := o.ExtractXYN()
	return fmt.Sprintf("DRW V%X, V%X, 0x%X\n", x, y, n)
}

func (o opSKP) Execute(s *System) {
	x, _ := o.ExtractVNN()
	if s.Key.HasState() && s.V[x]&s.Key.State() != 0 {
		s.PC += 4
		return
	}
	s.PC += 2
}
func (o opSKP) String() string {
	x, _ := o.ExtractVNN()
	return fmt.Sprintf("SKP V%X\n", x)
}

func (o opSKNP) Execute(s *System) {
	x, _ := o.ExtractVNN()
	if s.Key.HasState() && s.V[x]&s.Key.State() == 0 {
		s.PC += 4
		return
	}
	s.PC += 2
}
func (o opSKNP) String() string {
	x, _ := o.ExtractVNN()
	return fmt.Sprintf("SKNP V%X\n", x)
}
