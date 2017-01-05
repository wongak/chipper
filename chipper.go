package main

import (
	"time"

	"github.com/wongak/chipper/pkg/apparat"
)

func main() {
	tm := apparat.NewTimers()
	tm.SetSound(60)
	time.Sleep(20 * time.Second)
}
