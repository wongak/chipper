// Copyright 2016 The Ebiten Authors
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

package twenty48

import (
	"github.com/hajimehoshi/ebiten"
)

// Dir represents a direction.
type Dir int

const (
	DirUp Dir = iota
	DirRight
	DirDown
	DirLeft
)

// String returns a string representing the direction.
func (d Dir) String() string {
	switch d {
	case DirUp:
		return "Up"
	case DirRight:
		return "Right"
	case DirDown:
		return "Down"
	case DirLeft:
		return "Left"
	}
	panic("not reach")
}

// Vector returns a [-1, 1] value for each axis.
func (d Dir) Vector() (x, y int) {
	switch d {
	case DirUp:
		return 0, -1
	case DirRight:
		return 1, 0
	case DirDown:
		return 0, 1
	case DirLeft:
		return -1, 0
	}
	panic("not reach")
}

// Input represents the current key states.
type Input struct {
	keyState map[ebiten.Key]int
}

// NewInput generates a new Input object.
func NewInput() *Input {
	return &Input{
		keyState: map[ebiten.Key]int{},
	}
}

var (
	dirKeys = map[ebiten.Key]Dir{
		ebiten.KeyUp:    DirUp,
		ebiten.KeyRight: DirRight,
		ebiten.KeyDown:  DirDown,
		ebiten.KeyLeft:  DirLeft,
	}
)

// Update updates the current input states.
func (i *Input) Update() {
	for k := range dirKeys {
		if ebiten.IsKeyPressed(k) {
			i.keyState[k]++
		} else {
			i.keyState[k] = 0
		}
	}
}

// Dir returns a currenly pressed direction.
// Dir returns false if no direction key is pressed.
func (i *Input) Dir() (Dir, bool) {
	for k, d := range dirKeys {
		if i.keyState[k] == 1 {
			return d, true
		}
	}
	return 0, false
}
