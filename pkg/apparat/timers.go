package apparat

import (
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type (
	// Timers represents the CHIP-8 timers
	Timers struct {
		m      *sync.RWMutex
		close  chan struct{}
		paused bool
		t      *time.Ticker
		delay  uint32
		sound  uint32
	}
)

// NewTimers initializes and starts timers
func NewTimers() *Timers {
	t := &Timers{
		m:     &sync.RWMutex{},
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
			t.m.RLock()
			if t.paused {
				t.m.RUnlock()
				continue
			}
			t.m.RUnlock()
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

// Pause pauses the timers
func (t *Timers) Pause() {
	t.m.Lock()
	t.paused = true
	t.m.Unlock()
}

// Resume resumes the paused timers
func (t *Timers) Resume() {
	t.m.Lock()
	t.paused = false
	t.m.Unlock()
}

// Stop stops the timers and releases all resources
func (t *Timers) Stop() {
	close(t.close)
}
