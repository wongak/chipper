package main

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
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	o := s.Mem.FetchOpcode(s.PC)
	buf.WriteString(o.Executer().String())
	x := 65
	y := 2
	for _, c := range buf.String() {
		if c == '\n' {
			y++
			x = 65
			continue
		}
		termbox.SetCell(x, y, c, termbox.ColorWhite, termbox.ColorDefault)
		x++
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
