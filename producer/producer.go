package producer

import (
	"encoding/json"
	"github.com/underscorenico/tobcast/data"
	"log"
	"net"
	"strconv"
	"time"
)

type Producer struct {
	Ports       []int
	Connections []net.Conn
}

func New(ports []int) *Producer {
	return &Producer{Ports: ports, Connections: []net.Conn{}}
}

func (p *Producer) Multicast(time time.Time, value interface{}) {
	if len(p.Connections) == 0 {
		for _, s := range p.Ports {
			conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(s))
			if err != nil {
				log.Println("error connecting on port " + strconv.Itoa(s), err)
			}
			p.Connections = append(p.Connections, conn)
		}
	}
	for _, c := range p.Connections {
		message := data.Message{
			Timestamp: time.UnixNano(),
			Value:     value,
		}
		bytes, err := json.Marshal(message)
		bytes = append(bytes, '\n')
		if err == nil {
			c.Write(bytes)
		}
	}
}
