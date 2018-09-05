package main

import (
	"fmt"

	"github.com/master-g/playground/landmine"
)

func main() {
	landmine.Entry()
	b := make([]int, 1024)
	b = append(b, 99)
	fmt.Println("len:", len(b), "cap:", cap(b))
	b = append(b, 99)
	fmt.Println("len:", len(b), "cap:", cap(b))
}
