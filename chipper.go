package main

import (
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/wongak/chipper/pkg/apparat"
)

var (
	romFile string
	pause   bool
	debug   bool
)

var (
	// instructions per second
	ips int
)

var (
	s       *apparat.System
	paused  bool
	stopped bool

	initialized bool
)

var (
	squareBlack, squareWhite *ebiten.Image
	errStopped               error
)

func init() {
	errStopped = fmt.Errorf("stopped")
}

func main() {
	flag.StringVar(&romFile, "rom", "", "ROM file to load")
	flag.BoolVar(&pause, "p", true, "whether to pause on start")
	flag.BoolVar(&debug, "d", true, "debug mode")
	flag.Parse()

	ips = 10
	initSystem()

	err := ebiten.Run(update, 320, 240, 2, "chipper")
	if err != nil && err != errStopped {
		fmt.Fprintf(os.Stderr, "system err: %v", err)
		os.Exit(1)
	}

	//mainLoop:
	//	for {
	//		status.Draw()
	//
	//		ev := termbox.PollEvent()
	//		switch ev.Type {
	//		case termbox.EventError:
	//			fmt.Printf("error event: %v", ev.Err)
	//			os.Exit(1)
	//
	//		case termbox.EventInterrupt:
	//			for y := 0; y < apparat.DisplayHeight; y++ {
	//				l := s.Dsp.Line(uint8(y))
	//				for x := uint16(0); x < apparat.DisplayWidth; x++ {
	//					mask := uint64(1) << (63 - x)
	//					if mask&l != 0 {
	//						termbox.SetCell(int(x), y+1, ' ', termbox.ColorBlack, termbox.ColorWhite)
	//					} else {
	//						termbox.SetCell(int(x), y+1, ' ', termbox.ColorWhite, termbox.ColorBlack)
	//					}
	//				}
	//			}
	//			if debug {
	//				drawDebug()
	//			}
	//			termbox.Flush()
	//
	//		case termbox.EventKey:
	//			switch ev.Key {
	//			case termbox.KeyCtrlC:
	//				break mainLoop
	//
	//			case termbox.KeySpace:
	//				if paused {
	//					s.Step()
	//				}
	//
	//			default:
	//				if ev.Ch != 0 {
	//					switch ev.Ch {
	//					case 'p':
	//						if paused {
	//							s.SetSpeed(10)
	//						} else {
	//							s.SetSpeed(0)
	//						}
	//						paused = !paused
	//					case 'l':
	//						showMem = !showMem
	//
	//					case '1':
	//						s.Key.SetState(1)
	//					case '2':
	//						s.Key.SetState(2)
	//					case '3':
	//						s.Key.SetState(3)
	//					case 'q':
	//						s.Key.SetState(4)
	//					case 'w':
	//						s.Key.SetState(5)
	//					case 'e':
	//						s.Key.SetState(6)
	//					case 'a':
	//						s.Key.SetState(7)
	//					case 's':
	//						s.Key.SetState(8)
	//					case 'd':
	//						s.Key.SetState(9)
	//					case 'z', 'y':
	//						s.Key.SetState(0xA)
	//					case 'x':
	//						s.Key.SetState(0)
	//					case 'c':
	//						s.Key.SetState(0xB)
	//					case '4':
	//						s.Key.SetState(0xC)
	//					case 'r':
	//						s.Key.SetState(0xD)
	//					case 'f':
	//						s.Key.SetState(0xE)
	//					case 'v':
	//						s.Key.SetState(0xF)
	//					}
	//				}
	//			}
	//		}
	//	}
}

func initSystem() {
	s = apparat.NewSystem()
	if pause {
		paused = true
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

func update(screen *ebiten.Image) error {
	if stopped {
		return errStopped
	}
	var err error
	if !initialized {
		screen.Fill(color.NRGBA{0, 0, 0, 0xff})

		squareBlack, _ = ebiten.NewImage(5, 5, ebiten.FilterNearest)
		squareBlack.Fill(color.Black)

		squareWhite, _ = ebiten.NewImage(5, 5, ebiten.FilterNearest)
		squareWhite.Fill(color.White)

		initialized = true
	}

	err = handleInput()
	if err != nil {
		return fmt.Errorf("error handling input: %v", err)
	}

	for i := 0; i < ips; i++ {
		s.Step()
	}
	s.Tick()

	for y := 0; y < apparat.DisplayHeight; y++ {
		l := s.Dsp.Line(uint8(y))
		for x := uint16(0); x < apparat.DisplayWidth; x++ {
			mask := uint64(1) << (63 - x)
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(x)*5, float64(y)*5)
			if mask&l != 0 {
				screen.DrawImage(squareWhite, opts)
			} else {
				screen.DrawImage(squareBlack, opts)
			}
		}
	}

	return nil
}

type k uint8

func (k k) KeyPressed(comp uint8) bool {
	return uint8(k) == comp
}

func (k k) HasState() bool {
	return true
}

func (k k) State() uint8 {
	return uint8(k)
}

type noK struct{}

func (_ noK) KeyPressed(_ uint8) bool {
	return false
}

func (_ noK) HasState() bool {
	return false
}

func (_ noK) State() uint8 {
	return 0
}

func handleInput() error {
	switch true {
	case ebiten.IsKeyPressed(ebiten.Key1):
		s.Key = k(1)
	case ebiten.IsKeyPressed(ebiten.Key2):
		s.Key = k(2)
	case ebiten.IsKeyPressed(ebiten.Key3):
		s.Key = k(3)
	case ebiten.IsKeyPressed(ebiten.Key4):
		s.Key = k(0x12)
	case ebiten.IsKeyPressed(ebiten.KeyQ):
		s.Key = k(4)
	case ebiten.IsKeyPressed(ebiten.KeyW):
		s.Key = k(5)
	case ebiten.IsKeyPressed(ebiten.KeyE):
		s.Key = k(6)
	case ebiten.IsKeyPressed(ebiten.KeyR):
		s.Key = k(0x13)
	case ebiten.IsKeyPressed(ebiten.KeyA):
		s.Key = k(7)
	case ebiten.IsKeyPressed(ebiten.KeyS):
		s.Key = k(8)
	case ebiten.IsKeyPressed(ebiten.KeyD):
		s.Key = k(9)
	case ebiten.IsKeyPressed(ebiten.KeyF):
		s.Key = k(0x14)
	case ebiten.IsKeyPressed(ebiten.KeyZ), ebiten.IsKeyPressed(ebiten.KeyY):
		s.Key = k(0x10)
	case ebiten.IsKeyPressed(ebiten.KeyX):
		s.Key = k(0)
	case ebiten.IsKeyPressed(ebiten.KeyC):
		s.Key = k(0x11)
	case ebiten.IsKeyPressed(ebiten.KeyV):
		s.Key = k(0x15)

	case ebiten.IsKeyPressed(ebiten.KeyEscape):
		stopped = true
	default:
		s.Key = noK{}
	}
	return nil
}

/*
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
*/
