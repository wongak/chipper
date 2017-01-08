package main

import termbox "github.com/nsf/termbox-go"

func drawDebug() {
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
