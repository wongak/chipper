// Copyright 2014 Hajime Hoshi
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

package ebitenutil

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil/internal/assets"
)

type debugPrintState struct {
	textImage              *ebiten.Image
	debugPrintRenderTarget *ebiten.Image
}

var defaultDebugPrintState = &debugPrintState{}

// DebugPrint draws the string str on the image.
//
// DebugPrint always returns nil as of 1.5.0-alpha.
func DebugPrint(image *ebiten.Image, str string) error {
	defaultDebugPrintState.DebugPrint(image, str)
	return nil
}

func (d *debugPrintState) drawText(rt *ebiten.Image, str string, ox, oy int, c color.Color) {
	op := &ebiten.DrawImageOptions{}
	ur, ug, ub, ua := c.RGBA()
	const max = math.MaxUint16
	r := float64(ur) / max
	g := float64(ug) / max
	b := float64(ub) / max
	a := float64(ua) / max
	if 0 < a {
		r /= a
		g /= a
		b /= a
	}
	op.ColorM.Scale(r, g, b, a)
	x := 0
	y := 0
	for _, c := range str {
		const cw = assets.TextImageCharWidth
		const ch = assets.TextImageCharHeight
		if c == '\n' {
			x = 0
			y += ch
			continue
		}
		const n = assets.TextImageWidth / cw
		sx := (int(c) % n) * cw
		sy := (int(c) / n) * ch
		r := image.Rect(sx, sy, sx+cw, sy+ch)
		op.SourceRect = &r
		op.GeoM.Reset()
		op.GeoM.Translate(float64(x), float64(y))
		op.GeoM.Translate(float64(ox+1), float64(oy))
		_ = rt.DrawImage(d.textImage, op)
		x += cw
	}
}

// DebugPrint prints the given text str on the given image r.
func (d *debugPrintState) DebugPrint(r *ebiten.Image, str string) {
	if d.textImage == nil {
		img := assets.TextImage()
		d.textImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)
	}
	if d.debugPrintRenderTarget == nil {
		width, height := 256, 256
		d.debugPrintRenderTarget, _ = ebiten.NewImage(width, height, ebiten.FilterNearest)
	}
	d.drawText(r, str, 1, 1, color.NRGBA{0x00, 0x00, 0x00, 0x80})
	d.drawText(r, str, 0, 0, color.NRGBA{0xff, 0xff, 0xff, 0xff})
}
