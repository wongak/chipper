package apparat

import "testing"

func TestSpriteFromFontset(t *testing.T) {
	s := NewSystem()

	// 0 sprite
	sprite := s.Mem[0x50 : (0x50+(5*8))+1]
	expect := []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0,
	}
	for i := 0; i < 5; i++ {
		if sprite[i] != expect[i] {
			t.Errorf("expect sprite of 0 at byte %d, got %+v", i, sprite)
		}
	}
	sp := uint64(0xF0 << (56 - 28))
	if sp != 0xF00000000 {
		t.Errorf("line sprite error, got %X", sp)
	}
	sp = uint64(0XFF << (56 - 0))
	if sp != 0xFF00000000000000 {
		t.Errorf("line sprite 2, got %X", sp)
	}
	new := 0 ^ sp
	if new != sp {
		t.Errorf("expect XORd sp, got %X", new)
	}
}

func TestDsp(t *testing.T) {
	dsp := NewDisplay()
	flipped := dsp.draw(0, 0, 1, []byte{0xFF})
	if flipped != 0 {
		t.Error("expected no flip")
	}
	l := dsp.Line(0)
	if l != 0xFF00000000000000 {
		t.Errorf("expect drawing on line, got %X", l)
	}
}

func TestDrawOp(t *testing.T) {
	s := NewSystem()

	testProg := []byte{
		0xA0, 0x50, // load I 0x050
		// // draw 28, 13. 5
		0x60, 0x1C, // load V0 28 = 0x1C
		0x61, 0x0D, // load V1 13 = 0x0D
		0xD0, 0x15, // draw V0 V1 5
	}
	// ****
	// *  *
	// *  *
	// *  *
	// ****
	copy(s.Mem[0x200:], testProg)
	for i := 0; i < 4; i++ {
		s.executeOpcode()
	}
	if s.PC != 0x208 {
		t.Errorf("prog not run, got PC %X", s.PC)
		return
	}
	for i := uint8(0); i < 13; i++ {
		if s.Dsp.Line(i) != 0 {
			t.Errorf("expect blank line %d, got %X", i, s.Dsp.Line(i))
			t.Logf("regs:\n%+v\n", s.V)
			t.Log("Dsp:\n" + s.Dsp.Dump() + "\n")
			return
		}
	}
	if s.Dsp.Line(13) != 0xF00000000 {
		t.Errorf("expect line 1 to be drawn, got %X", s.Dsp.Line(13))
		t.Log("Dsp:\n" + s.Dsp.Dump() + "\n")
		return
	}
	if s.Dsp.Line(14) != 0x900000000 {
		t.Errorf("expect line 2 to be drawn, got %X", s.Dsp.Line(14))
		t.Log("Dsp:\n" + s.Dsp.Dump() + "\n")
		return
	}
	if s.Dsp.Line(15) != 0x900000000 {
		t.Errorf("expect line 3 to be drawn, got %X", s.Dsp.Line(15))
		t.Log("Dsp:\n" + s.Dsp.Dump() + "\n")
		return
	}
	if s.Dsp.Line(16) != 0x900000000 {
		t.Errorf("expect line 4 to be drawn, got %X", s.Dsp.Line(16))
		t.Log("Dsp:\n" + s.Dsp.Dump() + "\n")
		return
	}
	if s.Dsp.Line(17) != 0xF00000000 {
		t.Errorf("expect line 5 to be drawn, got %X", s.Dsp.Line(17))
		t.Log("Dsp:\n" + s.Dsp.Dump() + "\n")
		return
	}
	// ****
	// *  *
	// ****
	// *  *
	// ****
	copy(s.Mem[0x208:], []byte{
		0xA0, 0x78, // load I 0x78 = 8
		0xD0, 0x15, // draw V0, V1, 5
	})
	s.executeOpcode()
	s.executeOpcode()
	if s.V[0xF] != 1 {
		t.Error("expect flipped to be 1")
	}
	if s.Dsp.Line(15) != 0x600000000 {
		t.Errorf("expect line 3 to be drawn, got %X", s.Dsp.Line(15))
		t.Log("Dsp:\n" + s.Dsp.Dump() + "\n")
	}
}
