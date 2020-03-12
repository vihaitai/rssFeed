package pubServer

import (
	"fmt"
	"log"

	zmq "github.com/pebbe/zmq4"
)

type REQClient struct {
	socket *zmq.Socket
}

func NewREQClient(endpoint string) (*REQClient, error) {
	socket, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		log.Printf("[ERROR]zmq.NewSocket error %+v\n", err)
		return nil, err
	}

	if err := socket.Connect(endpoint); err != nil {
		log.Printf("[ERROR]zmq.Connect %s error %+v\n", endpoint, err)
		return nil, err
	}
	return &REQClient{socket: socket}, nil
}

func (client *REQClient) SendREQMsg(msg []byte) (string, error) {
	socket := client.socket
	total, err := socket.SendMessage(msg)
	if err != nil {
		log.Printf("[ERROR]zmq.SendMessage error %+v\n", err)
		return "", err
	}
	log.Printf("[DEBUG]zmq.SendMessage send total %d\n", total)

	messages, err := socket.RecvMessage(0)
	if err != nil {
		log.Printf("[ERROR]zmq.RecMessage error %+v\n", err)
		return "", err
	}
	if messages[0] != "ok" {
		err := fmt.Errorf("want ok got %+v", messages)
		log.Printf("[ERROR]%+v\n", err)
		return "", err
	}
	return messages[0], nil
}

func (client *REQClient) Close() {
	client.socket.Close()
}

func StartREP(endpoint string) {
	socket, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		log.Fatal("[ERROR]zmq.NewSocket zmq.REP failed %+v\n", err)
	}
	defer socket.Close()

	if err := socket.Bind(endpoint); err != nil {
		log.Fatal("[ERROR]zmq bind req server on %s failed %+v\n", endpoint, err)
	}
	log.Printf("[DEBUG]zmq bind req server on %s\n", endpoint)

	for {
		messages, err := socket.RecvMessageBytes(0)
		if err != nil {
			log.Printf("[ERROR]zmq.socket.RecvMessageBytes failed %+v\n", err)
			continue
		}
		for _, msg := range messages {
			log.Printf("[DEBUG]zmq.socket.RecvMessageBytes %s\n", msg)
			PubMsg(msg)
		}
		if _, err := socket.SendMessage("ok"); err != nil {
			log.Printf("[ERROR]zmq.socket.SendMessageBytes failed %+v\n", err)
		}
	}
}
