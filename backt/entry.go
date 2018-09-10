package backt

import (
	"fmt"
	"time"

	"github.com/jpillora/backoff"
)

// Entry for this package
func Entry() {
	b := &backoff.Backoff{
		// These are the defaults
		Min:    100 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2,
		Jitter: true,
	}

	for i := 0; i < 20; i++ {
		fmt.Printf("%s\n", b.Duration())
	}

	fmt.Printf("Reset!\n")
	b.Reset()

	fmt.Printf("%s\n", b.Duration())
}
