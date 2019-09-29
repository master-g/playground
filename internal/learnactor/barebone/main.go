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

package main

import (
	"fmt"
	"runtime"
	"time"

	"playground/internal/learnactor/barebone/barebone"
	"playground/pkg/signal"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"

	"github.com/oklog/run"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	remote.Start("127.0.0.1:8888")

	pid := barebone.RegisterBareBoneActor()

	var g run.Group

	// signal handler
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			signal.Start(cancel)
			return nil
		}, func(error) {
			// interrupt from other goroutine
			close(cancel)
		})
	}

	err := g.Run()
	if err != nil {
		fmt.Println(err)
	}
	actor.EmptyRootContext.Stop(pid)
	remote.Shutdown(true)
	fmt.Println("10 seconds before exit...")
	<-time.After(10 * time.Second)
}
