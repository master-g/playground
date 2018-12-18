package main

import (
	"math/rand"

	"github.com/master-g/playground/internal/cfgwatch"

	"github.com/sirupsen/logrus"
)

func premute() {
	poolorder := make([]int, 8)
	// loop through slice
	for i := 1; i < len(poolorder); i++ {
		j := rand.Int31n(int32(i + 1)) // choose a random index between 0~i
		poolorder[i] = poolorder[j]    // swap slice[i] and slice[j]
		poolorder[j] = i               // slice[j] = i
	}

	logrus.Info(poolorder)
}

func main() {
	cfgwatch.Execute()
}
