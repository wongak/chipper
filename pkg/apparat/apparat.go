package apparat

import (
	"os"
	"sync/atomic"
	"time"
)

type (
	// System describes the CHIP-8 system
	System struct {
		// V is the representation of the registers V0-VE
		V [16]byte
		// I is the index register
		I uint16
		// PC is the program counter
		PC uint16

		Stack *Stack

		Mem Memory

		Dsp Display

		Timers *Timers

		Key Keys
	}

	// Stack is the CHIP-8 stack with 16 levels
	Stack struct {
		s  [16]uint16
		sp byte
	}

	// Memory is the 4K CHIP-8 memory
	Memory [4096]byte

	// Display represents the display state
	Display [32]uint64

	// Keys represents the keyboard state
	Keys byte

	// Timers represents the CHIP-8 timers
	Timers struct {
		close chan struct{}
		t     *time.Ticker
		delay uint32
		sound uint32
	}

	// OpCode represents an opcode
	OpCode uint16
)

// NewSystem initializes a new system
func NewSystem() *System {
	return &System{
		PC:     0x200,
		Stack:  &Stack{},
		Timers: NewTimers(),
	}
}

// Reset resets the system
func (s *System) Reset() {
	s.V = [16]byte{}
	s.I = 0
	s.PC = 0x200
	s.Stack.Reset()

	s.Mem = Memory([4096]byte{})
	s.Dsp = [32]uint64{}
	s.Key = 0

	s.Timers.SetDelay(0)
	s.Timers.SetSound(0)
}

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

// NewTimers initializes and starts timers
func NewTimers() *Timers {
	t := &Timers{
		close: make(chan struct{}),
	}
	t.t = time.NewTicker(time.Second / 60)
	go t.run()
	return t
}

func (t *Timers) run() {
	for {
		select {
		case <-t.close:
			t.t.Stop()
			return
		case <-t.t.C:
			if !atomic.CompareAndSwapUint32(&t.delay, 0, 0) {
				atomic.AddUint32(&t.delay, ^uint32(0))
			}
			if !atomic.CompareAndSwapUint32(&t.sound, 0, 0) {
				atomic.AddUint32(&t.sound, ^uint32(0))
				os.Stdout.Write([]byte{7})
			}
		}
	}
}

// SetDelay sets the current delay value
func (t *Timers) SetDelay(a uint16) {
	atomic.StoreUint32(&t.delay, uint32(a))
}

// Delay returns the current delay timer value
func (t *Timers) Delay() uint16 {
	d := atomic.LoadUint32(&t.delay)
	return uint16(d)
}

// SetSound sets the sound counter
func (t *Timers) SetSound(a uint16) {
	atomic.StoreUint32(&t.sound, uint32(a))
}

// Sound returns the current sound counter value
func (t *Timers) Sound() uint16 {
	d := atomic.LoadUint32(&t.sound)
	return uint16(d)
}

// Stop stops the timers and releases all resources
func (t *Timers) Stop() {
	close(t.close)
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
