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
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	escapeSeqRegExp = regexp.MustCompile(`\x1b\[(\d+(;\d+)*)m`)
	boldMarker      = regexp.MustCompile(`\*([^*]+?)\*`)
	italicMarker    = regexp.MustCompile(`_([^_]+?)_`)
	underlineMarker = regexp.MustCompile(`~([^~]+?)~`)
	colorMarker     = regexp.MustCompile(`(#?\w+)\{([^}]+?)\}`)
)

// ParseOption defines parser options
type ParseOption func(*String) error

// ProcessTextAnnotations specifies whether during parsing bunt-specific text
// annotations like *bold*, or _italic_ should be processed.
func ProcessTextAnnotations() ParseOption {
	return func(s *String) error {
		return processTextAnnotations(s)
	}
}

// ParseString parses a string that can contain both ANSI escape code Select
// Graphic Rendition (SGR) codes and Markdown style text annotations, for
// example *bold* or _italic_.
// SGR details : https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
func ParseString(input string, opts ...ParseOption) (*String, error) {
	var (
		pointer int
		current uint64
		result  String
		err     error
	)

	// Special case: the escape sequence without any parameter is equivalent to
	// the reset escape sequence.
	input = strings.Replace(input, "\x1b[m", "\x1b[0m", -1)

	for _, submatch := range escapeSeqRegExp.FindAllStringSubmatchIndex(input, -1) {
		fullMatchStart, fullMatchEnd := submatch[0], submatch[1]
		settingsStart, settingsEnd := submatch[2], submatch[3]

		for _, r := range input[pointer:fullMatchStart] {
			result = append(result, ColoredRune{r, current})
		}

		current, err = parseSelectGraphicRenditionEscapeSequence(input[settingsStart:settingsEnd])
		if err != nil {
			return nil, err
		}

		pointer = fullMatchEnd
	}

	// Flush the remaining input string part into the result
	for _, r := range input[pointer:] {
		result = append(result, ColoredRune{r, current})
	}

	// Process optional parser options
	for _, opt := range opts {
		if err = opt(&result); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func parseSelectGraphicRenditionEscapeSequence(escapeSeq string) (uint64, error) {
	values := []uint8{}
	for _, x := range strings.Split(escapeSeq, ";") {
		// Note: This only works, because of the regular expression only matching
		// digits. Therefore, it should be okay to omit the error.
		value, _ := strconv.Atoi(x)
		values = append(values, uint8(value))
	}

	result := uint64(0)

	for i := 0; i < len(values); i++ {
		switch values[i] {
		case 1: // bold
			result |= boldMask

		case 3: // italic
			result |= italicMask

		case 4: // underline
			result |= underlineMask

		case 30: // Black
			result |= fgRGBMask(0, 0, 0)

		case 31: // Red
			result |= fgRGBMask(170, 0, 0)

		case 32: // Green
			result |= fgRGBMask(0, 170, 0)

		case 33: // Yellow
			result |= fgRGBMask(229, 229, 16)

		case 34: // Blue
			result |= fgRGBMask(0, 0, 170)

		case 35: // Magenta
			result |= fgRGBMask(170, 0, 170)

		case 36: // Cyan
			result |= fgRGBMask(0, 170, 170)

		case 37: // White
			result |= fgRGBMask(229, 229, 229)

		case 90: // Bright Black (Gray)
			result |= fgRGBMask(85, 85, 85)

		case 91: // Bright Red
			result |= fgRGBMask(255, 85, 85)

		case 92: // Bright Green
			result |= fgRGBMask(85, 255, 85)

		case 93: // Bright Yellow
			result |= fgRGBMask(255, 255, 85)

		case 94: // Bright Blue
			result |= fgRGBMask(85, 85, 255)

		case 95: // Bright Magenta
			result |= fgRGBMask(255, 85, 255)

		case 96: // Bright Cyan
			result |= fgRGBMask(85, 255, 255)

		case 97: // Bright White
			result |= fgRGBMask(255, 255, 255)

		case 40: // Black
			result |= bgRGBMask(0, 0, 0)

		case 41: // Red
			result |= bgRGBMask(170, 0, 0)

		case 42: // Green
			result |= bgRGBMask(0, 170, 0)

		case 43: // Yellow
			result |= bgRGBMask(229, 229, 16)

		case 44: // Blue
			result |= bgRGBMask(0, 0, 170)

		case 45: // Magenta
			result |= bgRGBMask(170, 0, 170)

		case 46: // Cyan
			result |= bgRGBMask(0, 170, 170)

		case 47: // White
			result |= bgRGBMask(229, 229, 229)

		case 100: // Bright Black (Gray)
			result |= bgRGBMask(85, 85, 85)

		case 101: // Bright Red
			result |= bgRGBMask(255, 85, 85)

		case 102: // Bright Green
			result |= bgRGBMask(85, 255, 85)

		case 103: // Bright Yellow
			result |= bgRGBMask(255, 255, 85)

		case 104: // Bright Blue
			result |= bgRGBMask(85, 85, 255)

		case 105: // Bright Magenta
			result |= bgRGBMask(255, 85, 255)

		case 106: // Bright Cyan
			result |= bgRGBMask(85, 255, 255)

		case 107: // Bright White
			result |= bgRGBMask(255, 255, 255)

		case 38: // foreground color
			switch {
			case len(values) > 4 && values[i+1] == 2:
				result |= fgRGBMask(uint64(values[i+2]), uint64(values[i+3]), uint64(values[i+4]))
				i += 4

			case len(values) > 2 && values[i+1] == 5:
				r, g, b := lookUp8bitColor(values[i+2])
				result |= fgRGBMask(uint64(r), uint64(g), uint64(b))
				i += 2

			default:
				return 0, fmt.Errorf("unsupported foreground color selection '%v'", values)
			}

		case 48: // background color
			switch {
			case len(values) > 4 && values[i+1] == 2:
				result |= bgRGBMask(uint64(values[i+2]), uint64(values[i+3]), uint64(values[i+4]))
				i += 4

			case len(values) > 2 && values[i+1] == 5:
				r, g, b := lookUp8bitColor(values[i+2])
				result |= bgRGBMask(uint64(r), uint64(g), uint64(b))
				i += 2

			default:
				return 0, fmt.Errorf("unsupported background color selection '%v'", values)

			}
		}
	}

	return result, nil
}

func processTextAnnotations(text *String) error {
	var buffer bytes.Buffer
	for _, coloredRune := range *text {
		_ = buffer.WriteByte(byte(coloredRune.Symbol))
	}

	raw := buffer.String()
	toBeDeleted := []int{}

	// Process text annotation markers for bold, italic and underline
	helperMap := map[uint64]*regexp.Regexp{
		boldMask:      boldMarker,
		italicMask:    italicMarker,
		underlineMask: underlineMarker,
	}

	for mask, regex := range helperMap {
		for _, match := range regex.FindAllStringSubmatchIndex(raw, -1) {
			fullMatchStart, fullMatchEnd := match[0], match[1]
			textStart, textEnd := match[2], match[3]

			for i := textStart; i < textEnd; i++ {
				(*text)[i].Settings |= mask
			}

			toBeDeleted = append(toBeDeleted, fullMatchStart, fullMatchEnd-1)
		}
	}

	// Process text annotation markers that specify a foreground color for a
	// specific part of the text
	for _, match := range colorMarker.FindAllStringSubmatchIndex(raw, -1) {
		fullMatchStart, fullMatchEnd := match[0], match[1]
		colorName := raw[match[2]:match[3]]
		textStart, textEnd := match[4], match[5]

		color := lookupColorByName(colorName)
		if color == nil {
			return fmt.Errorf("unable to find color by name: %s", colorName)
		}

		r, g, b := color.RGB255()
		for i := textStart; i < textEnd; i++ {
			(*text)[i].Settings |= fgMask
			(*text)[i].Settings |= uint64(r) << 8
			(*text)[i].Settings |= uint64(g) << 16
			(*text)[i].Settings |= uint64(b) << 24
		}

		for i := fullMatchStart; i < fullMatchEnd; i++ {
			if i < textStart || i > textEnd-1 {
				toBeDeleted = append(toBeDeleted, i)
			}
		}
	}

	// Finally, sort the runes to be deleted in descending order and delete them
	// one by one to get rid of the text annotation markers
	sort.Slice(toBeDeleted, func(i, j int) bool {
		return toBeDeleted[i] > toBeDeleted[j]
	})

	for _, idx := range toBeDeleted {
		(*text) = append((*text)[:idx], (*text)[idx+1:]...)
	}

	return nil
}

func fgRGBMask(r, g, b uint64) uint64 {
	return fgMask | r<<8 | g<<16 | b<<24
}

func bgRGBMask(r, g, b uint64) uint64 {
	return bgMask | r<<32 | g<<40 | b<<48
}

func lookUp8bitColor(n uint8) (r, g, b uint8) {
	return colorPalette8bit[n].RGB255()
}
