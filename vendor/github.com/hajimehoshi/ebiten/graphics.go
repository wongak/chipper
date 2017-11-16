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

package ebiten

import (
	"github.com/hajimehoshi/ebiten/internal/opengl"
)

// Filter represents the type of filter to be used when an image is maginified or minified.
type Filter int

const (
	// FilterNearest represents nearest (crisp-edged) filter
	FilterNearest Filter = iota

	// FilterLinear represents linear filter
	FilterLinear
)

func glFilter(filter Filter) opengl.Filter {
	switch filter {
	case FilterNearest:
		return opengl.Nearest
	case FilterLinear:
		return opengl.Linear
	}
	panic("not reach")
}

// CompositeMode represents Porter-Duff composition mode.
type CompositeMode int

// This name convention follows CSS compositing: https://drafts.fxtf.org/compositing-2/.
//
// In the comments,
// c_src, c_dst and c_out represent alpha-premultiplied RGB values of source, destination and output respectively. α_src and α_dst represent alpha values of source and destination respectively.
const (
	// Regular alpha blending
	// c_out = c_src + c_dst × (1 - α_src)
	CompositeModeSourceOver CompositeMode = CompositeMode(opengl.CompositeModeSourceOver)

	// c_out = 0
	CompositeModeClear = CompositeMode(opengl.CompositeModeClear)

	// c_out = c_src
	CompositeModeCopy = CompositeMode(opengl.CompositeModeCopy)

	// c_out = c_dst
	CompositeModeDestination = CompositeMode(opengl.CompositeModeDestination)

	// c_out = c_src × (1 - α_dst) + c_dst
	CompositeModeDestinationOver = CompositeMode(opengl.CompositeModeDestinationOver)

	// c_out = c_src × α_dst
	CompositeModeSourceIn = CompositeMode(opengl.CompositeModeSourceIn)

	// c_out = c_dst × α_src
	CompositeModeDestinationIn = CompositeMode(opengl.CompositeModeDestinationIn)

	// c_out = c_src × (1 - α_dst)
	CompositeModeSourceOut = CompositeMode(opengl.CompositeModeSourceOut)

	// c_out = c_dst × (1 - α_src)
	CompositeModeDestinationOut = CompositeMode(opengl.CompositeModeDestinationOut)

	// c_out = c_src × α_dst + c_dst × (1 - α_src)
	CompositeModeSourceAtop = CompositeMode(opengl.CompositeModeSourceAtop)

	// c_out = c_src × (1 - α_dst) + c_dst × α_src
	CompositeModeDestinationAtop = CompositeMode(opengl.CompositeModeDestinationAtop)

	// c_out = c_src × (1 - α_dst) + c_dst × (1 - α_src)
	CompositeModeXor = CompositeMode(opengl.CompositeModeXor)

	// Sum of source and destination (a.k.a. 'plus' or 'additive')
	// c_out = c_src + c_dst
	CompositeModeLighter = CompositeMode(opengl.CompositeModeLighter)
)
