// Copyright © 2019 mg
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	// InterruptChan as signal to control goroutine to finish
	InterruptChan = make(chan struct{})
	// interruptSignals defines system signals to catch
	interruptSignals = []os.Signal{syscall.SIGTERM, os.Interrupt}
)

func Start() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, interruptSignals...)
	defer close(InterruptChan)

	<-c
}

func StartWithContext(ctx context.Context) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, interruptSignals...)
	defer close(InterruptChan)

	select {
	case <-c:
	case <-ctx.Done():
	}
}

// Start waiting UNIX system signal
func StartWithCanel(cancel chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, interruptSignals...)
	defer close(InterruptChan)

	select {
	case <-c:
	case <-cancel:
	}
}
