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

package barebone

import (
	"fmt"
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/mailbox"
	"github.com/master-g/playground/internal/learnactor/messages"
)

type BareActor struct {
	stopCh chan struct{}
}

func bareActorProducer() actor.Producer {
	return func() actor.Actor {
		return &BareActor{
			stopCh: make(chan struct{}),
		}
	}
}

func RegisterBareBoneActor() *actor.PID {
	props := actor.PropsFromProducer(bareActorProducer()).WithMailbox(mailbox.Bounded(20000))
	ctx := actor.EmptyRootContext
	pid, err := ctx.SpawnNamed(props, "bare")
	if err != nil {
		log.Fatal(err)
	}

	return pid
}

func (bareActor *BareActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		fmt.Println("actor started")
		bareActor.startTicking(ctx)
	case *actor.Stopping:
		fmt.Println("actor stopping")
		bareActor.stopTicking()
	case *TickMessage:
		fmt.Println("tick...")
	case *messages.RemoteMessage:
		msg := ctx.Message().(*messages.RemoteMessage)
		fmt.Println(msg.Content)
		ctx.Respond(&messages.RemoteMessage{Content: "response from remote"})
	}
}

func (bareActor *BareActor) startTicking(ctx actor.Context) {
	fmt.Println("start ticking")
	go func() {
		tick := time.NewTicker(time.Second)
		for {
			select {
			case <-bareActor.stopCh:
				fmt.Println("stopCh closed")
				return
			case <-tick.C:
				ctx.Send(ctx.Self(), &TickMessage{})
			}
		}
	}()
}

func (bareActor *BareActor) stopTicking() {
	select {
	case <-bareActor.stopCh:
		return
	default:
		close(bareActor.stopCh)
	}
}
