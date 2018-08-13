package main

import (
	"fmt"
	"log"

	"github.com/master-g/playground/mnemonic"
)

func main() {
	m, err := mnemonic.NewMnemonicWords(128)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m.Words)
}
