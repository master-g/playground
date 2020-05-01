package main

// http://www.btechsmartclass.com/data_structures/knuth-morris-pratt-algorithm.html

import (
	"fmt"
	"log"
	"os"
)

func genLPS(pattern string) []int {
	if len(pattern) == 0 {
		return nil
	}

	lps := make([]int, len(pattern))
	lps[0] = 0
	i := 0

	for j := 1; j < len(pattern); j++ {
		if pattern[i] == pattern[j] {
			i++
			lps[j] = i
		} else {
			lps[j] = 0
			if i != 0 {
				i = lps[i-1]
			}
		}
	}

	return lps
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s string", os.Args[0])
	}

	lps := genLPS(os.Args[1])
	fmt.Println(os.Args[1])
	for _, v := range lps {
		fmt.Print(v)
	}
	fmt.Println()
}
