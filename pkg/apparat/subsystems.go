package apparat

import (
	"encoding/binary"
	"encoding/hex"
	"sync"
)

type (
	// Stack is the CHIP-8 stack with 16 levels
	Stack struct {
		s  [16]uint16
		sp byte
	}

	// Memory is the 4K CHIP-8 memory
	Memory [4096]byte

	// Displayer represents a display
	Displayer interface {
		Draw(x, y, h uint8, sprite []byte) uint8
		Clear()
		Line(y uint8) uint64
		Dump() string
	}

	// Display represents the display state
	Display struct {
		RWM *sync.RWMutex
		d   [32]uint64
	}

	// KeyStater can report key states
	// If k & 0xF0 != 0, then report true on any key
	// Otherwise report true if key 0x0F & k is pressed.
	KeyStater interface {
		KeyPressed(k uint8) bool
		HasState() bool
		State() uint8
	}

	defaultKeys struct{}
)

// Reset resets the stack
func (s *Stack) Reset() {
	s.s = [16]uint16{}
	s.sp = 0
}

// Push pushes onto the stack
func (s *Stack) Push(a uint16) {
	if s.sp > 14 {
		panic("stack overflow")
	}
	s.s[s.sp] = a
	s.sp++
}

// Pop pops from the stack
func (s *Stack) Pop() uint16 {
	if s.sp == 0 {
		panic("stack underflow")
	}
	s.sp--
	v := s.s[s.sp]
	return v
}

// NewMemory creates a memory with the preloaded CHIP-8 data
func NewMemory() Memory {
	m := [4096]byte{}
	for i, f := range fontset {
		m[fontsetStartAddress+uint16(i)] = f
	}
	return Memory(m)
}

// FetchOpcode retrieves an opcode from the given address
//
// An opcode consists of 2 bytes, therefore we need to retrieve
// the bytes from PC and PC+1
func (m *Memory) FetchOpcode(addr uint16) OpCode {
	var opcode uint16
	opcode = uint16(m[addr]) << 8
	opcode = opcode | uint16(m[addr+1])
	return OpCode(opcode)
}

// Dump dumps the memory as a hexdump
func (m *Memory) Dump(addr uint16) string {
	return hex.Dump(m[addr:])
}

// NewDisplay creates a new display
//
// It is implemented as a memory map, which can be read
// by other processes to perform the actual drawing
func NewDisplay() *Display {
	return &Display{
		RWM: &sync.RWMutex{},
		d:   [32]uint64{},
	}
}

func (d *Display) Draw(x, y, h uint8, sprite []byte) uint8 {
	d.RWM.Lock()
	defer d.RWM.Unlock()
	var flipped uint8
	var l uint8
	for l = 0; l < h; l++ {
		if l+y >= DisplayHeight {
			break
		}
		// get line bitmap
		m := d.d[l+y]
		// shift sprite to xa position
		shft := 56 - int(x)
		var sp uint64
		if shft < 0 {
			sp = uint64(sprite[l]) >> uint(shft*-1)
		} else {
			sp = uint64(sprite[l]) << (56 - x)
		}
		// XOR bitmap and sprite
		d.d[l+y] = m ^ sp
		// AND bitmap and sprite will give us any
		// flips
		if flipped == 0 && (m&sp != 0) {
			flipped = 0x1
		}
	}
	return flipped
}

// Line returns the bitmap for line y
func (d *Display) Line(y uint8) uint64 {
	d.RWM.RLock()
	l := d.d[y]
	d.RWM.RUnlock()
	return l

}

// Clear clearst the display
func (d *Display) Clear() {
	d.RWM.Lock()
	for i := 0; i < 32; i++ {
		d.d[i] = 0
	}
	d.RWM.Unlock()
}

// Dump dumps the display buffer in the same format
// as hexdump
func (d *Display) Dump() string {
	buf := make([]byte, 32*8)
	d.RWM.RLock()
	for i, r := range d.d {
		binary.BigEndian.PutUint64(buf[i*8:], r)
	}
	d.RWM.RUnlock()
	return hex.Dump(buf)
}

func (k defaultKeys) KeyPressed(_ uint8) bool {
	return false
}

func (k defaultKeys) HasState() bool {
	return false
}

func (k defaultKeys) State() uint8 {
	return 0
}
