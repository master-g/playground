package binary11

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
)

func Entry() {
	entropy := "f8f7c5fa92d77b7cbf8909f936ad3329"
	h, err := hex.DecodeString(entropy)
	if err != nil {
		log.Fatal(err)
	}
	sum := sha256.Sum256(h)
	fmt.Println(hex.EncodeToString(sum[0:1]))

	for _, v := range h {
		fmt.Printf("%v ", ByteToBinary(v))
	}
	fmt.Println()

	a := make([]int, 12)
	fmt.Println("---------------")

	for i := 0; i < 11; i++ {
		startByteIdx := i * 11 / 8
		endByteIdx := (i + 1) * 11 / 8
		startBitIdx := i * 11 % 8
		endBitIdx := (i + 1) * 11 % 8
		var value int
		value = 0
		head := 1<<uint(8-startBitIdx) - 1
		head &= int(h[startByteIdx])
		head <<= uint(3 + startBitIdx)
		value |= head
		if endByteIdx-startByteIdx > 1 {
			middle := int(h[startByteIdx+1])
			middle <<= uint(endBitIdx)
			value |= middle
		}
		tail := uint(h[endByteIdx]) >> uint(8-endBitIdx)
		value |= int(tail)

		a[i] = value
		fmt.Println(fmt.Sprintf("%v,%v | %v %v", startByteIdx, endByteIdx, startBitIdx, endBitIdx))
	}
	for _, v := range a {
		fmt.Printf("%v ", IntToBinary11(v))
	}
	fmt.Println()
	for _, v := range a {
		fmt.Printf("%d ", v)
	}
}

func IntToBinary11(value int) string {
	sb := strings.Builder{}
	for i := 0; i < 11; i++ {
		mask := 1 << uint(10-i)
		if value&mask == 0 {
			sb.WriteString("0")
		} else {
			sb.WriteString("1")
		}
	}
	return sb.String()
}

func ByteToBinary(value byte) string {
	sb := strings.Builder{}
	for i := 0; i < 8; i++ {
		mask := byte(1 << uint(7-i))
		if value&mask == 0 {
			sb.WriteString("0")
		} else {
			sb.WriteString("1")
		}
	}
	return sb.String()
}
