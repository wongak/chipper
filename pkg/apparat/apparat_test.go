package apparat

import "testing"

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
