package mnemonic

import (
	"crypto/sha256"
	"errors"
)

type Mnemonic struct {
	BitArray     []byte
	BitLength    int
	checksumBits int
	checksum     byte
	Words        []string
}

const (
	checksumBytesPerBit = 32
	bitsPerWord         = 11
)

var (
	ErrorInvalidBitLength = errors.New("invalid bit length")
)

func NewMnemonicWords(bitLength int) (*Mnemonic, error) {
	if bitLength != 128 && bitLength != 160 && bitLength != 192 && bitLength != 224 && bitLength != 256 {
		return nil, ErrorInvalidBitLength
	}

	bitArray, err := NewRandomBits(bitLength)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(*bitArray)

	m := &Mnemonic{
		BitArray:     *bitArray,
		BitLength:    bitLength,
		checksumBits: bitLength / checksumBytesPerBit,
		checksum:     hash[0],
	}
	m.checksum = m.checksum >> uint(8-m.checksumBits)
	m.Words = make([]string, (m.BitLength+m.checksumBits)/bitsPerWord)

	for i := 0; i < cap(m.Words); i++ {

	}
	return m, nil
}
