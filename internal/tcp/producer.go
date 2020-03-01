package tcp

import (
	"encoding/json"
	"log"
	"net"
	"strconv"

	"github.com/underscorenico/tobcast/internal/data"
)

type Producer struct {
	Ports       []int
	connections map[int]net.Conn
}

func NewProducer(broadcastPorts []int) *Producer {
	connections := make(map[int]net.Conn)
	return &Producer{Ports: broadcastPorts, connections: connections}
}

func (p *Producer) Broadcast(message *data.Message) error {
	p.connectBroacast()
	for _, c := range p.connections {

		bytes, err := getBytes(*message)
		if err != nil {
			return err
		}
		c.Write(bytes)
	}
	return nil
}

func (p *Producer) Stop() {
	log.Println("closing tcp producer")
	for _, conn := range p.connections {
		conn.Close()
	}
}

func (p *Producer) connectBroacast() {
	if len(p.connections) != len(p.Ports) {
		for _, port := range p.Ports {
			if p.connections[port] == nil {
				conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(port))
				if err != nil {
					log.Println("error connecting on port "+strconv.Itoa(port), err)
				} else {
					p.connections[port] = conn
				}
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
