package apparat

import "sync"

type (
	// System describes the CHIP-8 system
	System struct {
		// opcodes per second to be run (clock speed)
		ops int
		m   *sync.RWMutex

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
)

// NewSystem initializes a new system
func NewSystem() *System {
	return &System{
		ops: 600,
		m:   &sync.RWMutex{},

		PC:     0x200,
		Stack:  &Stack{},
		Mem:    NewMemory(),
		Timers: NewTimers(),
	}
}

// Reset resets the system
func (s *System) Reset() {
	s.m.Lock()
	s.V = [16]byte{}
	s.I = 0
	s.PC = 0x200
	s.Stack.Reset()

	s.Mem = NewMemory()
	s.Dsp = [32]uint64{}
	s.Key = 0

	s.Timers.SetDelay(0)
	s.Timers.SetSound(0)
	s.m.Unlock()
}

// SetSpeed sets the opcodes per second execution speed
func (s *System) SetSpeed(ops int) {
	s.m.Lock()
	s.ops = ops
	s.m.Unlock()
}
