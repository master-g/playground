package mnemonic

import (
	"fmt"
	"log"
)

func Entry() {
	m, err := NewMnemonicWords(128)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m.Words[0])
}
