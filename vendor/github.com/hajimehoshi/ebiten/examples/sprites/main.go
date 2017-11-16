// Copyright 2015 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build example

package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
	maxAngle     = 256
)

var (
	ebitenImage *ebiten.Image
)

type Sprite struct {
	imageWidth  int
	imageHeight int
	x           int
	y           int
	vx          int
	vy          int
	angle       int
}

func (s *Sprite) Update() {
	s.x += s.vx
	s.y += s.vy
	if s.x < 0 {
		s.x = -s.x
		s.vx = -s.vx
	} else if screenWidth <= s.x+s.imageWidth {
		s.x = 2*(screenWidth-s.imageWidth) - s.x
		s.vx = -s.vx
	}
	if s.y < 0 {
		s.y = -s.y
		s.vy = -s.vy
	} else if screenHeight <= s.y+s.imageHeight {
		s.y = 2*(screenHeight-s.imageHeight) - s.y
		s.vy = -s.vy
	}
	s.angle++
	s.angle %= maxAngle
}

type Sprites struct {
	sprites []*Sprite
	num     int
}

func (s *Sprites) Update() {
	for _, sprite := range s.sprites {
		sprite.Update()
	}
}

const (
	MinSprites = 0
	MaxSprites = 50000
)

var (
	sprites = &Sprites{make([]*Sprite, MaxSprites), 500}
	op      = &ebiten.DrawImageOptions{}
)

func update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		sprites.num -= 20
		if sprites.num < MinSprites {
			sprites.num = MinSprites
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		sprites.num += 20
		if MaxSprites < sprites.num {
			sprites.num = MaxSprites
		}
	}
	sprites.Update()

	if ebiten.IsRunningSlowly() {
		return nil
	}
	w, h := ebitenImage.Size()
	for i := 0; i < sprites.num; i++ {
		s := sprites.sprites[i]
		op.GeoM.Reset()
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Rotate(2 * math.Pi * float64(s.angle) / maxAngle)
		op.GeoM.Translate(float64(w)/2, float64(h)/2)
		op.GeoM.Translate(float64(s.x), float64(s.y))
		screen.DrawImage(ebitenImage, op)
	}
	msg := fmt.Sprintf(`FPS: %0.2f
Num of sprites: %d
Press <- or -> to change the number of sprites`, ebiten.CurrentFPS(), sprites.num)
	ebitenutil.DebugPrint(screen, msg)
	return nil
}

func main() {
	var err error
	img, _, err := ebitenutil.NewImageFromFile("_resources/images/ebiten.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	w, h := img.Size()
	ebitenImage, _ = ebiten.NewImage(w, h, ebiten.FilterNearest)
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.5)
	ebitenImage.DrawImage(img, op)
	for i := range sprites.sprites {
		w, h := ebitenImage.Size()
		x, y := rand.Intn(screenWidth-w), rand.Intn(screenHeight-h)
		vx, vy := 2*rand.Intn(2)-1, 2*rand.Intn(2)-1
		a := rand.Intn(maxAngle)
		sprites.sprites[i] = &Sprite{
			imageWidth:  w,
			imageHeight: h,
			x:           x,
			y:           y,
			vx:          vx,
			vy:          vy,
			angle:       a,
		}
	}
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Sprites (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}
