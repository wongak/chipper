package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	termbox "github.com/nsf/termbox-go"
	"github.com/wongak/chipper/pkg/apparat"
)

var (
	romFile string
	pause   bool
	debug   bool
)

var (
	s      *apparat.System
	paused bool
	stop   chan struct{}
)

func main() {
	flag.StringVar(&romFile, "rom", "", "ROM file to load")
	flag.BoolVar(&pause, "p", true, "whether to pause on start")
	flag.BoolVar(&debug, "d", true, "debug mode")
	flag.Parse()

	initSystem()
	err := termbox.Init()
	if err != nil {
		fmt.Printf("error initializing term: %v", err)
		os.Exit(1)
	}
	termbox.SetOutputMode(termbox.OutputNormal)
	defer termbox.Close()

	stop = make(chan struct{})

	status := statusBar{s}
	status.Draw()
	termbox.Flush()

	s.Draw = func(dsp apparat.Displayer) {
		go termbox.Interrupt()
	}

	go s.Run()

mainLoop:
	for {
		status.Draw()

		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventError:
			fmt.Printf("error event: %v", ev.Err)
			os.Exit(1)

		case termbox.EventInterrupt:
			for y := 0; y < apparat.DisplayHeight; y++ {
				l := s.Dsp.Line(uint8(y))
				for x := uint16(0); x < apparat.DisplayWidth; x++ {
					mask := uint64(1) << (63 - x)
					if mask&l != 0 {
						termbox.SetCell(int(x), y+1, ' ', termbox.ColorBlack, termbox.ColorWhite)
					} else {
						termbox.SetCell(int(x), y+1, ' ', termbox.ColorWhite, termbox.ColorBlack)
					}
				}
			}
			if debug {
				drawDebug()
			}
			termbox.Flush()

		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlC:
				break mainLoop

			case termbox.KeySpace:
				if paused {
					s.Step()
				}

			default:
				if ev.Ch != 0 {
					switch ev.Ch {
					case 'p':
						if paused {
							s.SetSpeed(10)
						} else {
							s.SetSpeed(0)
						}
						paused = !paused
					case 'l':
						showMem = !showMem

					case '1':
						s.Key.SetState(1)
					case '2':
						s.Key.SetState(2)
					case '3':
						s.Key.SetState(3)
					case 'q':
						s.Key.SetState(4)
					case 'w':
						s.Key.SetState(5)
					case 'e':
						s.Key.SetState(6)
					case 'a':
						s.Key.SetState(7)
					case 's':
						s.Key.SetState(8)
					case 'd':
						s.Key.SetState(9)
					case 'z', 'y':
						s.Key.SetState(0xA)
					case 'x':
						s.Key.SetState(0)
					case 'c':
						s.Key.SetState(0xB)
					case '4':
						s.Key.SetState(0xC)
					case 'r':
						s.Key.SetState(0xD)
					case 'f':
						s.Key.SetState(0xE)
					case 'v':
						s.Key.SetState(0xF)
					}
				}
			}
		}
	}
}

func initSystem() {
	s = apparat.NewSystem()
	if pause {
		paused = true
		s.SetSpeed(0)
	}

	if romFile != "" {
		f, err := os.Open(romFile)
		if err != nil {
			fmt.Printf("error opening ROM file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Printf("error reading ROM contents: %v\n", err)
			os.Exit(1)
		}
		s.LoadROM(b)
	}
}

type statusBar struct {
	sys *apparat.System
}

func (s statusBar) Draw() {
	bg := termbox.ColorWhite
	fg := termbox.ColorBlack
	wr := " chipper - running: "
	last := 0
	for i, c := range wr {
		termbox.SetCell(i, 0, c, fg, bg)
		last = i
	}
	if s.sys.IsRunning() {
		termbox.SetCell(last+1, 0, 'y', termbox.ColorGreen, bg)
	} else {
		termbox.SetCell(last+1, 0, 'n', termbox.ColorRed, bg)
	}
	wr = " paused: "
	last = last + 2
	for _, c := range wr {
		termbox.SetCell(last, 0, c, fg, bg)
		last++
	}
	if paused {
		termbox.SetCell(last, 0, 'y', termbox.ColorRed, bg)
	} else {
		termbox.SetCell(last, 0, 'n', termbox.ColorGreen, bg)
	}
	termbox.SetCell(last+1, 0, ' ', fg, bg)
	w, _ := termbox.Size()
	for i := last + 2; i <= w; i++ {
		termbox.SetCell(i, 0, ' ', fg, bg)
	}
}
