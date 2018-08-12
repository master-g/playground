package mnemonic

import (
	"crypto/rand"
	"io"
)

const (
	BYTES_OF_128_BITS = 16
	BYTES_OF_160_BITS = 20
	BYTES_OF_192_BITS = 24
	BYTES_OF_224_BITS = 28
	BYTES_OF_256_BITS = 32
)

func byteSizeOfBits(bitLen int) (l int) {
	l = bitLen / 8
	if bitLen%8 != 0 {
		l++
	}
	return
}

func NewRandomBytes(size int) (*[]byte, error) {
	data := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, data[:])
	return &data, err
}

func NewRandomBits(bitLen int) (*[]byte, error) {
	return NewRandomBytes(byteSizeOfBits(bitLen))
}
