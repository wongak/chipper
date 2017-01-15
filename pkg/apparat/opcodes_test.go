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

func TestCALLRET(t *testing.T) {
	s := NewSystem()

	s.Mem[0x200] = 0x2A // CALL 0xA22
	s.Mem[0x201] = 0x22
	s.Mem[0x202] = 0x7E // ADD VE 1
	s.Mem[0x203] = 0x01
	s.Mem[0xA22] = 0x6E // LD VE 23
	s.Mem[0xA23] = 0x23
	s.Mem[0xA24] = 0x00 // RET
	s.Mem[0xA25] = 0xEE

	s.executeOpcode()
	if s.PC != 0xA22 {
		t.Errorf("expect PC to be 0xA22, got %X", s.PC)
		return
	}
	if s.V[0xE] != 0 {
		t.Error("VE init err")
		return
	}
	s.executeOpcode()
	if s.V[0xE] != 0x23 {
		t.Error("invalid value in VE")
		return
	}
	s.executeOpcode()
	if s.PC != 0x202 {
		t.Error("RET call error")
		return
	}
	s.executeOpcode()
	if s.V[0xE] != 0x24 {
		t.Error("finish error")
	}
}

func TestJP(t *testing.T) {
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

	// SE V4, A1
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
		0x200 LD V0, 0xA1
		0x202 JP 0xA10
		0xC00 LD V1, 0xFF
		0xA10 SE V0, V1
		0xA12 SNE V1, 0x00
		0xA14 JP 0xC00
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
		t.Errorf("expect V0 A1, got V0 %X", s.V[0])
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
		t.Errorf("expect V1 FF, got V1 %X", s.V[1])
		return
	}
}

func TestSERegisters(t *testing.T) {
	s := NewSystem()
	s.Mem[0x200] = 0x50 // SE V0, V1
	s.Mem[0x201] = 0xA0

	s.executeOpcode()
	if s.PC != 0x204 {
		t.Error("SE V0, V1 not executed")
	}
}

