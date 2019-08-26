// Copyright © 2019 mg
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

package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

const (
	imgWidth  = 512
	imgHeight = 512
)

var (
	seed = 123456789
)

func simpleRnd() int {
	a := 1103515245
	c := 12345
	m := 1 << 31

	seed = (a*seed + c) % m
	return seed
}

func main() {
	m := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			r := simpleRnd()
			p := uint8(r % 2 * 255)
			m.Set(x, y, color.RGBA{R: p, G: p, B: p, A: 255})
		}
	}

	o, _ := os.Create("rand_crypto.png")
	png.Encode(o, m)
	o.Close()
}
