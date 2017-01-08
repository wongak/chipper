package apparat

import (
	"crypto/rand"
	"io"
	"sync"
	"time"
)

const (
	clockSpeed = 60

	// DisplayWidth is the width of the CHIP-8 display
	DisplayWidth = 64
	// DisplayHeight is the height of the CHIP-8 display
	DisplayHeight = 32
)

type (
	// System describes the CHIP-8 system
	System struct {
		clockSpeed time.Duration
		// opcodes per tick to be run (clock speed)
		ops int
		// controls
		RWM     *sync.RWMutex
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

		Key *Keys
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
		RWM:        &sync.RWMutex{},

		Stop: make(chan struct{}),
		Draw: func(_ *Display) {},

		rndSource: rand.Reader,

		PC:     0x200,
		Stack:  &Stack{},
		Mem:    NewMemory(),
		Dsp:    NewDisplay(),
		Timers: &Timers{},
		Key:    NewKeys(),
	}
}

// Reset resets the system
func (s *System) Reset() {
	s.RWM.Lock()
	s.V = [16]byte{}
	s.I = 0
	s.PC = 0x200
	s.Stack.Reset()

	s.Mem = NewMemory()
	s.Dsp.d = [32]uint64{}
	s.Key.Reset()

	s.Timers.Delay = 0
	s.Timers.Sound = 0
	s.RWM.Unlock()
}

// SetSpeed sets the opcodes per second execution speed
func (s *System) SetSpeed(ops int) {
	s.RWM.Lock()
	s.ops = ops
	s.RWM.Unlock()
}

// IsRunning returns true when the system is running
func (s *System) IsRunning() bool {
	s.RWM.RLock()
	running := s.running
	s.RWM.RUnlock()
	return running
}

// Step performs a single step
func (s *System) Step() {
	s.RWM.Lock()
	s.executeOpcode()
	s.Draw(s.Dsp)
	s.RWM.Unlock()
}

// Tick manually ticks the timer (currently ticking at
// around 60Hz)
// Tick is supposed to be used in conjunction with Step
// for debugging
func (s *System) Tick() {
	s.RWM.Lock()
	s.Timers.Tick()
	s.RWM.Unlock()
}

// LoadROM loads the contents of b into the execution space
// 0x200
func (s *System) LoadROM(b []byte) {
	s.RWM.Lock()
	copy(s.Mem[0x200:], b)
	s.RWM.Unlock()
}

// MemDump dumps the current mem in a hexdump format
func (s *System) MemDump(addr uint16) string {
	s.RWM.Lock()
	d := s.Mem.Dump(addr)
	s.RWM.Unlock()
	return d
}

func (s *System) executeOpcode() {
	op := s.Mem.FetchOpcode(s.PC)
	exec, err := op.Executer()
	if err != nil {
		panic(err)
	}
	exec.Execute(s)
}

// Run runs the system
func (s *System) Run() {
	s.RWM.Lock()
	s.running = true
	s.RWM.Unlock()
	t := time.Now()
	var wait time.Duration
	for {
		s.RWM.Lock()
		// run n opcodes in tick
		for i := 0; i < s.ops; i++ {
			s.executeOpcode()
		}
		s.Key.Reset()
		s.RWM.Unlock()
		// wait for tick
		// everything, opcodes, drawing, tick and loop
		// should be finished in an almost constant
		// time
		wait = time.Since(t)
		if wait < s.clockSpeed {
			time.Sleep(s.clockSpeed - wait)
		}
		s.RWM.Lock()
		s.Timers.Tick()
		s.RWM.Unlock()
		s.Draw(s.Dsp)
		select {
		case <-s.Stop:
			s.RWM.Lock()
			s.running = false
			s.RWM.Unlock()
			return
		default:
		}
		// update new time
		// we start measuring from the last cycle
		t = time.Now()
	}
}
