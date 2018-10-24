package main

import (
	"runtime"

	"github.com/master-g/playground/internal/abyss"
)

func main() {
	runtime.LockOSThread()

	abyss.Entry()
}
