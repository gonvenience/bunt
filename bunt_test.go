// Copyright © 2020 The Homeport Team
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

package bunt_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gonvenience/bunt"
)

var _ = Describe("settings tests", func() {
	var parse = func(setting string) (string, error) {
		var tmp = &SwitchState{}
		err := tmp.Set(setting)
		return tmp.String(), err
	}

	Context("parse color settings", func() {
		It("should parse auto as setting auto", func() {
			setting, err := parse("auto")
			Expect(err).ToNot(HaveOccurred())
			Expect(setting).To(Equal(AUTO.String()))
		})

		It("should parse off as setting off", func() {
			setting, err := parse("off")
			Expect(err).ToNot(HaveOccurred())
			Expect(setting).To(Equal(OFF.String()))
		})

		It("should parse on as setting on", func() {
			setting, err := parse("on")
			Expect(err).ToNot(HaveOccurred())
			Expect(setting).To(Equal(ON.String()))
		})

		It("should fail to parse unknown setting", func() {
			_, err := parse("foo")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(BeEquivalentTo("invalid state 'foo' used, supported modes are: auto, on, or off"))
		})
	})
})
