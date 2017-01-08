package apparat

import (
	"crypto/rand"
	"io"
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
		Stop    chan struct{}
		Draw    DrawCall
		running bool

		// emulated
		rndSource io.Reader

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

	// DrawCall is the connection to the actual implementation
	// of the display.
	// The func will receive a display, it should read lock the mutex,
	// then draw.
	// The DrawCall will be called each time the system wants
	// to update a frame. The display acts like the framebufer.
	// The call itself should not block.
	DrawCall func(dsp *Display)
)

// NewSystem initializes a new system
func NewSystem() *System {
	return &System{
		clockSpeed: time.Second / clockSpeed,
		ops:        10,
		m:          &sync.RWMutex{},

		Stop: make(chan struct{}),
		Draw: func(_ *Display) {},

		rndSource: rand.Reader,

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

// IsRunning returns true when the system is running
func (s *System) IsRunning() bool {
	s.m.RLock()
	running := s.running
	s.m.RUnlock()
	return running
}

// Step performs a single step
func (s *System) Step() {
	s.m.Lock()
	s.executeOpcode()
	s.Draw(s.Dsp)
	s.m.Unlock()
}

// Tick manually ticks the timer (currently ticking at
// around 60Hz)
// Tick is supposed to be used in conjunction with Step
// for debugging
func (s *System) Tick() {
	s.m.Lock()
	s.Timers.Tick()
	s.m.Unlock()
}

// Run runs the system
func (s *System) Run() {
	s.m.Lock()
	s.running = true
	s.m.Unlock()
	t := time.Now()
	var wait time.Duration
	for {
		s.m.Lock()
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
		s.Draw(s.Dsp)
		s.m.Unlock()
		select {
		case <-s.Stop:
			s.m.Lock()
			s.running = false
			s.m.Unlock()
			return
		default:
		}
		// update new time
		// we start measuring from the last cycle
		t = time.Now()
	}
}
