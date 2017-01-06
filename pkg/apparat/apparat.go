package apparat

import (
	"sync"
	"time"
)

const (
	clockSpeed = 60
)

type (
	// System describes the CHIP-8 system
	System struct {
		clockSpeed time.Duration
		// opcodes per tick to be run (clock speed)
		ops int
		m   *sync.RWMutex

		// controls
		Stop chan struct{}

		// V is the representation of the registers V0-VE
		V [16]byte
		// I is the index register
		I uint16
		// PC is the program counter
		PC uint16

		Stack *Stack

		Mem Memory

		Dsp *Display

		Timers *Timers

		Key Keys
	}
)

// NewSystem initializes a new system
func NewSystem() *System {
	return &System{
		clockSpeed: time.Second / clockSpeed,
		ops:        10,
		m:          &sync.RWMutex{},

		Stop: make(chan struct{}),

		PC:     0x200,
		Stack:  &Stack{},
		Mem:    NewMemory(),
		Dsp:    NewDisplay(),
		Timers: &Timers{},
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
	s.Dsp.d = [32]uint64{}
	s.Key = 0

	s.Timers.Delay = 0
	s.Timers.Sound = 0
	s.m.Unlock()
}

// SetSpeed sets the opcodes per second execution speed
func (s *System) SetSpeed(ops int) {
	s.m.Lock()
	s.ops = ops
	s.m.Unlock()
}

// Run runs the system
func (s *System) Run() {
	t := time.Now()
	var wait time.Duration
	for {
		s.m.RLock()
		// run n opcodes in tick
		for i := 0; i < s.ops; i++ {
			s.executeOpcode()
		}
		// wait for tick
		// everything, opcodes, drawing, tick and loop
		// should be finished in an almost constant
		// time
		wait = time.Since(t)
		if wait < s.clockSpeed {
			time.Sleep(s.clockSpeed - wait)
		}
		s.Timers.Tick()
		s.draw()
		s.m.RUnlock()
		select {
		case <-s.Stop:
			return
		default:
		}
		// update new time
		// we start measuring from the last cycle
		t = time.Now()
	}
}

func (s *System) executeOpcode() {
}

func (s *System) draw() {
}
