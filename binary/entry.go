package binary

import (
	"fmt"
	"log"

	"github.com/master-g/playground/mnemonic"
)

// Entry for this package
func Entry() {
	m, err := mnemonic.NewMnemonicWords(128)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)
}
