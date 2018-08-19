package mnemonic

import (
	"crypto/sha256"
	"errors"
	"math/big"
)

type Mnemonic struct {
	Entropy      []byte
	BitLength    uint
	checksumBits uint
	checksum     byte
	Binary       *big.Int
	Words        []string
}

const (
	checksumBytesPerBit = 32
	bitsPerWord         = 11
)

var (
	ErrorInvalidEntropyBitLength = errors.New("invalid bit length for entropy, must be [128, 256] and a multiple of 32")
)

func NewMnemonicWords(bitLength uint) (*Mnemonic, error) {
	if bitLength < 128 || bitLength > 256 || bitLength%32 != 0 {
		return nil, ErrorInvalidEntropyBitLength
	}

	bitArray, err := NewRandomBits(bitLength)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(*bitArray)

	m := &Mnemonic{
		Entropy:      *bitArray,
		BitLength:    bitLength,
		checksumBits: bitLength / checksumBytesPerBit,
	}
	m.checksum = hash[0] >> (8 - m.checksumBits)
	m.Words = make([]string, (m.BitLength+m.checksumBits)/bitsPerWord)

	bin := big.NewInt(0).SetBytes(*bitArray)
	bin = bin.Lsh(bin, m.checksumBits).Or(bin, big.NewInt(int64(m.checksum)))
	m.Binary = big.NewInt(0).Set(bin)

	mask := big.NewInt(0x7FF)
	for i := len(m.Words) - 1; i >= 0; i-- {
		v := big.NewInt(0).And(bin, mask)
		m.Words[i] = wordlist_en[v.Int64()]
		bin = bin.Rsh(bin, bitsPerWord)
	}

	return m, nil
}

func SeedFromMnemonicWords(words []string, passphrase string) []byte {
	seed := make([]byte, 0)

	return seed
}

func MnemonicToEntropy(words []string, passphrase string) []byte {
	return nil
}
