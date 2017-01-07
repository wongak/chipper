package apparat

import (
	"strings"
	"testing"
)

func TestInstructionOpcode(t *testing.T) {
	op := OpCode(0x2A01)
	if op.Instruction() != 0x2 {
		t.Errorf("expect instruction bits to be 0x2, got %X", op.Instruction())
	}
}

func TestRETStackUnderflow(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Error("expect panic")
			return
		}
		if !strings.Contains(err.(string), "underflow") {
			t.Errorf("expect underflow panic, got %s", err.(string))
		}
	}()
	s := NewSystem()
	s.Mem[0x200] = 0x00
	s.Mem[0x201] = 0xEE

	s.executeOpcode()
}

func TestJMP(t *testing.T) {
	s := NewSystem()
	// program: JMP 0xA10
	// opcpde 0x1A10
	s.Mem[0x200] = 0x1A
	s.Mem[0x201] = 0x10

	s.executeOpcode()
	if s.PC != 0xA10 {
		t.Errorf("expect PC to be 0xA10, got %X", s.PC)
	}
}

func TestSIE(t *testing.T) {
	s := NewSystem()
	s.V[4] = 0xA1

	s.Mem[0x200] = 0x34
	s.Mem[0x201] = 0xA1
	s.Mem[0x204] = 0x34
	s.Mem[0x205] = 0xA2

	s.executeOpcode()
	if s.PC != 0x204 {
		t.Errorf("expect PC to be 0x204, got %X", s.PC)
		return
	}
	s.executeOpcode()
	if s.PC != 0x206 {
		t.Errorf("expect PC to be 0x206, got %X", s.PC)
	}
}

func TestProg(t *testing.T) {
	s := NewSystem()

	/*
		0x200 SRG 0 A1
		0x202 JMP 0xA10
		0xC00 SRG 1 FF
		0xA10 SRE 0 1
		0xA14 SNE 0 FF
		0xA16 JMP 0xC00
	*/
	s.Mem[0x200] = 0x61
	s.Mem[0x201] = 0xA1
	s.Mem[0x202] = 0x1A
	s.Mem[0x203] = 0x10
	s.Mem[0xC00] = 0x61
	s.Mem[0xC01] = 0xFF
	s.Mem[0xA10] = 0x50
	s.Mem[0xA11] = 0x10
	s.Mem[0xA14] = 0x40
	s.Mem[0xA15] = 0xFF
	s.Mem[0xA16] = 0x1C
}
