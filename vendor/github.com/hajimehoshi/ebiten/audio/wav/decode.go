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

// Package wav provides WAV (RIFF) decoder.
package wav

import (
	"bytes"
	"fmt"
	"io"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/internal/resampling"
)

type readSeekCloseSizer interface {
	audio.ReadSeekCloser
	Size() int64
}

// Stream is a decoded audio stream.
type Stream struct {
	inner readSeekCloseSizer
}

// Read is implementation of io.Reader's Read.
func (s *Stream) Read(p []byte) (int, error) {
	return s.inner.Read(p)
}

// Seek is implementation of io.Seeker's Seek.
//
// Note that Seek can take long since decoding is a relatively heavy task.
func (s *Stream) Seek(offset int64, whence int) (int64, error) {
	return s.inner.Seek(offset, whence)
}

// Read is implementation of io.Closer's Close.
func (s *Stream) Close() error {
	return s.inner.Close()
}

// Size returns the size of decoded stream in bytes.
func (s *Stream) Size() int64 {
	return s.inner.Size()
}

type stream struct {
	src        audio.ReadSeekCloser
	headerSize int64
	dataSize   int64
	remaining  int64
}

// Read is implementation of io.Reader's Read.
func (s *stream) Read(p []byte) (int, error) {
	if s.remaining <= 0 {
		return 0, io.EOF
	}
	if s.remaining < int64(len(p)) {
		p = p[0:s.remaining]
	}
	n, err := s.src.Read(p)
	s.remaining -= int64(n)
	return n, err
}

// Seek is implementation of io.Seeker's Seek.
func (s *stream) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart {
		offset += s.headerSize
	}
	n, err := s.src.Seek(offset, whence)
	if err != nil {
		return 0, err
	}
	if n-s.headerSize < 0 {
		return 0, fmt.Errorf("wav: invalid offset")
	}
	s.remaining = s.dataSize - (n - s.headerSize)
	// There could be a tail in wav file.
	if s.remaining < 0 {
		s.remaining = 0
		return s.dataSize, nil
	}
	return n - s.headerSize, nil
}

// Read is implementation of io.Closer's Close.
func (s *stream) Close() error {
	return s.src.Close()
}

// Size returns the size of decoded stream in bytes.
func (s *stream) Size() int64 {
	return s.dataSize
}

// Decode decodes WAV (RIFF) data to playable stream.
//
// The format must be 2 channels, 16bit little endian PCM.
//
// Sample rate is automatically adjusted to fit with the audio context.
func Decode(context *audio.Context, src audio.ReadSeekCloser) (*Stream, error) {
	buf := make([]byte, 12)
	n, err := io.ReadFull(src, buf)
	if n != len(buf) {
		return nil, fmt.Errorf("wav: invalid header")
	}
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(buf[0:4], []byte("RIFF")) {
		return nil, fmt.Errorf("wav: invalid header: 'RIFF' not found")
	}
	if !bytes.Equal(buf[8:12], []byte("WAVE")) {
		return nil, fmt.Errorf("wav: invalid header: 'WAVE' not found")
	}

	// Read chunks
	dataSize := int64(0)
	headerSize := int64(0)
	sampleRateFrom := 0
	sampleRateTo := 0
chunks:
	for {
		buf := make([]byte, 8)
		n, err := io.ReadFull(src, buf)
		if n != len(buf) {
			return nil, fmt.Errorf("wav: invalid header")
		}
		if err != nil {
			return nil, err
		}
		headerSize += 8
		size := int64(buf[4]) | int64(buf[5])<<8 | int64(buf[6])<<16 | int64(buf[7])<<24
		switch {
		case bytes.Equal(buf[0:4], []byte("fmt ")):
			if size != 16 {
				return nil, fmt.Errorf("wav: invalid header: maybe non-PCM file?")
			}
			buf := make([]byte, size)
			n, err := io.ReadFull(src, buf)
			if n != len(buf) {
				return nil, fmt.Errorf("wav: invalid header")
			}
			if err != nil {
				return nil, err
			}
			format := int(buf[0]) | int(buf[1])<<8
			if format != 1 {
				return nil, fmt.Errorf("wav: format must be linear PCM")
			}
			channelNum := int(buf[2]) | int(buf[3])<<8
			// TODO: Remove this magic number
			if channelNum != 2 {
				return nil, fmt.Errorf("wav: channel num must be 2")
			}
			bitsPerSample := int(buf[14]) | int(buf[15])<<8
			// TODO: Remove this magic number
			if bitsPerSample != 16 {
				return nil, fmt.Errorf("wav: bits per sample must be 16")
			}
			sampleRate := int64(buf[4]) | int64(buf[5])<<8 | int64(buf[6])<<16 | int64(buf[7])<<24
			if int64(context.SampleRate()) != sampleRate {
				sampleRateFrom = int(sampleRate)
				sampleRateTo = context.SampleRate()
			}
			headerSize += size
		case bytes.Equal(buf[0:4], []byte("data")):
			dataSize = size
			break chunks
		default:
			buf := make([]byte, size)
			n, err := io.ReadFull(src, buf)
			if n != len(buf) {
				return nil, fmt.Errorf("wav: invalid header")
			}
			if err != nil {
				return nil, err
			}
			headerSize += size
		}
	}
	s := &stream{
		src:        src,
		headerSize: headerSize,
		dataSize:   dataSize,
		remaining:  dataSize,
	}
	if sampleRateFrom != sampleRateTo {
		fixed := resampling.NewStream(s, s.dataSize, sampleRateFrom, sampleRateTo)
		return &Stream{fixed}, nil
	}
	return &Stream{s}, nil
}