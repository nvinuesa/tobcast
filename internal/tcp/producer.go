package tcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/underscorenico/tobcast/internal/data"
)

type Producer struct {
	Ports       []int
	connections []net.Conn
}

func NewProducer(broadcastPorts []int) *Producer {
	return &Producer{Ports: broadcastPorts, connections: []net.Conn{}}
}

func (p *Producer) Broadcast(message data.Message) error {
	p.connectBroacast()
	for _, c := range p.connections {
		log.Println("broadcasting msg '{}'", fmt.Sprintf("%v", message))

		bytes, err := getBytes(message)
		if err != nil {
			return err
		}
		c.Write(bytes)
	}
	return nil
}

func (p *Producer) connectBroacast() {
	if len(p.connections) == 0 {
		for _, s := range p.Ports {
			conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(s))
			if err != nil {
				log.Println("error connecting on port "+strconv.Itoa(s), err)
			} else {
				p.connections = append(p.connections, conn)
			}
		}
	}
}

func getBytes(message data.Message) ([]byte, error) {
	bytes, err := json.Marshal(message)
	if err != nil {
		log.Fatal("error marshalling message")
		return nil, err
	}
	bytes = append(bytes, '\n')
	return bytes, nil
}
