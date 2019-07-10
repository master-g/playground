package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func premute() {
	poolorder := make([]int, 8)
	// loop through slice
	for i := 1; i < len(poolorder); i++ {
		j := rand.Int31n(int32(i + 1)) // choose a random index between 0~i
		poolorder[i] = poolorder[j]    // swap slice[i] and slice[j]
		poolorder[j] = i               // slice[j] = i
	}

	log.Info(poolorder)
}

func fileServer() {
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", ".", "directory of static files to host")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	log.Infof("Serving %v on HTTP port: %v", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func main() {
	// fileServer()
	v := 1
	for i := 0; i < 10; i++ {
		v *= 10
		fmt.Println(magic(v, 7))
	}
}

func stringifyWithComma(num int) string {
	strNum := strconv.Itoa(num)
	sb := &strings.Builder{}
	for i := len(strNum) - 1; i >= 0; i-- {
		sb.WriteByte(strNum[i])
		if i > 0 && i%3 == 0 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

func magic(num, maxLen int) string {
	if maxLen <= 0 {
		return stringifyWithComma(num)
	}

	raw := stringifyWithComma(num)
	postfix := "KMGTPEZY"
	curPostfixIndex := -1
	newValue := num
	for len(raw) > maxLen {
		newValue /= 1000
		curPostfixIndex++

		raw = stringifyWithComma(newValue)
		raw = fmt.Sprintf("%s%s", raw, string(postfix[curPostfixIndex]))
	}

	return raw
}
