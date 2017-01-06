package apparat

import "os"

type (
	// Timers represents the CHIP-8 timers
	Timers struct {
		Delay byte
		Sound byte
	}
)

// Tick simulates a tick of the timers
func (t *Timers) Tick() {
	if t.Delay > 0 {
		t.Delay--
	}
	if t.Sound > 0 {
		t.Sound--
		os.Stdout.Write([]byte{7})
	}
}
