// Copyright © 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package bunt

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"

	ciede2000 "github.com/mattn/go-ciede2000"
)

func (s String) String() string {
	var (
		buffer    = &bytes.Buffer{}
		current   = uint64(0)
		useColors = UseColors()
	)

	for _, coloredRune := range s {
		if useColors && (coloredRune.Settings != current) {
			buffer.WriteString(renderSelectGraphicRenditionEscapeSequence(coloredRune.Settings))
			current = coloredRune.Settings
		}

		buffer.WriteRune(coloredRune.Symbol)
	}

	// Make sure to finish with a reset escape sequence
	if current != 0 {
		buffer.WriteString(renderSelectGraphicRenditionEscapeSequence(0))
	}

	return buffer.String()
}

func renderSelectGraphicRenditionEscapeSequence(setting uint64) string {
	if setting == 0 {
		return renderEscapeSequence(0)
	}

	parameters := []uint8{}

	if (setting & 0x04) != 0 {
		parameters = append(parameters, 1)
	}

	if (setting & 0x08) != 0 {
		parameters = append(parameters, 3)
	}

	if (setting & 0x10) != 0 {
		parameters = append(parameters, 4)
	}

	if (setting & 0x1) != 0 {
		r, g, b := uint8((setting>>8)&0xFF), uint8((setting>>16)&0xFF), uint8((setting>>24)&0xFF)
		if UseTrueColor() {
			parameters = append(parameters, 38, 2, r, g, b)

		} else {
			parameters = append(parameters, closest4bitColorParameter(r, g, b))
		}
	}

	if (setting & 0x2) != 0 {
		r, g, b := uint8((setting>>32)&0xFF), uint8((setting>>40)&0xFF), uint8((setting>>48)&0xFF)
		if UseTrueColor() {
			parameters = append(parameters, 48, 2, r, g, b)

		} else {
			parameters = append(parameters, 10+closest4bitColorParameter(r, g, b))
		}
	}

	return renderEscapeSequence(parameters...)
}

func renderEscapeSequence(a ...uint8) string {
	values := make([]string, len(a))
	for i := range a {
		values[i] = strconv.Itoa(int(a[i]))
	}

	return fmt.Sprintf("\x1b[%sm", strings.Join(values, ";"))
}

// closest4bitColorParameter returns the color attribute which matches the best
// with the provided RGB color
func closest4bitColorParameter(r, g, b uint8) uint8 {
	var (
		result    = uint8(0)
		target    = &color.RGBA{r, g, b, 0xFF}
		min       = math.MaxFloat64
		helperMap = map[uint8]*color.RGBA{
			30: &color.RGBA{0x00, 0x00, 0x00, 0xFF},
			31: &color.RGBA{0xAA, 0x00, 0x00, 0xFF},
			32: &color.RGBA{0x00, 0xAA, 0x00, 0xFF},
			33: &color.RGBA{0xFF, 0xFF, 0x00, 0xFF},
			34: &color.RGBA{0x00, 0x00, 0xAA, 0xFF},
			35: &color.RGBA{0xAA, 0x00, 0xAA, 0xFF},
			36: &color.RGBA{0x00, 0xAA, 0xAA, 0xFF},
			37: &color.RGBA{0xAA, 0xAA, 0xAA, 0xFF},
			90: &color.RGBA{0x55, 0x55, 0x55, 0xFF},
			91: &color.RGBA{0xFF, 0x55, 0x55, 0xFF},
			92: &color.RGBA{0x55, 0xFF, 0x55, 0xFF},
			93: &color.RGBA{0xFF, 0xFF, 0x55, 0xFF},
			94: &color.RGBA{0x55, 0x55, 0xFF, 0xFF},
			95: &color.RGBA{0xFF, 0x55, 0xFF, 0xFF},
			96: &color.RGBA{0x55, 0xFF, 0xFF, 0xFF},
			97: &color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
		}
	)

	// Calculate the distance between the target color and the available 4-bit
	// colors using the `deltaE` algorithm to find the best match.
	for attribute, candidate := range helperMap {
		if distance := ciede2000.Diff(target, candidate); distance < min {
			min, result = distance, attribute
		}
	}

	return result
}