func TestADD(t *testing.T) {
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

func TestLDRegisters(t *testing.T) {
	s := NewSystem()
	s.V[0xA] = 0xAE

	s.Mem[0x200] = 0x83 // LD V3, VA
	s.Mem[0x201] = 0xA0
	s.executeOpcode()
	if s.V[3] != 0xAE {
		t.Error("failed LD V3, VA")
	}
}

func TestBinaryOps(t *testing.T) {
	s := NewSystem()
	s.V[0] = 1
	s.V[1] = 4
	s.V[2] = 1
	s.V[3] = 241
	s.V[4] = 0xFF

	s.Mem[0x200] = 0x80 // OR V0, V1
	s.Mem[0x201] = 0x11
	s.Mem[0x202] = 0x80 // OR V0, V2
	s.Mem[0x203] = 0x21
	s.Mem[0x204] = 0x80 // AND V0, V2
	s.Mem[0x205] = 0x22
	s.Mem[0x206] = 0x83 // XOR V3, V4
	s.Mem[0x207] = 0x43

	s.executeOpcode()
	if s.V[0] != 5 {
		t.Error("4 OR 1 error")
		return
	}
	s.executeOpcode()
	if s.V[0] != 5 {
		t.Error("5 OR 1 error")
		return
	}
	s.executeOpcode()
	if s.V[0] != 1 {
		t.Error("5 AND 1 error")
		return
	}
	s.executeOpcode()
	if s.V[3] != 0xE {
		t.Error("F1 XOR FF error")
		return
	}
}

func TestADDCarry(t *testing.T) {
	s := NewSystem()
	s.V[0] = 0xFE
	s.V[1] = 0x12

	s.Mem[0x200] = 0x80 // ADD V0, V1
	s.Mem[0x201] = 0x14
	s.Mem[0x202] = 0x80
	s.Mem[0x203] = 0x14

	s.executeOpcode()
	if s.V[0] != 0x10 {
		t.Error("ADD FE 12 error")
		return
	}
	if s.V[0xF] != 1 {
		t.Error("carry flag not set")
		return
	}
	s.executeOpcode()
	if s.V[0] != 0x22 {
		t.Error("ADD 10 12 error")
		return
	}
	if s.V[0xF] != 0 {
		t.Error("carry flag not unset")
		return
	}
}

func TestSUB(t *testing.T) {
	s := NewSystem()

	s.Mem[0x200] = 0x6A // LD VA, 02
	s.Mem[0x201] = 0x02
	s.Mem[0x202] = 0x6C // LD VC, 03
	s.Mem[0x203] = 0x03
	s.Mem[0x204] = 0x8A // SUB VA, VC
	s.Mem[0x205] = 0xC5
	s.Mem[0x206] = 0x8A
	s.Mem[0x207] = 0xC5
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
	s.executeOpcode()
	if s.V[0xA] != 0xFC {
		t.Error("SUB FF 03 error")
		return
	}
	if s.V[0xF] != 1 {
		t.Error("NOT borrow should be 1")
		return
	}
}

func TestSHR(t *testing.T) {
	s := NewSystem()
	s.V[3] = 13
	s.V[6] = 2

	s.Mem[0x200] = 0x83 // SHR V3
	s.Mem[0x201] = 0x06
	s.Mem[0x202] = 0x86 // SHR V6
	s.Mem[0x203] = 0x06

	s.executeOpcode()
	if s.V[0xF] != 1 {
		t.Error("LSB error")
		return
	}
	if s.V[3] != 6 {
		t.Error("SHR not executed")
		return
	}
	s.executeOpcode()
	if s.V[0x0] != 0 {
		t.Error("2 LSB error")
		return
	}
	if s.V[6] != 1 {
		t.Error("SHR 2 = 1 error")
		return
	}
}

func TestSUBN(t *testing.T) {
	s := NewSystem()

	s.V[0] = 0x04
	s.V[1] = 0x01
	s.V[2] = 0x02
	s.V[3] = 0x0C
	copy(s.Mem[0x200:], []byte{
		0x80, 0x17, // SUBN V0, V1
		0x82, 0x37, // SUBN V2, V3
	})
	s.executeOpcode()
	if s.V[0] != 0xFD {
		t.Error("01 - 04 error")
		return
	}
	if s.V[0xF] != 0 {
		t.Error("wraparound flag error")
		return
	}
	s.executeOpcode()
	if s.V[2] != 0x0A {
		t.Error("0C - 02 error")
		return
	}
	if s.V[0xF] != 1 {
		t.Error("not borrow error")
		return
	}
}

func TestSHL(t *testing.T) {
	s := NewSystem()

	copy(s.Mem[:], []byte{
		0x64, 0xA0, // LD V4, A0 (1010 0000)
		0x84, 0x0E, // SHL V4
		0x62, 0x10, // LD V2, 10 (0001 0000)
		0x82, 0x0E, // SHL V2
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
	s.executeOpcode()
	s.executeOpcode()
	if s.V[0xF] != 0 {
		t.Error("MSB error")
		return
	}
	if s.V[2] != 32 {
		t.Error("SHL 16 error")
		return
	}
}

func TestSNERegs(t *testing.T) {
	s := NewSystem()

	copy(s.Mem[0x200:], []byte{
		0x90, 0x10, // SNE V0, V1
		0x71, 0x05, // LD V1, 0x05
		0x90, 0x10, // SNE V0, V1
		0x75, 0x01, // LD V5, 0x01
		0x75, 0x02, // LD V5, 0x02
	})
	s.executeOpcode()
	s.executeOpcode()
	if s.V[1] != 5 {
		t.Error("first SNE error")
		return
	}
	s.executeOpcode()
	s.executeOpcode()
	if s.V[5] != 2 {
		t.Error("second SNE error")
		return
	}
}

func TestLDI(t *testing.T) {
	s := NewSystem()

	s.Mem[0x200] = 0xA1
	s.Mem[0x201] = 0x23
	s.executeOpcode()
	if s.I != 0x123 {
		t.Error("error on LD I")
		return
	}
}

func TestJPV0NNN(t *testing.T) {
	s := NewSystem()

	s.V[0] = 5
	s.Mem[0x200] = 0xBA // JP V0, A08
	s.Mem[0x201] = 0x08
	s.executeOpcode()
	if s.PC != 0xA0D {
		t.Error("JP 5 + A08 error")
		return
	}
}

type rndMock struct {
	w []byte
}

func (r rndMock) Read(p []byte) (int, error) {
	return copy(p, r.w), nil
}

func TestRND(t *testing.T) {
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

type mockKeyState struct {
	ks uint8
}

func (k mockKeyState) KeyPressed(p uint8) bool {
	return k.ks&0x0F == p
}

func (k mockKeyState) HasState() bool {
	return k.ks&0xF0 != 0
}

func (k mockKeyState) State() uint8 {
	return k.ks & 0x0F
}

func TestSKP(t *testing.T) {
	s := NewSystem()

	s.Key = &mockKeyState{ks: 0x1E}
	if !s.Key.HasState() {
		t.Error("has state err")
		return
	}
	if !s.Key.KeyPressed(0x0E) {
		t.Error("pressed err")
		return
	}
	s.V[1] = 0x0E
	s.V[2] = 0x01
	s.Mem[0x200] = 0xE1 // SKP V1
	s.Mem[0x201] = 0x9E
	s.Mem[0x204] = 0xA1 // LD I, 0x123
	s.Mem[0x205] = 0x23
	s.Mem[0x206] = 0xE2 // SKP V2
	s.Mem[0x207] = 0x9E
	s.Mem[0x208] = 0xA5 // LD I, 0x543
	s.Mem[0x209] = 0x43
	s.executeOpcode()
	if s.PC != 0x204 {
		t.Error("SKP V1 PC err")
		return
	}
	s.executeOpcode()
	if s.I != 0x123 {
		t.Error("I reg error")
		return
	}
	s.executeOpcode()
	if s.PC != 0x208 {
		t.Error("SKP V2 PC err")
		return
	}
	s.executeOpcode()
	if s.I != 0x543 {
		t.Error("SKP no match error")
		return
	}
}

func TestSKNP(t *testing.T) {
	s := NewSystem()
	s.Key = &mockKeyState{ks: 0x1A}

	s.V[4] = 0x0A
	s.V[5] = 0x0B
	s.Mem[0x200] = 0xE4 // SKNP V4
	s.Mem[0x201] = 0xA1
	s.Mem[0x202] = 0xE5 // SKNP V5
	s.Mem[0x203] = 0xA1
	s.executeOpcode()
	s.executeOpcode()
	if s.PC != 0x206 {
		t.Error("PC err")
		return
	}
}

func TestDT(t *testing.T) {
	s := NewSystem()
	s.V[0] = 1
	s.V[2] = 2

	copy(s.Mem[0x200:], []byte{
		0xF0, 0x15, // LD DT, V0
		0xF1, 0x07, // LD V1, DT
		0xF2, 0x07, // LD V2, DT
	})
	s.executeOpcode()
	s.executeOpcode()
	if s.V[1] != 1 {
		t.Error("no tick delay err V1")
		return
	}
	s.Timers.Tick()
	s.executeOpcode()
	if s.V[2] != 0 {
		t.Error("tick delay err V2")
		return
	}
}

func TestLoadK(t *testing.T) {
	s := NewSystem()
	s.Key = &mockKeyState{0x01}

	s.Mem[0x200] = 0xF0 // LD V0, K
	s.Mem[0x201] = 0x0A

	s.executeOpcode()
	if s.PC != 0x200 {
		t.Error("PC err if no key pressed")
		return
	}
	s.Key.(*mockKeyState).ks = 0x12
	s.executeOpcode()
	if s.PC != 0x202 {
		t.Error("PC err if key pressed")
		return
	}
	if s.V[0] != 2 {
		t.Error("wrong key state stored")
		return
	}
}

func TestADDIVX(t *testing.T) {
	s := NewSystem()

	copy(s.Mem[0x200:], []byte{
		0xA0, 0x42, // LD I, 0x042
		0x61, 0x08, // LD V1, 0x08
		0xF1, 0x1E, // ADD I, V1
	})
	s.executeOpcode()
	s.executeOpcode()
	s.executeOpcode()
	if s.I != 0x04A {
		t.Error("ADD I, V1 err")
		return
	}
}

func TestLoadF(t *testing.T) {
	s := NewSystem()

	copy(s.Mem[0x200:], []byte{
		0xF0, 0x29, // LD F, V0
		0x60, 0x01, // LD V0, 0x01
		0xF0, 0x29, // LD F, V0
	})
	s.executeOpcode()
	if s.I != fontsetStartAddress {
		t.Error("error font start")
		return
	}
	s.executeOpcode()
	s.executeOpcode()
	if s.I != fontsetStartAddress+fontSize {
		t.Error("error font +1")
		return
	}
}

func TestLoadB(t *testing.T) {
	s := NewSystem()

	copy(s.Mem[0x200:], []byte{
		0xA3, 0x00, // LD I, 0x300
		0x60, 0x7B, // LD V0, 0x7B
		0xF0, 0x33, // LD B, V0
	})
	for i := 0; i < 3; i++ {
		s.executeOpcode()
	}
	if s.Mem[0x300] != 1 {
		t.Error("expect decimal 1")
		return
	}
	if s.Mem[0x301] != 2 {
		t.Error("expect decimal 2")
		return
	}
	if s.Mem[0x302] != 3 {
		t.Error("expect decimal 3")
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
