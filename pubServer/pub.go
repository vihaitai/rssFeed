package pubServer

import (
	"fmt"
	"log"
	"sync"

	zmq "github.com/pebbe/zmq4"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var initOnce sync.Once
var _pubInstance *Publisher

type Publisher struct {
	*zmq.Socket
}

// 保证pub socket是一个单例
func initPub() *Publisher {
	initOnce.Do(func() {
		socket, err := zmq.NewSocket(zmq.PUB)
		if err != nil {
			panic(err)
		}
		_pubInstance = &Publisher{
			Socket: socket,
		}
	})
	return _pubInstance
}

var msgChan chan []byte

func PubMsg(msg []byte) error {
	// A send to a nil channel blocks forever
	// A receive from a nil channel blocks forever
	// A send to a closed channel panics
	// A receive from a closed channel returns the zero value immediately
	if msgChan == nil {
		return fmt.Errorf("start server before pubMsg")
	}
	msgChan <- msg
	return nil
}

func Start(endpoint string) {
	socket, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.Fatal("[ERROR]zmq.NewSocket error %+v\n", err)
	}
	if err := socket.Bind(endpoint); err != nil {
		log.Fatal("[ERROR]socket.Bind error %+v\n", err)
	}
	defer socket.Close()
	log.Printf("[DEBUG]start pub server on %s\n", endpoint)

	msgChan = make(chan []byte)
	defer close(msgChan)
	for msg := range msgChan {
		count, err := socket.SendMessage("article", msg)
		if err != nil {
			log.Printf("[ERROR]send message %s error %+v\n", string(msg), err)
		} else {
			log.Printf("[DEBUG]send message %s, send byte %d\n", string(msg), count)
		}
	}
}
