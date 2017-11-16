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

// +build example

package blocks

import (
	"github.com/hajimehoshi/ebiten"
)

const maxFlushCount = 20

type Field struct {
	blocks        [fieldBlockNumX][fieldBlockNumY]BlockType
	flushCount    int
	onEndFlushing func(int)
}

func NewField() *Field {
	return &Field{}
}

func (f *Field) IsBlocked(x, y int) bool {
	if x < 0 || fieldBlockNumX <= x {
		return true
	}
	if y < 0 {
		return false
	}
	if fieldBlockNumY <= y {
		return true
	}
	return f.blocks[x][y] != BlockTypeNone
}

func (f *Field) collides(piece *Piece, x, y int, angle Angle) bool {
	return piece.collides(f, x, y, angle)
}

func (f *Field) MovePieceToLeft(piece *Piece, x, y int, angle Angle) int {
	if f.collides(piece, x-1, y, angle) {
		return x
	}
	return x - 1
}

func (f *Field) MovePieceToRight(piece *Piece, x, y int, angle Angle) int {
	if f.collides(piece, x+1, y, angle) {
		return x
	}
	return x + 1
}

func (f *Field) PieceDroppable(piece *Piece, x, y int, angle Angle) bool {
	return !f.collides(piece, x, y+1, angle)
}

func (f *Field) DropPiece(piece *Piece, x, y int, angle Angle) int {
	if f.collides(piece, x, y+1, angle) {
		return y
	}
	return y + 1
}

func (f *Field) RotatePieceRight(piece *Piece, x, y int, angle Angle) Angle {
	if f.collides(piece, x, y, angle.RotateRight()) {
		return angle
	}
	return angle.RotateRight()
}

func (f *Field) RotatePieceLeft(piece *Piece, x, y int, angle Angle) Angle {
	if f.collides(piece, x, y, angle.RotateLeft()) {
		return angle
	}
	return angle.RotateLeft()
}

func (f *Field) AbsorbPiece(piece *Piece, x, y int, angle Angle) {
	piece.AbsorbInto(f, x, y, angle)
	if f.flushable() {
		f.flushCount = maxFlushCount
	}
}

func (f *Field) Flushing() bool {
	return 0 < f.flushCount
}

func (f *Field) SetEndFlushing(fn func(lines int)) {
	f.onEndFlushing = fn
}

func (f *Field) flushable() bool {
	for j := fieldBlockNumY - 1; 0 <= j; j-- {
		if f.flushableLine(j) {
			return true
		}
	}
	return false
}

func (f *Field) flushableLine(j int) bool {
	for i := 0; i < fieldBlockNumX; i++ {
		if f.blocks[i][j] == BlockTypeNone {
			return false
		}
	}
	return true
}

func (f *Field) setBlock(x, y int, blockType BlockType) {
	f.blocks[x][y] = blockType
}

func (f *Field) endFlushing() int {
	flushedLineNum := 0
	for j := fieldBlockNumY - 1; 0 <= j; j-- {
		if f.flushLine(j + flushedLineNum) {
			flushedLineNum++
		}
	}
	return flushedLineNum
}

func (f *Field) flushLine(j int) bool {
	for i := 0; i < fieldBlockNumX; i++ {
		if f.blocks[i][j] == BlockTypeNone {
			return false
		}
	}
	for j2 := j; 1 <= j2; j2-- {
		for i := 0; i < fieldBlockNumX; i++ {
			f.blocks[i][j2] = f.blocks[i][j2-1]
		}
	}
	for i := 0; i < fieldBlockNumX; i++ {
		f.blocks[i][0] = BlockTypeNone
	}
	return true
}

func (f *Field) Update() {
	if 0 <= f.flushCount {
		f.flushCount--
		if f.flushCount == 0 {
			l := f.endFlushing()
			if f.onEndFlushing != nil {
				f.onEndFlushing(l)
			}
		}
	}
}

func min(a, b float64) float64 {
	if a > b {
		return b
	}
	return a
}

func (f *Field) flushingColor() ebiten.ColorM {
	c := f.flushCount
	if c < 0 {
		c = 0
	}
	rate := float64(c) / maxFlushCount
	clr := ebiten.ColorM{}
	alpha := min(1, rate*2)
	clr.Scale(1, 1, 1, alpha)
	r := min(1, (1-rate)*2)
	clr.Translate(r, 0, 0, 0)
	return clr
}

func (f *Field) Draw(r *ebiten.Image, x, y int) {
	blocks := make([][]BlockType, len(f.blocks))
	flushingBlocks := make([][]BlockType, len(f.blocks))
	for i := 0; i < fieldBlockNumX; i++ {
		blocks[i] = make([]BlockType, fieldBlockNumY)
		flushingBlocks[i] = make([]BlockType, fieldBlockNumY)
	}
	for j := 0; j < fieldBlockNumY; j++ {
		if f.flushableLine(j) {
			for i := 0; i < fieldBlockNumX; i++ {
				flushingBlocks[i][j] = f.blocks[i][j]
			}
		} else {
			for i := 0; i < fieldBlockNumX; i++ {
				blocks[i][j] = f.blocks[i][j]
			}
		}
	}
	drawBlocks(r, blocks, x, y, ebiten.ColorM{})
	drawBlocks(r, flushingBlocks, x, y, f.flushingColor())
}
