package binary

import (
	"fmt"
	"log"

	"playground/internal/mnemonic"
)

// Entry for this package
func Entry() {
	m, err := mnemonic.NewMnemonicWords(128)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)
}
