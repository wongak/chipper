// Copyright 2016 Hajime Hoshi
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

// +build darwin freebsd linux windows
// +build !js
// +build !android
// +build !ios

package ebitenutil

import (
	"os"
)

// OpenFile opens a file and returns a stream for its data.
//
// This function is available both on desktops and browsers.
// Note that this doesn't work on mobiles.
func OpenFile(path string) (ReadSeekCloser, error) {
	return os.Open(path)
}
