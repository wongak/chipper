package main

import (
	"fmt"
	"os"

	termbox "github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		fmt.Printf("error initializing term: %v", err)
		os.Exit(1)
	}
	defer termbox.Close()
}
