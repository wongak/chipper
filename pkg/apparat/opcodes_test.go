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
	x, y, o := op.ExtractXYN()
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

func TestRetStackUnderflow(t *testing.T) {
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

type testDsp struct {
	calledClear bool
}

func (t *testDsp) Clear() {
	t.calledClear = true
}
func (t *testDsp) Draw(x, y, h uint8, sprite []byte) uint8 {
	return 0
}
func (t *testDsp) Line(y uint8) uint64 {
	return 0
}
func (t *testDsp) Dump() string {
	return ""
}

func TestCLS(t *testing.T) {
	s := NewSystem()
	dsp := &testDsp{}
	s.Dsp = dsp
	s.Mem[0x200] = 0x00
	s.Mem[0x201] = 0xE0
	s.executeOpcode()
	if !dsp.calledClear {
		t.Error("expect clear to be called")
	}
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

func TestSkipIfEq(t *testing.T) {
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
		0x200 load V0 0xA1
		0x202 jump 0xA10
		0xC00 load V1 0xFF
		0xA10 skip.eq V0 V1
		0xA12 skip.ne V1 0xFF
		0xA14 jump 0xC00
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

func TestAdd(t *testing.T) {
	s := NewSystem()

	s.Mem[0x200] = 0x6F // load VF FF
	s.Mem[0x201] = 0xFF
	s.Mem[0x202] = 0x7F // add VF 0x01
	s.Mem[0x203] = 0x01

	s.executeOpcode()
	s.executeOpcode()

	if s.V[0xF] != 0 {
		t.Errorf("expect add overflow (0), got %X", s.V[0xF])
		return
	}

	s.Mem[0x204] = 0x6A // load VA FE
	s.Mem[0x205] = 0xFE
	s.Mem[0x206] = 0x6B // load VB 04
	s.Mem[0x207] = 0x04
	s.Mem[0x208] = 0x8A // add VA, VB
	s.Mem[0x209] = 0xB4
	for i := 0; i < 3; i++ {
		s.executeOpcode()
	}
	if s.V[0xA] != 0x02 {
		t.Errorf("expect VA to be 0x02, got %X", s.V[0xA])
		return
	}
	if s.V[0xF] != 1 {
		t.Errorf("expect carry to be 1, got %X", s.V[0xF])
		return
	}
}

func TestSub(t *testing.T) {
	s := NewSystem()

	s.Mem[0x200] = 0x6A // load VA 02
	s.Mem[0x201] = 0x02
	s.Mem[0x202] = 0x6C // load VC 03
	s.Mem[0x203] = 0x03
	s.Mem[0x204] = 0x8A // sub VA VC
	s.Mem[0x205] = 0xC5
	for i := 0; i < 3; i++ {
		s.executeOpcode()
	}
	if s.V[0xA] != 0xFF {
		t.Errorf("expect VA to be 0xFF (wraparound), got %X", s.V[0xA])
		return
	}
	if s.V[0xF] != 0 {
		t.Errorf("expect borrow flag to be 0, got %X", s.V[0xF])
		return
	}
}

func TestShl(t *testing.T) {
	s := NewSystem()

	copy(s.Mem[:], []byte{
		0x64, 0xA0, // load V4 A0 (1010 0000)
		0x84, 0x0E, // shl V4
	})
	s.PC = 0
	s.executeOpcode()
	s.executeOpcode()

	if s.V[0xF] != 1 {
		t.Errorf("expect most siginificant bit to be 1, got VF %X|", s.V[0xF])
		return
	}
	if s.V[4] != 0x40 {
		t.Errorf("expect shifted 0x40, got %X", s.V[4])
		return
	}
}

type rndMock struct {
	w []byte
}

func (r rndMock) Read(p []byte) (int, error) {
	return copy(p, r.w), nil
}

func TestRnd(t *testing.T) {
	s := NewSystem()
	s.rndSource = rndMock{
		w: []byte{0xA1},
	}
	s.Mem[0x200] = 0xC0 // 0xC0FF rnd V0 FF
	s.Mem[0x201] = 0xFF

	s.executeOpcode()
	if s.V[0] != 0xA1 {
		t.Errorf("expect rand source result in V0, got %X", s.V[1])
		return
	}
	for i := 1; i < 16; i++ {
		if s.V[i] != 0 {
			t.Errorf("expect V%d to be 0, got %X", i, s.V[i])
		}
	}
	s.Mem[0x202] = 0xC0 // 0xC0E0 rnd V0 E0
	s.Mem[0x203] = 0xE0
	// is A1 & E0 = A0
	s.executeOpcode()
	if s.V[0] != 0xA0 {
		t.Errorf("expect rand source result in V0 0xA0, got %X", s.V[0])
		return
	}
}

func TestLoadRestoreReg(t *testing.T) {
	s := NewSystem()
	// test data
	copy(s.Mem[0x500:], []byte{
		0x01, 0x02, 0x03, 0x04,
	})
	testProg := []byte{
		0xA5, 0x00, // load I 0x500
		0xF3, 0x65, // restore V3
		0x71, 0xF0, // add V1 F0
		0x82, 0x34, // add V2 V3
		0x6F, 0x0F, // load VF, 0F
		0xFF, 0x55, // save VF
	}
	copy(s.Mem[0x200:], testProg)
	s.executeOpcode()
	s.executeOpcode()
	for i := 0; i < 4; i++ {
		if s.V[i] != uint8(i+1) {
			t.Errorf("expect reg %d to be %d, got %X", i, i, s.V[i])
		}
	}
	for i := 0; i < 4; i++ {
		s.executeOpcode()
	}
	if s.Mem[0x501] != 0xF2 {
		t.Errorf("expect 0x501 to be 0xF2, got %X", s.Mem[0x501])
	}
	if s.Mem[0x502] != 0x07 {
		t.Errorf("expect 0x502 to be 0x07, got %X", s.Mem[0x502])
	}
	if s.Mem[0x50F] != 0x0F {
		t.Errorf("expect 0x50F to be 0x0F, got %X", s.Mem[0x50F])
		t.Logf("reg: %+v\n", s.V)
		t.Log("mem:\n" + s.Mem.Dump(0x200))
	}
}
