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
	"fmt"
	"strings"

	"github.com/gonvenience/term"
)

// Internal bit mask to mark feature states, e.g. foreground coloring
const (
	fgMask        = 0x1
	bgMask        = 0x2
	boldMask      = 0x4
	italicMask    = 0x8
	underlineMask = 0x10
)

// ColorSetting defines the coloring setting to be used
var ColorSetting = AUTO

// TrueColorSetting defines the true color usage setting to be used
var TrueColorSetting = AUTO

// switchState is the type to cover different preferences/settings like: on, off, or auto
type switchState int

// Supported setting states
const (
	ON = switchState(iota)
	OFF
	AUTO
)

func (s switchState) String() string {
	switch s {
	case ON:
		return "on"

	case OFF:
		return "off"

	case AUTO:
		return "auto"
	}

	panic("unsupported switch state")
}

func (s *switchState) Set(setting string) error {
	switch strings.ToLower(setting) {
	case "auto":
		*s = AUTO

	case "off", "no", "false":
		*s = OFF

	case "on", "yes", "true":
		*s = ON

	default:
		return fmt.Errorf("invalid state '%s' used, supported modes are: auto, on, or off", setting)
	}

	return nil
}

func (s switchState) Type() string {
	return ""
}

// UseColors return whether colors are used or not based on the configured color
// setting or terminal capabilities
func UseColors() bool {
	return (ColorSetting == ON) ||
		(ColorSetting == AUTO && term.IsTerminal() && !term.IsDumbTerminal())
}

// UseTrueColor returns whether true color colors should be used or not based on
// the configured true color usage setting or terminal capabilities
func UseTrueColor() bool {
	return (TrueColorSetting == ON) ||
		(TrueColorSetting == AUTO && term.IsTrueColor())
}
