package main

import (
	"fmt"
	"log"
	"strings"
)

var patterns = [][]int{
	{1, 1, 1}, {3, 2, 1}, {1, 2, 1}, {3, 3, 1}, {1, 1, 2}, {3, 1, 1}, {1, 2, 2}, {1, 3, 1}, {1, 3, 2},
	{2, 2, 2}, {1, 2, 3}, {2, 3, 2}, {2, 2, 1}, {2, 2, 3}, {2, 1, 1}, {2, 3, 3}, {3, 1, 3}, {3, 1, 2},
	{3, 3, 3}, {2, 1, 2}, {3, 2, 3}, {3, 3, 2}, {1, 1, 3}, {3, 2, 2}, {1, 3, 3}, {2, 1, 3}, {2, 3, 1},
}

var sequence = []int{
	0, 3, 6, 9, 12, 15, 18, 21, 24,
	1, 4, 7, 10, 13, 16, 19, 22, 25,
	2, 5, 8, 11, 14, 17, 20, 23, 26,
}

var magicSequence = []int{
	0, 9, 18,
	1, 10, 19,
	2, 11, 20,
	3, 12, 21,
	4, 13, 22,
	5, 14, 23,
	6, 15, 24,
	7, 16, 25,
	8, 17, 26,
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
		for x := 0; x < 5; x++ {
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
	sb.WriteString(fmt.Sprintf("%d,%d,%d", p.data[0], p.data[1], p.data[2]))
	sb.WriteString("] \n")

	return sb.String()
}

// 0123456789
// | * * * |
// | * * * |
// | * * * |
// 27:[3,3,3]

// 11 * 9

func main() {
	buf := make([]rune, 99*15)
	for i := 0; i < len(buf); i++ {
		buf[i] = ' '
	}

	patternChecker := make(map[string]bool)
	patternObjs := make([]*Pattern, 0)

	index := 0
	for y := 0; y < 3; y++ {
		for x := 0; x < 9; x++ {
			pattern := NewPattern(sequence[index]+1, patterns[index])

			sx := x * 11
			sy := y * 5
			s := pattern.String()
			ss := strings.Split(s, "\n")
			for yy := 0; yy < len(ss); yy++ {
				for xx, r := range ss[yy] {
					buf[(sy+yy)*99+(sx+xx)] = r
				}
			}
			patternObjs = append(patternObjs, pattern)

			patternChecker[fmt.Sprintf("%v", patterns[index])] = true
			index++
		}
	}

	if len(patternChecker) != len(patterns) {
		log.Fatalf("invalid patterns, expect %d but got %v", len(patterns), len(patternChecker))
	}

	sb := &strings.Builder{}
	for y := 0; y < 15; y++ {
		for x := 0; x < 99; x++ {
			sb.WriteRune(buf[y*99+x])
		}
		sb.WriteRune('\n')
	}
	fmt.Println(sb.String())
	fmt.Println("// indices")

	for s := 0; s < len(magicSequence); s++ {
		magic := make([]string, 9)
		for i := 0; i < len(magic); i++ {
			magic[i] = "0"
		}

		p := patternObjs[magicSequence[s]]
		for i, v := range p.data {
			magic[(v-1)*3+i] = "1"
		}
		fmt.Printf("{%s},\n", strings.Join(magic, ","))
	}

	fmt.Println("\n// magics")

	for s := 0; s < len(magicSequence); s++ {
		i := magicSequence[s]
		p := patterns[i]
		fmt.Printf("{%d,%d,%d},\n", p[0], p[1], p[2])
	}
}
