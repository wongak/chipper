package apparat

import (
	"strings"
	"testing"
)

func TestReset(t *testing.T) {
	s := NewSystem()
	s.PC = 0xFFA
	s.V[1] = 0xAA
	s.Mem[24] = 0xFF
	s.Dsp.(*Display).d[31] = 16
	s.Reset()
	for i := 0; i < 16; i++ {
		if s.V[i] != 0 {
			t.Errorf("expect V%d to be 0, got %d", i, s.V[i])
		}
	}
	for i := 0x200; i < 4096; i++ {
		if s.Mem[i] != 0 {
			t.Errorf("expect mem %X to be 0, got %d", i, s.Mem[i])
		}
	}
	for i := 0; i < 32; i++ {
		if s.Dsp.(*Display).d[i] != 0 {
			t.Errorf("expect display line %d to be 0, got %d", i, s.Dsp.(*Display).d[i])
		}
	}
	if s.PC != 0x200 {
		t.Errorf("expect PC to be 0x200, got %X", s.PC)
	}
}

func TestStackPanicsUnderflow(t *testing.T) {
	s := &Stack{}
	defer func() {
		err := recover()
		if err == nil {
			t.Error("expected panic")
		}
		if !strings.Contains(err.(string), "underflow") {
			t.Error("expect panic message stack underflow")
		}
	}()
	s.Pop()
}

func TestStackPanicsOverflow(t *testing.T) {
	s := &Stack{}
	var i int
	defer func() {
		err := recover()
		if err == nil {
			t.Error("expected panic")
		}
		if !strings.Contains(err.(string), "overflow") {
			t.Error("expect panic message stack overflow")
		}
		if i != 15 {
			t.Errorf("expect panic on last push, got i %d", i)
		}
	}()
	for i = 0; i < 16; i++ {
		s.Push(0xA00)
	}
}

func TestStack(t *testing.T) {
	s := &Stack{}
	s.Push(0xA10)
	s.Push(0xF12)
	s.Push(0x123)
	x := s.Pop()
	if x != 0x123 {
		t.Errorf("expect pop last 0x123")
	}
	s.Push(0x423)
	x = s.Pop()
	if x != 0x423 {
		t.Errorf("expect pop last 0x423")
	}
	s.Push(0x000)
	s.Push(0xFF0)
	expect := []uint16{
		0xFF0,
		0x000,
		0xF12,
		0xA10,
	}
	for _, e := range expect {
		x = s.Pop()
		if x != e {
			t.Errorf("expect %X, got %X", e, x)
		}
	}
}

func TestFetchOpCode(t *testing.T) {
	m := Memory{}
	m[0] = 0x42
	m[1] = 0xAE
	m[2] = 0xE1
	m[1024] = 0x00
	m[1025] = 0xA1

	type in struct {
		addr   uint16
		expect OpCode
	}
	cases := []in{
		{0, 0x42AE},
		{1, 0xAEE1},
		{1024, 0x00A1},
	}
	for _, c := range cases {
		op := m.FetchOpcode(c.addr)
		if op != c.expect {
			t.Errorf("expected opcode to be %X, got %X", c.expect, op)
			return
		}
	}
}
