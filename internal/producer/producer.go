package producer

import (
	"encoding/json"
	"log"
	"net"
	"strconv"

	"github.com/underscorenico/tobcast/internal/data"
)

type Producer struct {
	BroadcastPorts     []int
	DeliverPorts       []int
	Connections        []net.Conn
	DeliverConnections []net.Conn
}

func New(broadcastPorts []int, deliverPorts []int) *Producer {
	return &Producer{BroadcastPorts: broadcastPorts, DeliverPorts: deliverPorts, Connections: []net.Conn{}}
}

func (p *Producer) Broadcast(message data.Message) {
	if len(p.Connections) == 0 {
		for _, s := range p.BroadcastPorts {
			conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(s))
			if err != nil {
				log.Println("error connecting on port "+strconv.Itoa(s), err)
			}
			p.Connections = append(p.Connections, conn)
		}
	}
	for _, c := range p.Connections {

		bytes, err := json.Marshal(message)
		bytes = append(bytes, '\n')
		if err != nil {
			log.Fatal("error marshalling message")
		}
		c.Write(bytes)
	}
}

func (p *Producer) Deliver(message data.Message) error {
	if len(p.DeliverConnections) == 0 {
		for _, s := range p.DeliverPorts {
			log.Println("connecting on port " + strconv.Itoa(s))
			conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(s))
			if err != nil {
				log.Println("error connecting on port "+strconv.Itoa(s), err)
			}
			p.DeliverConnections = append(p.DeliverConnections, conn)
		}
	}
	for _, c := range p.DeliverConnections {
		log.Println("sending msg")

		bytes, err := json.Marshal(message)
		bytes = append(bytes, '\n')
		if err == nil {
			c.Write(bytes)
		}
		return err
	}
	return nil
}
