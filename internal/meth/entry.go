package meth

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

// Entry for this package
func Entry() {
	key, err := newECDSAPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(key[2:])

	fmt.Println("---------------------------------")

	convertPrivateKeyHexToAddress()
}

func newECDSAPrivateKey() (key string, err error) {
	var privateKey *ecdsa.PrivateKey
	privateKey, err = crypto.GenerateKey()
	if err != nil {
		return
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	key = hexutil.Encode(privateKeyBytes)
	return
}

func convertPrivateKeyHexToAddress() {
	privateKeyString := "0026fe00e006d5b318bd2aacd2b57e1819065042e3c329929e300853f8d787bb"
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:])

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)

	hash := sha3.NewKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))
}
