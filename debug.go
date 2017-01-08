package main

import termbox "github.com/nsf/termbox-go"

var (
	showListing bool
	memPos      = 0x200
)

func drawDebug() {
	if showListing {
		drawListing()
	} else {
		drawMem()
	}
}

func drawListing() {
}

func drawMem() {
	mem := s.MemDump()
	x := 65
	y := 1
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
