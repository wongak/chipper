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

func TestExtractXY(t *testing.T) {
	op := OpCode(0x8AFC)
	x, y, o := op.ExtractXY()
	if x != 0xA {
		t.Errorf("expect x to be 0xA, got %X", x)
	}
	if y != 0xF {
		t.Errorf("expect y to be 0xF, got %X", y)
	}
	if o != 0xC {
		t.Errorf("expect o to be 0xC, got %X", o)
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
		0xA12 SNE 1 FF
		0xA14 JMP 0xC00
	*/
	s.Mem[0x200] = 0x60
	s.Mem[0x201] = 0xA1
	s.Mem[0x202] = 0x1A
	s.Mem[0x203] = 0x10
	s.Mem[0xC00] = 0x61
	s.Mem[0xC01] = 0xFF
	s.Mem[0xA10] = 0x50
	s.Mem[0xA11] = 0x10
	s.Mem[0xA12] = 0x41
	s.Mem[0xA13] = 0x00
	s.Mem[0xA14] = 0x1C

	s.executeOpcode()
	if s.V[0] != 0xA1 {
		t.Errorf("expect SRG 0 A1, got V0 %X", s.V[0])
		return
	}
	s.executeOpcode()
	if s.PC != 0xA10 {
		t.Errorf("expect PC 0xA10, got %X", s.PC)
		return
	}
	s.executeOpcode()
	if s.PC != 0xA12 {
		t.Errorf("expect PC 0xA12, got %X", s.PC)
		return
	}
	s.executeOpcode()
	if s.PC != 0xA14 {
		t.Errorf("expect PC 0xA14, got %X", s.PC)
		return
	}
	s.executeOpcode()
	if s.PC != 0xC00 {
		t.Errorf("expect PC 0xC00, got %X", s.PC)
		return
	}
	s.executeOpcode()
	if s.V[1] != 0xFF {
		t.Errorf("expect SRG 1 FF, got V1 %X", s.V[1])
		return
	}
}

func TestADR(t *testing.T) {
	s := NewSystem()

	s.Mem[0x200] = 0x6F // SRG F FF
	s.Mem[0x201] = 0xFF
	s.Mem[0x202] = 0x7F // ADR F 01
	s.Mem[0x203] = 0x01

	s.executeOpcode()
	s.executeOpcode()

	if s.V[0xF] != 0 {
		t.Errorf("expect ADR overflow (0), got %X", s.V[0xF])
	}
}
