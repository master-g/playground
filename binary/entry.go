package binary

import (
	"fmt"
	"log"

	"github.com/master-g/playground/mnemonic"
)

func Entry() {
	m, err := mnemonic.NewMnemonicWords(128)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)
}
