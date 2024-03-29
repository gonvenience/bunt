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

package bunt_test

import (
	"bufio"
	"bytes"
	"io"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/gonvenience/bunt"
)

var _ = Describe("print functions", func() {
	Context("process markdown style in Print functions", func() {
		var captureStdout = func(f func()) string {
			r, w, err := os.Pipe()
			Expect(err).ToNot(HaveOccurred())

			tmp := os.Stdout
			defer func() {
				os.Stdout = tmp
			}()

			os.Stdout = w
			f()
			w.Close()

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			Expect(err).ToNot(HaveOccurred())

			return buf.String()
		}

		BeforeEach(func() {
			SetColorSettings(ON, AUTO)
		})

		AfterEach(func() {
			SetColorSettings(AUTO, AUTO)
		})

		It("should parse and process markdown style in Print", func() {
			Expect(captureStdout(func() {
				_, _ = Print("This should be *bold*.")
			})).To(BeEquivalentTo("This should be \x1b[1mbold\x1b[0m."))
		})

		It("should parse and process markdown style in Printf", func() {
			Expect(captureStdout(func() {
				_, _ = Printf("This should be *%s*.", "bold")
			})).To(BeEquivalentTo("This should be \x1b[1mbold\x1b[0m."))
		})

		It("should parse and process markdown style in Println", func() {
			Expect(captureStdout(func() {
				_, _ = Println("This should be *bold*.")
			})).To(BeEquivalentTo("This should be \x1b[1mbold\x1b[0m.\n"))
		})
	})

	Context("process markdown style in Fprint functions", func() {
		var (
			buf bytes.Buffer
			out *bufio.Writer
		)

		BeforeEach(func() {
			SetColorSettings(ON, AUTO)
			buf = bytes.Buffer{}
			out = bufio.NewWriter(&buf)
		})

		AfterEach(func() {
			buf.Reset()
			out = nil
			SetColorSettings(AUTO, AUTO)
		})

		It("should parse and process markdown style in Fprint", func() {
			_, _ = Fprint(out, "This should be *bold*.")
			out.Flush()
			Expect(buf.String()).To(BeEquivalentTo("This should be \x1b[1mbold\x1b[0m."))
		})

		It("should parse and process markdown style in Fprintf", func() {
			_, _ = Fprintf(out, "This should be *%s*.", "bold")
			out.Flush()
			Expect(buf.String()).To(BeEquivalentTo("This should be \x1b[1mbold\x1b[0m."))
		})

		It("should parse and process markdown style in Fprintln", func() {
			_, _ = Fprintln(out, "This should be *bold*.")
			out.Flush()
			Expect(buf.String()).To(BeEquivalentTo("This should be \x1b[1mbold\x1b[0m.\n"))
		})
	})

	Context("process markdown style in Sprint functions", func() {
		BeforeEach(func() {
			SetColorSettings(ON, AUTO)
		})

		AfterEach(func() {
			SetColorSettings(AUTO, AUTO)
		})

		It("should parse and process markdown style in Sprint", func() {
			Expect(Sprint("This should be *bold*.")).To(BeEquivalentTo("This should be \x1b[1mbold\x1b[0m."))
		})

		It("should parse and process markdown style in Sprintf", func() {
			Expect(Sprintf("This should be *%s*.", "bold")).To(BeEquivalentTo("This should be \x1b[1mbold\x1b[0m."))
		})

		It("should parse and process markdown style in Sprintln", func() {
			Expect(Sprintln("This should be *bold*.")).To(BeEquivalentTo("This should be \x1b[1mbold\x1b[0m.\n"))
		})
	})

	Context("weird use cases and issues", func() {
		BeforeEach(func() {
			SetColorSettings(ON, AUTO)
		})

		AfterEach(func() {
			SetColorSettings(AUTO, AUTO)
		})

		It("should ignore escape sequences that cannot be processed", func() {
			Expect(Sprint("ok", "\x1b[38;2;1;2mnot ok\x1b[0m")).To(
				BeEquivalentTo("ok\x1b[38;2;1;2mnot ok\x1b[0m"))
		})
		It("should not fail writing simple types", func() {
			Expect(Sprint(42)).To(Equal("42"))
		})
		It("should not fail writing slices", func() {
			Expect(Sprint([]int{42, 1})).To(Equal("[42 1]"))
		})
	})
})
