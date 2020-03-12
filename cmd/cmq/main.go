package main

import (
	"rssFeed/pubServer"
	"time"

	zmq "github.com/pebbe/zmq4"
)

// pub模式下
// 发布topic: data
// 订阅topic

func main() {
	go func() {
		socket, err := zmq.NewSocket(zmq.SUB)
		if err != nil {
			panic(err)
		}
		defer socket.Close()

		if err := socket.Connect("tcp://127.0.0.1:8001"); err != nil {
			panic(err)
		}
		if err := socket.SetSubscribe("article"); err != nil {
			panic(err)
		}
		for {
			messages, err := socket.RecvMessage(0)
			if err != nil {
				panic(err)
			}
			topic, data := messages[0], messages[1]
			println(topic, data)
			time.Sleep(time.Duration(500 * time.Millisecond))
		}
	}()

	client, err := pubServer.NewREQClient("tcp://127.0.0.1:8002")
	if err != nil {
		panic(err)
	}
	defer client.Close()
	client.SendREQMsg([]byte("hello word"))

	time.Sleep(time.Duration(3 * time.Second))
}
