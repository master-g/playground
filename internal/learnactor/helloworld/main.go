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

// 01. Hello World
// 最基本的 actor 基本概念与创建, 简单的消息发送接收
//
// actor.PID ProcessID
// * Id: actor 的唯一标识
// * Address: actor 地址
// * process: actor.Process 引用
//
// actor.Props 创建 actor 的时候提供的属性配置
// * spawner: Spawn 方法， 创建 process， mailbox 等 actor 内部组件， 注册 pid 等. 基本上仅在 actor 包内使用
// * producer: actor 工厂方法, 用户自行创建并初始化 actor， 然后返回 actor 实例
// * mailboxProducer: mailbox 工厂方法
// * guardianStrategy: 守护策略, 创建 actor 的时候如果指定守护策略， 一个父 actor 会被同时创建
// * supervisionStrategy: 监管策略， 创建 actor 后， 此 actor 会采用指定的监管策略来管理其子 actor
// * dispatcher: 指定邮箱调度器
// * receiverMiddleware: 消息接受者中间件
// * senderMiddleware: 消息发送者中间件
// * spawnMiddleware: spawner 中间件
// * receiverMiddlewareChain: 消息接受中间件调用函数
// * senderMiddlewareChain: 消息发送中间件调用函数
// * spawnMiddlewareChain: spawn 中间件调用函数
// * contextDecorator: context 装饰器
// * contextDecoratorChain: context 装饰器调用函数
//
// actor.Context actor 的上下文， 包括
// * infoPart： 此 context 当前相关联的actor， 父 actor
// * basePart: 接收超时， 子 actor， 监控指定 actor
// * messagePart: 获取当前处理的消息和消息头
// * senderPart: 获取当前消息的发送者， 向指定 pid 发送消息等
// * receiverPart: 实现处理消息的 Receive 方法
// * spawnerPart: 创建子 actor 的方法
// * stopperPart: 终止 actor 的方法
//
// actor.Process 合约规定的 actor 交互接口
// * SendUserMessage(pid *PID, message interface{}) 发送用户消息
// * SendSystemMessage(pid *PID, message interface{}) 发送系统消息
// * Stop(pid *PID) 停止指定的 process

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/master-g/playground/internal/learnactor/messages"
)

// 用户消息类型, http://proto.actor/docs/messages
type hello struct {
	Who string
}

// Actor 实现
type helloActor struct {
}

// Receive implement actor's Receive interface in helloActor
func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		fmt.Printf("Hello %v\n", msg.Who)
	}
}

// main entry
func main() {
	// actor 的创建
	// actor.Props 是 actor 的各种配置选项的集合
	// 要创建并取得 actor 的引用， 可以将 actor.Props 传递给 actor.Context.SpawnXXX 方法
	props := actor.PropsFromProducer(func() actor.Actor {
		return &helloActor{}
	})

	// https://github.com/AsynkronIT/protoactor-go/blob/dev/actor/root_context.go
	// actor.RootContext 包含了各种中间件， 消息头与监控策略
	rootContext := actor.EmptyRootContext
	// 创建新的 actor， 并返回其 PID
	// PID 包含了一个 actor 的唯一标识， actor 的地址和一个对 actor.Process 的引用
	pid := rootContext.Spawn(props)

	// 与旧版不同， 新版发送消息必须通过 context， 这样能够更好的支持中间件
	rootContext.Send(pid, &hello{Who: "World"})

	<-time.After(time.Second)

	// 在 props 中指定 actor 的 receive 方法
	props = actor.PropsFromFunc(func(c actor.Context) {
		switch msg := c.Message().(type) {
		case *hello:
			fmt.Printf("Hi %v\n", msg.Who)
		}
	})
	pid = rootContext.Spawn(props)
	rootContext.Send(pid, &hello{Who: "PropsFromFunc"})

	// remote
	remote.Start("127.0.0.1:9001")
	remotePid := actor.NewPID("127.0.0.1:8888", "bare")
	future := actor.EmptyRootContext.RequestFuture(remotePid, &messages.RemoteMessage{Content: "hi from remote client"}, 5*time.Second)
	resp, err := future.Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)

	<-time.After(time.Second)
}
