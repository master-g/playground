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
	"fmt"
	"strings"
)

var patterns = [][]int{
	{2, 2, 2, 2, 2}, {2, 1, 1, 1, 2}, {2, 1, 2, 3, 2}, {2, 2, 1, 2, 2}, {3, 1, 1, 1, 3}, {3, 1, 2, 3, 1}, {1, 2, 3, 3, 2}, {1, 2, 1, 2, 3}, {3, 1, 1, 1, 1}, {2, 3, 2, 3, 2},
	{1, 1, 1, 1, 1}, {2, 3, 3, 3, 2}, {1, 2, 2, 2, 1}, {2, 2, 3, 2, 2}, {2, 3, 1, 3, 2}, {1, 3, 2, 1, 3}, {1, 1, 3, 3, 3}, {3, 2, 3, 2, 1}, {1, 3, 3, 3, 3}, {1, 2, 3, 3, 3},
	{3, 3, 3, 3, 3}, {1, 1, 2, 3, 3}, {3, 2, 2, 2, 3}, {1, 1, 3, 1, 1}, {2, 1, 3, 1, 2}, {1, 3, 2, 3, 1}, {3, 3, 1, 1, 1}, {2, 3, 3, 1, 1}, {3, 3, 3, 3, 1}, {3, 2, 1, 1, 1},
	{1, 2, 3, 2, 1}, {3, 3, 2, 1, 1}, {1, 2, 1, 2, 1}, {3, 3, 1, 3, 3}, {1, 3, 1, 3, 1}, {3, 1, 2, 1, 3}, {2, 1, 3, 2, 3}, {1, 1, 2, 2, 3}, {1, 1, 1, 1, 3}, {1, 2, 2, 2, 2},
	{3, 2, 1, 2, 3}, {2, 3, 2, 1, 2}, {3, 2, 3, 2, 3}, {1, 3, 3, 3, 1}, {3, 2, 3, 2, 3}, {3, 2, 1, 1, 2}, {2, 3, 1, 2, 1}, {3, 3, 2, 2, 1}, {2, 1, 2, 1, 2}, {3, 2, 2, 2, 2},
}

var sequence = []int{
	0, 5, 10, 15, 20, 25, 30, 35, 40, 45,
	1, 6, 11, 16, 21, 26, 31, 36, 41, 46,
	2, 7, 12, 17, 22, 27, 32, 37, 42, 47,
	3, 8, 13, 18, 23, 28, 33, 38, 43, 48,
	4, 9, 14, 19, 24, 29, 34, 39, 44, 49,
}

var magicSequence = []int{
	0, 10, 20, 30, 40,
	1, 11, 21, 31, 41,
	2, 12, 22, 32, 42,
	3, 13, 23, 33, 43,
	4, 14, 24, 34, 44,
	5, 15, 25, 35, 45,
	6, 16, 26, 36, 46,
	7, 17, 27, 37, 47,
	8, 18, 28, 38, 48,
	9, 19, 29, 39, 49,
}

type Pattern struct {
	index int
	data  []int
}

func NewPattern(index int, data []int) *Pattern {
	p := &Pattern{
		index: index,
		data:  make([]int, len(data)),
	}
	copy(p.data, data)
	return p
}

func (p *Pattern) String() string {
	sb := &strings.Builder{}
	for y := 0; y < 3; y++ {
		sb.WriteString("| ")
		for x := 0; x < 9; x++ {
			if x%2 == 0 {
				if p.data[x/2]-1 == y {
					sb.WriteRune('*')
				} else {
					sb.WriteRune(' ')
				}
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteString(" | \n")
	}
	sb.WriteString(fmt.Sprintf("%d:[", p.index))
	sb.WriteString(fmt.Sprintf("%d,%d,%d,%d,%d", p.data[0], p.data[1], p.data[2], p.data[3], p.data[4]))
	sb.WriteString("] \n")

	return sb.String()
}

// 0123456789ABCDE
// | * * * * * |
// | * * * * * |
// | * * * * * |
// 49:[1,2,3,4,5]

func main() {
	buf := make([]rune, 150*25)
	for i := 0; i < len(buf); i++ {
		buf[i] = ' '
	}

	patternObjs := make([]*Pattern, 0)

	index := 0
	for y := 0; y < 5; y++ {
		for x := 0; x < 10; x++ {
			pattern := NewPattern(sequence[index]+1, patterns[index])
			index++

			sx := x * 15
			sy := y * 5
			s := pattern.String()
			ss := strings.Split(s, "\n")
			for yy := 0; yy < len(ss); yy++ {
				for xx, r := range ss[yy] {
					buf[(sy+yy)*150+(sx+xx)] = r
				}
			}
			patternObjs = append(patternObjs, pattern)
		}
	}

	sb := &strings.Builder{}
	for y := 0; y < 25; y++ {
		for x := 0; x < 150; x++ {
			sb.WriteRune(buf[y*150+x])
		}
		sb.WriteRune('\n')
	}
	fmt.Println(sb.String())

	magic := make([]string, 15)
	for s := 0; s < len(magicSequence); s++ {
		p := patternObjs[magicSequence[s]]
		for i := 0; i < 15; i++ {
			x := i % 5
			y := i / 5
			if y == p.data[x]-1 {
				magic[i] = "1"
			} else {
				magic[i] = "0"
			}
		}
		fmt.Printf("{%s},\n", strings.Join(magic, ","))
	}

	for s := 0; s < len(magicSequence); s++ {
		i := magicSequence[s]
		p := patterns[i]
		fmt.Printf("{%d,%d,%d,%d,%d},\n", p[0], p[1], p[2], p[3], p[4])
	}
}
