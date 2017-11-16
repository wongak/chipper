package main

/*
import (
	"bytes"
	"fmt"

	termbox "github.com/nsf/termbox-go"
)

var (
	showMem bool
	memPos  uint16 = 0x200
)

func drawDebug() {
	if showMem {
		drawMem()
	} else {
		drawListing()
	}
}

func drawListing() {
	w, h := termbox.Size()
	x := 35

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	s.RWM.RLock()
	pc := s.PC
	v := make([]byte, 16)
	copy(v, s.V[:])
	vi := s.I
	s.RWM.RUnlock()
	buf.WriteString(fmt.Sprintf("PC: 0x%X", pc))
	for i, va := range v {
		buf.WriteString(fmt.Sprintf(" V%X:0x%X", i, va))
	}
	buf.WriteString(fmt.Sprintf(" I: 0x%X", vi))
	for i, c := range buf.String() {
		termbox.SetCell(x+i, 0, c, termbox.ColorWhite, termbox.ColorBlack)
	}

	buf.Reset()
	s.RWM.RLock()
	for i := uint16(0); i < 80; i++ {
		o := s.Mem.FetchOpcode(s.PC + i*2)
		exec, err := o.Executer()
		if err != nil {
			break
		}
		buf.WriteString(fmt.Sprintf("0x%X ", s.PC+i*2))
		buf.WriteString(exec.String())
	}
	s.RWM.RUnlock()
	x = 65
	y := 3
	for _, c := range buf.String() {
		if c == '\n' {
			if x < w {
				for i := x; i < w; i++ {
					termbox.SetCell(i, y, ' ', termbox.ColorWhite, termbox.ColorDefault)
				}
			}
			y++
			x = 65
			continue
		}
		termbox.SetCell(x, y, c, termbox.ColorWhite, termbox.ColorDefault)
		x++
	}
	if y < h {
		for i := y; i < h; i++ {
			for j := x; j < w; j++ {
				termbox.SetCell(j, i, ' ', termbox.ColorWhite, termbox.ColorDefault)
			}
		}
	}
}

func drawMem() {
	mem := s.MemDump(memPos)
	x := 65
	y := 2
	w := fmt.Sprintf("Mem 0x%X", memPos)
	for i, c := range w {
		termbox.SetCell(x+i, 1, c, termbox.ColorWhite, termbox.ColorDefault)
	}
	for _, c := range mem {
		if c == '\n' {
			y++
			x = 65
			continue
		}
		termbox.SetCell(x, y, c, termbox.ColorWhite, termbox.ColorDefault)
		x++
	}
}
*/
