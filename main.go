package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	entropy := "6da3de924b850d4e1bcaa8174d232d91d041fa99"
	h, err := hex.DecodeString(entropy)
	if err != nil {
		log.Fatal(err)
	}
	sum := sha256.Sum256(h)
	fmt.Println(hex.EncodeToString(sum[0:1]))

	a := make([]int, 10)
	fmt.Println(cap(a))
}
