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

import colorful "github.com/lucasb-eyer/go-colorful"

// StyleOption defines style option for strings
type StyleOption func(*String)

// PlainTextLength returns the length of the input text without any escape
// sequences. The function will panic in the unlikely case of a parse issue.
func PlainTextLength(text string) int {
	result, err := ParseString(text)
	if err != nil {
		panic(err)
	}

	return len(*result)
}

// Substring returns a substring of a text that may contains escape sequences.
// The function will panic in the unlikely case of a parse issue.
func Substring(text string, start int, end int) string {
	result, err := ParseString(text)
	if err != nil {
		panic(err)
	}

	result.Substring(start, end)

	return result.String()
}

// Bold applies the bold text parameter
func Bold() StyleOption {
	return func(s *String) {
		for i := range *s {
			(*s)[i].Settings |= 1 << 2
		}
	}
}

// Italic applies the italic text parameter
func Italic() StyleOption {
	return func(s *String) {
		for i := range *s {
			(*s)[i].Settings |= 1 << 3
		}
	}
}

// Foreground sets the given color as the foreground color of the text
func Foreground(color colorful.Color) StyleOption {
	r, g, b := color.RGB255()
	return func(s *String) {
		for i := range *s {
			(*s)[i].Settings |= 1
			(*s)[i].Settings |= uint64(r) << 8
			(*s)[i].Settings |= uint64(g) << 16
			(*s)[i].Settings |= uint64(b) << 24
		}
	}
}

// Style is a multi-purpose function to programmatically apply styles and other
// changes to an input text. The function will panic in the unlikely case of a
// parse issue.
func Style(text string, styleOptions ...StyleOption) string {
	result, err := ParseString(text)
	if err != nil {
		panic(err)
	}

	for _, styleOption := range styleOptions {
		styleOption(result)
	}

	return result.String()
}
